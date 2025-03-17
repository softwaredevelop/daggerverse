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

func Test_Hello(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_hello_container", func(t *testing.T) {
		t.Parallel()

		client, err := getClient()
		require.NoError(t, err)
		t.Cleanup(func() { client.Close() })

		stringArg := "Hello, Daggerverse!"
		out, err := client.Container().
			From("busybox:uclibc").
			WithExec([]string{"echo", stringArg}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Equal(t, "Hello, Daggerverse!\n", out)
	})
}
