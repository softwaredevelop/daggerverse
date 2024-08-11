// A generated module for Editorconfig functions
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
	"dagger/editorconfig/internal/dagger"
)

// Editorconfig is a module for checking editorconfig files.
type Editorconfig struct{}

// Check runs the editorconfig-checker command.
func (m *Editorconfig) Check(
	// source is an optional argument that specifies a directory.
	source *dagger.Directory,
	// excludeDirectoryPattern is an optional argument that specifies a pattern to exclude directories.
	// +default=".git"
	excludeDirectoryPattern string,
) *dagger.Container {
	return base().
		WithMountedDirectory("/tmp", source.WithoutDirectory(excludeDirectoryPattern)).
		WithWorkdir("/tmp").
		WithExec([]string{"/editorconfig-checker"})
}

// base returns a container with the editorconfig-checker binary installed.
func base() *dagger.Container {
	install := dag.
		Container().
		From("golang:alpine").
		WithExec([]string{
			"go", "install",
			"github.com/editorconfig-checker/editorconfig-checker/cmd/editorconfig-checker@latest",
		})

	return dag.Container().
		WithFile("/", install.File("/go/bin/editorconfig-checker"))
}
