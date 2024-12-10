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
)

const defaultImageRepository = "ghcr.io/quarto-dev/quarto"

// Quarto is a Dagger module that provides functions for running Quarto
type Quarto struct {
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Quarto struct
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Quarto {
	var ctr *dagger.Container

	if image != "" {
		ctr = dag.Container().From(image)
	} else {
		ctr = dag.Container().From(defaultImageRepository)
	}

	return &Quarto{ctr}
}

// Cli runs the quarto cli
func (m *Quarto) Cli(
	// commands to run
	args []string,
) *dagger.Container {
	return m.Ctr.
		WithExec(args)
}
