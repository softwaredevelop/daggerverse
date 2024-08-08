// Package main provides a generated module for Hello functions.
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.
package main

import (
	"context"
	"dagger/hello/internal/dagger"
)

// New creates a new instance of the Hello struct.
// If the ctr parameter is nil, it will create a new dagger.Container using the "busybox:uclibc" image.
// The stringArg parameter is an optional argument that specifies a string value.
// If not provided, it will default to "Hello, Daggerverse!".
// The function returns a pointer to the created Hello struct.
func New(
	// ctr is an optional argument that specifies a container.
	// +optional
	ctr *dagger.Container,
	// stringArg is an optional argument that specifies a string value.
	// +optional
	// +default="Hello, Daggerverse!"
	stringArg string,
) *Hello {
	if ctr == nil {
		ctr = dag.Container().
			From("busybox:uclibc")
	}
	return &Hello{
		Ctr:       *ctr,
		StringArg: stringArg,
	}
}

// Hello is a struct that provides echo functions.
type Hello struct {
	Ctr       dagger.Container
	StringArg string
}

// HelloString returns the string value provided to the Hello struct.
func (m *Hello) HelloString(
	ctx context.Context,
) string {
	return m.StringArg
}

// HelloContainer executes a container with the provided string value.
func (m *Hello) HelloContainer(
	ctx context.Context,
) (string, error) {
	return dag.Container().
		From("cgr.dev/chainguard/wolfi-base:latest").
		WithExec([]string{"echo", m.StringArg}).
		Stdout(ctx)
}
