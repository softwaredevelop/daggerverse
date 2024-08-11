package main_test

import (
	"context"
	"flag"
	"os"
	"path/filepath"
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

func Test_Editorconfig(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_mounted_host_directory_without_git_directory", func(t *testing.T) {
		t.Parallel()
		container := base(c)
		require.NotNil(t, container)

		dir, err := os.Getwd()
		require.NoError(t, err)
		require.NotEmpty(t, dir)

		out, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory(filepath.Join(dir, ".."), dagger.HostDirectoryOpts{
				Exclude: []string{".git"},
			})).
			WithWorkdir("/tmp").
			WithExec([]string{"/editorconfig-checker", "-dry-run"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotNil(t, out)
	})
	t.Run("Test_mounted_host_directory", func(t *testing.T) {
		t.Parallel()
		container := base(c)
		require.NotNil(t, container)

		out, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory("./test/testdata/")).
			WithWorkdir("/tmp").
			WithExec([]string{"/editorconfig-checker", "-dry-run"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "t.txt")
	})
	t.Run("Test_editorconfig_checker_help", func(t *testing.T) {
		t.Parallel()
		container := base(c)
		require.NotNil(t, container)

		_, err := container.
			WithExec([]string{"/editorconfig-checker", "-help"}).
			Stdout(ctx)
		require.NoError(t, err)
	})
}

func base(c *dagger.Client) *dagger.Container {
	install := c.
		Container().
		From("golang:alpine").
		WithExec([]string{
			"go", "install",
			"github.com/editorconfig-checker/editorconfig-checker/cmd/editorconfig-checker@latest",
		})

	return c.Container().
		WithFile("/", install.File("/go/bin/editorconfig-checker"))
}
