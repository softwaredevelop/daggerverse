// A generated module for Quarto functions
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
	"dagger/quarto/internal/dagger"
	"strings"
)

const (
	defaultImageRepository = "ghcr.io/quarto-dev/quarto"
)

// Quarto is a module for running Quarto
type Quarto struct {
	// +private
	Image string
	// +private
	Extensions []string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Actionlint struct
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
	// List of extensions to add to the container.
	// +optional
	extensions []string,
) *Quarto {
	return &Quarto{
		Image:      image,
		Extensions: extensions,
	}
}

// Container returns the underlying Dagger container
func (m *Quarto) Container() *dagger.Container {
	if m.Ctr != nil {
		return m.Ctr
	}

	image := m.Image
	if image == "" {
		image = defaultImageRepository
	}

	ctr := dag.Container().From(image)

	for _, ext := range m.Extensions {
		ctr = ctr.WithExec([]string{"quarto", "add", "--no-prompt", ext})
	}

	m.Ctr = ctr
	return m.Ctr
}

// Build runs the quarto render command exporting to a directory
func (m *Quarto) Build(
	// source directory.
	source *dagger.Directory,
) *dagger.Directory {

	return m.Container().
		WithDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithExec([]string{"quarto", "render"}).Directory("/tmp/_book")
}

// Render runs the quarto render command
func (m *Quarto) Render(
	// source directory.
	source *dagger.Directory,
) *dagger.Container {

	return m.Container().
		WithDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithExec([]string{"quarto", "render"})
}

// Cli runs the quarto cli
func (m *Quarto) Cli(
	// commands to run
	args string,
) *dagger.Container {
	parsedArgs := strings.Split(args, " ")

	return m.Container().
		WithExec(parsedArgs)
}
