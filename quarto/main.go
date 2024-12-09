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

// Quarto is a Dagger module that provides functions for running Quarto
type Quarto struct{}

// Cli runs the quarto cli
func (m *Quarto) Cli(
	// commands to run
	command []string,
) *dagger.Container {
	return base().
		WithExec(command)
}

// base returns the quarto container
func base() *dagger.Container {
	return dag.Container().
		From("ghcr.io/quarto-dev/quarto:latest")
}
