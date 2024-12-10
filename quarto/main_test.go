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

func Test_Quatro(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_quatro_version", func(t *testing.T) {
		t.Parallel()
		container := base(nil)
		require.NotNil(t, container)

		out, err := container.
			WithExec([]string{"quarto", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Regexp(t, `\d+\.\d+\.\d+`, out)
	})
}

func base(
	image *string,
) *dagger.Container {

	defaultImageRepository := "ghcr.io/quarto-dev/quarto"
	var ctr *dagger.Container

	if image != nil {
		ctr = c.Container().From(*image)
	} else {
		ctr = c.Container().From(defaultImageRepository)
	}

	return ctr
}
