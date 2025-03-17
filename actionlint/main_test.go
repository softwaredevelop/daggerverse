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

func Test_Actionlint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_actionlint_error", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		_, err = container.
			WithMountedDirectory("/tmp", client.Host().Directory("./test/testdata/")).
			WithWorkdir("/tmp").
			WithExec([]string{"sh", "-c", "find . -type f -name '*.yml' -print0 | xargs -0 actionlint"}).
			Stderr(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 123")
	})
	t.Run("Test_actionlint_version", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"actionlint", "-version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
	t.Run("Test_shellcheck_version", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"shellcheck", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "version")
	})
}

func base(
	image string,
	client *dagger.Client,
) *dagger.Container {

	defaultImageRepository := "rhysd/actionlint"
	var ctr *dagger.Container

	if image != "" {
		ctr = client.Container().From(image)
	} else {
		ctr = client.Container().From(defaultImageRepository)
	}

	return ctr
}
