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

func Test_Revive(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_revive_check", func(t *testing.T) {
		t.Parallel()
		container := base()
		require.NotNil(t, container)

		_, err := container.
			WithMountedDirectory("/go", c.Host().Directory("./test/testdata/")).
			WithWorkdir("/go").
			WithExec([]string{"revive", "-set_exit_status", "./..."}).
			Stdout(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "don't use underscores in Go names")
	})
	t.Run("Test_revive_version", func(t *testing.T) {
		t.Parallel()
		container := base()
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"revive", "-version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "version")
	})
}

func base() *dagger.Container {
	image := c.Container().
		From("ghcr.io/mgechev/revive:latest")

	return c.Container().
		From("golang:alpine").
		WithFile("/usr/bin/revive", image.File("/usr/bin/revive"))
}
