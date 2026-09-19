// A Dagger module for linting Go source code using Revive.
//
// Revive is a fast, configurable, extensible, and flexible linter for Go.
// This module runs revive with exit status enforcement on target Go packages.
package main

import (
	"dagger/revive/internal/dagger"
)

const (
	defaultImageRepository = "ghcr.io/mgechev/revive:latest"
)

// Revive provides functions for running the Revive linter.
type Revive struct {
	// +private
	Image string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Revive struct.
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Revive {
	return &Revive{
		Image: image,
	}
}

// container returns the underlying Dagger container, lazily initialized.
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

// Check runs the revive linter on the target Go packages.
func (m *Revive) Check(
	// Source directory containing the Go module.
	source *dagger.Directory,
	// Optional configuration file (e.g. revive.toml).
	// +optional
	config *dagger.File,
	// Target Go packages or file pattern to check (defaults to "./...").
	// +default="./..."
	// +optional
	packages string,
) *dagger.Container {
	ctr := m.container().
		WithMountedDirectory("/work", source).
		WithWorkdir("/work")

	args := []string{"/revive", "-set_exit_status"}

	if config != nil {
		ctr = ctr.WithFile("/etc/revive.toml", config)
		args = append(args, "-config", "/etc/revive.toml")
	}

	args = append(args, packages)

	return ctr.WithExec(args)
}
