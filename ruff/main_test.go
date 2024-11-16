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

func Test_Ruff(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_ruff_with_config", func(t *testing.T) {
		t.Parallel()
		container := base()
		require.NotNil(t, container)

		_, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithMountedFile("/.config/.ruff.toml", c.Host().File("./test/testdata/.config/.ruff.toml")).
			WithExec([]string{"/ruff", "check", "--config", "/.config/.ruff.toml"}).
			Stdout(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "F821")
	})
	t.Run("Test_ruff_host_directory_check_file_error", func(t *testing.T) {
		t.Parallel()
		container := base()
		require.NotNil(t, container)

		_, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithExec([]string{"/ruff", "check"}).
			Stdout(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "F401")
	})
	t.Run("Test_ruff_error", func(t *testing.T) {
		t.Parallel()
		container := base()
		require.NotNil(t, container)

		_, err := container.
			WithNewFile(
				"example.py",
				"import os",
			).
			WithExec([]string{"/ruff", "check", "example.py"}).
			Stdout(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "F401")
	})
	t.Run("Test_ruff_version", func(t *testing.T) {
		t.Parallel()
		container := base()
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"/ruff", "version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
}

func base() *dagger.Container {
	return c.Container().
		From("ghcr.io/astral-sh/ruff:latest")
}
