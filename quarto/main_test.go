package main_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

const (
	defaultImageRepository = "ghcr.io/quarto-dev/quarto"
	tlmgrUpdateURL         = "https://mirror.ctan.org/systems/texlive/tlnet/update-tlmgr-latest.sh"
)

func getClient() (*dagger.Client, error) {
	ctx := context.Background()
	return dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
}

func Test_Quarto(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_quarto_render_export", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("ghcr.io/quarto-dev/quarto-full", nil, nil, client)
		require.NotNil(t, container)

		out, err := container.
			WithDirectory("/tmp", client.Host().Directory("./test/testdata/")).
			WithWorkdir("/tmp").
			WithExec([]string{"quarto", "render"}).
			Directory("/tmp/_output").Sync(ctx)
		require.NoError(t, err)
		require.NotNil(t, out)

		outputDir := "./test/outputdata/"
		_, err = out.Export(ctx, outputDir)
		require.NoError(t, err)

		files, err := os.ReadDir(outputDir)
		require.NoError(t, err)

		for _, file := range files {
			require.Regexp(t, `\.pdf$`, file)
		}

		err = os.RemoveAll(outputDir)
		require.NoError(t, err)
	})
	t.Run("Test_quarto_full_render", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("ghcr.io/quarto-dev/quarto-full", nil, nil, client)
		require.NotNil(t, container)

		out, err := container.
			WithDirectory("/tmp", client.Host().Directory("./test/testdata/")).
			WithWorkdir("/tmp").
			WithExec([]string{"quarto", "render"}).
			WithExec([]string{"ls", "-1", "_output"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\.pdf\s*$`, out)
	})
	t.Run("Test_quarto_full_version", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("ghcr.io/quarto-dev/quarto-full", nil, nil, client)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"quarto", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
	t.Run("Test_quarto_add_latex_packages", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base(
			"ghcr.io/quarto-dev/quarto-full",
			nil,
			[]string{"babel-english", "babel-german"},
			client,
		)
		require.NotNil(t, container)

		for _, pkg := range []string{"babel-english", "babel-german"} {
			out, err := container.
				WithExec([]string{"tlmgr", "info", pkg}).
				Stdout(ctx)
			require.NoError(t, err, "Failed to get info for package %s", pkg)
			require.Regexp(t, `installed:\s+Yes`, out, "Package %s is not installed", pkg)
		}
	})
	t.Run("Test_quarto_add_extensions", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base(
			"",
			[]string{"quarto-ext/latex-environment", "quarto-ext/include-code-files"},
			nil,
			client,
		)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"quarto", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
	t.Run("Test_quarto_version", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", nil, nil, client)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"quarto", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
}

func base(
	image string,
	extensions []string,
	latexpackages []string,
	client *dagger.Client,
) *dagger.Container {

	var ctr *dagger.Container

	if image == "" {
		image = defaultImageRepository
	}

	ctr = client.Container().From(image)

	if strings.Contains(image, "quarto-full") {
		ctr = ctr.WithExec([]string{
			"sh", "-c",
			"curl -fsSL " + tlmgrUpdateURL + " -o update-tlmgr-latest.sh && sh update-tlmgr-latest.sh -- --update",
		})

		for _, pkg := range latexpackages {
			ctr = ctr.WithExec([]string{"tlmgr", "install", pkg})
		}
	}

	for _, ext := range extensions {
		ctr = ctr.WithExec([]string{"quarto", "add", "--no-prompt", ext})
	}

	return ctr
}
