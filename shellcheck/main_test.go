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

func Test_Shellcheck(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_mounted_host_directory_check", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		_, err = container.
			WithDirectory("/tmp", client.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithExec([]string{"sh", "-c", "find . -type f -name '*.sh' -print0 | xargs -0 shellcheck"}).
			Stdout(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 123")
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
		require.Regexp(t, `version:\s*v\d+\.\d+\.\d+(?:-\d+-[a-f0-9]+)?`, out)
	})
}

func base(
	image string,
	client *dagger.Client,
) *dagger.Container {

	defaultImageRepository := "koalaman/shellcheck-alpine"
	var ctr *dagger.Container

	if image == "" {
		image = defaultImageRepository
	}

	ctr = client.Container().From(image)

	return ctr
}
