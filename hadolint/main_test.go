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

func Test_Hadolint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_hadolint_with_config", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("hadolint/hadolint:latest-alpine", client)
		require.NotNil(t, container)

		_, err = container.
			WithMountedDirectory("/tmp", client.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithMountedFile("/.config/.hadolint.yaml", client.Host().File("./test/testdata/.config/.hadolint.yaml")).
			WithExec([]string{"sh", "-c", "find . -type f \\( -name 'Dockerfile' -o -name 'Dockerfile.*' \\) -print0 | xargs -0 hadolint --config /.config/.hadolint.yaml"}).
			Stdout(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 123")
	})
	t.Run("Test_hadolint_dockerfile_error", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("hadolint/hadolint:latest-alpine", client)
		require.NotNil(t, container)

		_, err = container.
			WithMountedDirectory("/tmp", client.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithExec([]string{"sh", "-c", "find . -type f \\( -name 'Dockerfile' -o -name 'Dockerfile.*' \\) -print0 | xargs -0 hadolint"}).
			Stdout(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 123")
	})
	t.Run("Test_hadolint_error", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("hadolint/hadolint:latest-alpine", client)
		require.NotNil(t, container)

		_, err = container.
			WithNewFile(
				"Dockerfile",
				"FROM docker.io/library/alpine:$VARIANT",
			).
			WithNewFile(
				"Dockerfile.test",
				"FROM docker.io/library/alpine:latest",
			).
			WithExec([]string{"sh", "-c", "find . -type f \\( -name 'Dockerfile' -o -name 'Dockerfile.*' \\) -print0 | xargs -0 hadolint"}).
			Stdout(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 123")
	})
	t.Run("Test_hadolint_version", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"hadolint", "-v"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
}

func base(
	image string,
	client *dagger.Client,
) *dagger.Container {

	defaultImageRepository := "hadolint/hadolint"
	var ctr *dagger.Container

	if image == "" {
		image = defaultImageRepository
	}

	ctr = client.Container().From(image)

	return ctr
}
