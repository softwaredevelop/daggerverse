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

func Test_Hadolint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_hadolint_with_config", func(t *testing.T) {
		t.Parallel()
		container := base("hadolint/hadolint:latest-alpine")
		require.NotNil(t, container)

		_, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithMountedFile("/.config/.hadolint.yaml", c.Host().File("./test/testdata/.config/.hadolint.yaml")).
			WithExec([]string{"sh", "-c", "find . -type f \\( -name 'Dockerfile' -o -name 'Dockerfile.*' \\) -print0 | xargs -0 hadolint --config /.config/.hadolint.yaml"}).
			Stdout(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 123")
	})
	t.Run("Test_hadolint_dockerfile_error", func(t *testing.T) {
		t.Parallel()
		container := base("hadolint/hadolint:latest-alpine")
		require.NotNil(t, container)

		_, err := container.
			WithMountedDirectory("/tmp", c.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithExec([]string{"sh", "-c", "find . -type f \\( -name 'Dockerfile' -o -name 'Dockerfile.*' \\) -print0 | xargs -0 hadolint"}).
			Stdout(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 123")
	})
	t.Run("Test_hadolint_error", func(t *testing.T) {
		t.Parallel()
		container := base("hadolint/hadolint:latest-alpine")
		require.NotNil(t, container)

		_, err := container.
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
		container := base("")
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
) *dagger.Container {

	defaultImageRepository := "hadolint/hadolint"
	var ctr *dagger.Container

	if image == "" {
		image = defaultImageRepository
	}

	ctr = c.Container().From(image)

	return ctr
}
