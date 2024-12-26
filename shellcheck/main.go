// A generated module for Shellcheck functions
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
	"dagger/shellcheck/internal/dagger"
)

const (
	defaultImageRepository = "koalaman/shellcheck-alpine:latest"
)

// Shellcheck is a module for checking shell scripts.
type Shellcheck struct {
	// +private
	Image string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Shellcheck struct
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Shellcheck {
	return &Shellcheck{
		Image: image,
	}
}

// Container returns the underlying Dagger container
func (m *Shellcheck) Container() *dagger.Container {
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

// Check runs the shellcheck command.
func (m *Shellcheck) Check(
	// source is an optional argument that specifies a directory.
	source *dagger.Directory,
) *dagger.Container {
	return m.Container().
		WithMountedDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithExec([]string{"sh", "-c", "find . -type f -name '*.sh' -print0 | xargs -0 shellcheck"})
}
