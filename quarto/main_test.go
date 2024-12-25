package main_test

import (
	"context"
	"flag"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

var c *dagger.Client

func TestMain(m *testing.M) {
	flag.Parse()

	ctx := context.Background()

	c, _ = dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	defer c.Close()

	code := m.Run()
	defer c.Close()
	os.Exit(code)
}

func Test_Quarto(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_quarto_render_export", func(t *testing.T) {
		t.Parallel()
		container := base("ghcr.io/quarto-dev/quarto-full", nil)
		require.NotNil(t, container)

		out, err := container.
			WithDirectory("/tmp", c.Host().Directory("./test/testdata/")).
			WithWorkdir("/tmp").
			WithExec([]string{"quarto", "render"}).
			Directory("/tmp/_book").Sync(ctx)
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
	t.Run("Test_quarto_tlmgr_mirror_and_render", func(t *testing.T) {
		t.Parallel()
		container := base("ghcr.io/quarto-dev/quarto-full", nil)
		require.NotNil(t, container)

		out, err := container.
			WithDirectory("/tmp", c.Host().Directory("./test/testdata/")).
			WithWorkdir("/tmp").
			WithExec([]string{"tlmgr", "option", "repository", "http://mirror.ctan.org/systems/texlive/tlnet"}).
			WithExec([]string{"quarto", "render"}).
			WithExec([]string{"ls", "-1", "_book"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\.pdf\s*$`, out)
	})
	t.Run("Test_quarto_full_render", func(t *testing.T) {
		t.Parallel()
		container := base("ghcr.io/quarto-dev/quarto-full", nil)
		require.NotNil(t, container)

		out, err := container.
			WithDirectory("/tmp", c.Host().Directory("./test/testdata/")).
			WithWorkdir("/tmp").
			WithExec([]string{"quarto", "render"}).
			WithExec([]string{"ls", "-1", "_book"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\.pdf\s*$`, out)
	})
	t.Run("Test_quarto_full_version", func(t *testing.T) {
		t.Parallel()
		container := base("ghcr.io/quarto-dev/quarto-full", nil)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"quarto", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
	t.Run("Test_quarto_add_extensions", func(t *testing.T) {
		t.Parallel()
		container := base("", []string{"quarto-ext/latex-environment", "quarto-ext/include-code-files"})
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"quarto", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
	t.Run("Test_quarto_version", func(t *testing.T) {
		t.Parallel()
		container := base("", nil)
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
) *dagger.Container {

	defaultImageRepository := "ghcr.io/quarto-dev/quarto"
	var ctr *dagger.Container

	if image != "" {
		ctr = c.Container().From(image)
	} else {
		ctr = c.Container().From(defaultImageRepository)
	}

	for _, ext := range extensions {
		ctr = ctr.WithExec([]string{"quarto", "add", "--no-prompt", ext})
	}

	return ctr
}
