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

func Test_Revive(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_revive_check", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		_, err = container.
			WithMountedDirectory("/go", client.Host().Directory("./test/testdata/")).
			WithWorkdir("/go").
			WithExec([]string{"/revive", "-set_exit_status", "./..."}).
			Stdout(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 1")
	})
	t.Run("Test_revive_version", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"/revive", "-version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `Version:\s*\d+\.\d+\.\d+`, out)
	})
}

func base(
	image string,
	client *dagger.Client,
) *dagger.Container {

	defaultImageRepository := "ghcr.io/mgechev/revive"
	var ctr *dagger.Container

	if image == "" {
		image = defaultImageRepository
	}

	ctr = client.Container().From(image)

	return ctr
}
