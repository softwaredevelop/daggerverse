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

func Test_Hello(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_hello_container", func(t *testing.T) {
		t.Parallel()
		stringArg := "Hello, Daggerverse!"
		out, err := c.Container().
			From("busybox:uclibc").
			WithExec([]string{"echo", stringArg}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Equal(t, "Hello, Daggerverse!\n", out)
	})
}
