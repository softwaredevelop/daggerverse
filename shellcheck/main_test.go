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

func Test_Shellcheck(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_mounted_host_directory_check", func(t *testing.T) {
		t.Parallel()
		container := base()
		require.NotNil(t, container)

		_, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithExec([]string{"sh", "-c", "find . -type f -name '*.sh' -print0 | xargs -0 shellcheck"}).
			Stdout(ctx)
		require.Error(t, err)
		errorIDs := []string{"SC2283", "SC2154", "SC2086"}
		for _, id := range errorIDs {
			require.Contains(t, err.Error(), id)
		}
	})
	t.Run("Test_shellcheck_version", func(t *testing.T) {
		t.Parallel()
		container := base()
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"shellcheck", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "version")
	})
}

func base() *dagger.Container {
	return c.
		Container().
		From("koalaman/shellcheck-alpine:latest").
		WithoutEntrypoint()
}
