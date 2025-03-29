package main_test

import (
	"context"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func getClient() (*dagger.Client, error) {
	ctx := context.Background()
	return dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
}

func Test_Ruff(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_ruff_with_config", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		_, err = container.
			WithMountedDirectory("/tmp", client.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithMountedFile("/.config/.ruff.toml", client.Host().File("./test/testdata/.config/.ruff.toml")).
			WithExec([]string{"/ruff", "check", "--config", "/.config/.ruff.toml"}).
			Stdout(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 1")
	})
	t.Run("Test_ruff_host_directory_check_file_error", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		_, err = container.
			WithMountedDirectory("/tmp", client.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithExec([]string{"/ruff", "check"}).
			Stdout(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "exit code: 1")
	})
	t.Run("Test_ruff_error", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		_, err = container.
			WithNewFile(
				"example.py",
				"import os",
			).
			WithExec([]string{"/ruff", "check", "example.py"}).
			Stderr(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "exit code: 1")
	})
	t.Run("Test_ruff_version", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"/ruff", "version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
}

func base(
	image string,
	client *dagger.Client,
) *dagger.Container {

	defaultImageRepository := "ghcr.io/astral-sh/ruff"
	var ctr *dagger.Container

	if image != "" {
		ctr = client.Container().From(image)
	} else {
		ctr = client.Container().From(defaultImageRepository)
	}

	return ctr
}
