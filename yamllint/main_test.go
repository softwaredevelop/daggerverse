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

func Test_yamllint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_yamllint_config_with_host_directory", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		_, err = container.
			WithMountedDirectory("/tmp", client.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithExec([]string{"sh", "-c", "find . -type f \\( -name '*.yaml' -o -name '*.yml' \\) -print0 | xargs -0 yamllint -c /tmp/.config/.yamllint"}).
			Stderr(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 1")
	})
	t.Run("Test_yamllint_directory_with_host_directory", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		_, err = container.
			WithMountedDirectory("/tmp", client.Host().Directory("./test/testdata")).
			WithWorkdir("/tmp").
			WithExec([]string{"yamllint",
				"--config-data",
				"{extends: default, rules: {line-length: {level: warning}}}",
				"--no-warnings",
				"."}).
			Stderr(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 1")
	})
	t.Run("Test_yamllint_error", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		// editorconfig-checker-disable
		badYAML := `
	foo: "bar"
	  baz: "qux"
	- item1
	- item2
	key : value
	`
		// editorconfig-checker-enable

		container = container.WithNewFile("/tmp/bad.yaml", badYAML)
		_, err = container.
			WithExec([]string{"yamllint",
				"--config-data",
				"{extends: default, rules: {line-length: {level: warning}}}",
				"--no-warnings",
				"/tmp/bad.yaml"}).
			Stderr(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "exit code: 1")
	})
	t.Run("Test_yamllint_version", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		container := base("", client)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"yamllint", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
}

func base(
	image string,
	client *dagger.Client,
) *dagger.Container {

	defaultImageRepository := "pipelinecomponents/yamllint"
	var ctr *dagger.Container

	if image == "" {
		image = defaultImageRepository
	}

	ctr = client.Container().From(image)

	return ctr
}
