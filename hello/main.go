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
	"dagger/hello/internal/dagger"
)

const (
	defaultImageRepository = "busybox:uclibc"
)

// Hello is a struct that provides echo functions.
type Hello struct {
	// +private
	Image string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Hello struct.
// If the ctr parameter is nil, it will create a new dagger.Container using the "busybox:uclibc" image.
// The stringArg parameter is an optional argument that specifies a string value.
// If not provided, it will default to "Hello, Daggerverse!".
// The function returns a pointer to the created Hello struct.
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Hello {
	return &Hello{
		Image: image,
	}
}

// Container returns the underlying Dagger container
func (m *Hello) container() *dagger.Container {
	if m.Ctr != nil {
		return m.Ctr
	}

	image := m.Image
	if image == "" {
		image = defaultImageRepository
	}

	m.Ctr = dag.Container().From(image)
	return m.Ctr
}

// HelloString returns the string value provided to the Hello struct.
func (m *Hello) HelloString(
	// StringArg is the string value to be printed.
	StringArg string,
) *dagger.Container {
	return m.container().
		WithExec([]string{"echo", StringArg})
}

// HelloContainer executes a container with the provided string value.
func (m *Hello) HelloContainer(
	// StringArg is the string value to be printed.
	StringArg string,
) *dagger.Container {
	return m.container().
		From("cgr.dev/chainguard/wolfi-base:latest").
		WithExec([]string{"echo", StringArg})
}
