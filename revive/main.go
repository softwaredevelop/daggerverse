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

// Revive is a Dagger module that provides functions for running Revive linter
type Revive struct{}

// Check runs the revive command.
func (m *Revive) Check(
	// source is an optional argument that specifies a directory.
	source *dagger.Directory,
) *dagger.Container {
	return base().
		WithMountedDirectory("/go", source).
		WithWorkdir("/go").
		WithExec([]string{"revive", "-set_exit_status", "./..."})
}

// Test_Revive runs tests for the Revive module
func base() *dagger.Container {
	image := dag.
		Container().
		From("ghcr.io/mgechev/revive:latest")

	return dag.Container().
		From("golang:alpine").
		WithFile("/usr/bin/revive", image.File("/usr/bin/revive"))
}
