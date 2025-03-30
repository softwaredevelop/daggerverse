// A generated module for Revive functions
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
	"dagger/revive/internal/dagger"
)

const (
	defaultImageRepository = "ghcr.io/mgechev/revive:latest"
)

// Revive is a Dagger module that provides functions for running Revive linter
type Revive struct {
	Image string
	Ctr   *dagger.Container
}

// New creates a new instance of the Revive struct
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Revive {
	return &Revive{
		Image: image,
	}
}

// Container returns the underlying Dagger container
func (m *Revive) container() *dagger.Container {
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

// Check runs the revive command
func (m *Revive) Check(
	// Source directory
	source *dagger.Directory,
) *dagger.Container {
	return m.container().
		WithMountedDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithExec([]string{"/revive", "-set_exit_status", "./..."})
}
