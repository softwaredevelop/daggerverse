package main_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func getClient() (*dagger.Client, error) {
	ctx := context.Background()
	return dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
}

func Test_Editorconfig(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_mounted_host_directory_without_git_directory", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		dir, err := os.Getwd()
		require.NoError(t, err)
		require.NotEmpty(t, dir)

		out, err := container.
			WithMountedDirectory("/tmp", client.Host().Directory(filepath.Join(dir, ".."), dagger.HostDirectoryOpts{
				Exclude: []string{".git"},
			})).
			WithWorkdir("/tmp").
			WithExec([]string{"ec", "-dry-run"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotNil(t, out)
		require.NotRegexp(t, `\.git/`, out)
	})
	t.Run("Test_mounted_host_directory", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		out, err := container.
			WithMountedDirectory("/tmp", client.Host().Directory("./test/testdata/")).
			WithWorkdir("/tmp").
			WithExec([]string{"ec", "-dry-run"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `t\.txt`, out)
	})
	t.Run("Test_editorconfig-checker_help", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		_, err = container.
			WithExec([]string{"ec", "-help"}).
			Stderr(ctx)
		require.NoError(t, err)
	})
}

func base(
	image string,
	client *dagger.Client,
) *dagger.Container {

	defaultImageRepository := "mstruebing/editorconfig-checker"
	var ctr *dagger.Container

	if image == "" {
		image = defaultImageRepository
	}

	ctr = client.Container().From(image)

	return ctr
}
