// A generated module for Actionlint functions
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
	"dagger/actionlint/internal/dagger"
)

// Actionlint is a module for checking GitHub Actions workflows.
type Actionlint struct{}

// Check runs the actionlint command.
func (m *Actionlint) Check(
	// source is an optional argument that specifies a directory.
	source *dagger.Directory,
) *dagger.Container {
	return base().
		WithMountedDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithExec([]string{"sh", "-c", "find . -type f -name '*.yml' -print0 | xargs -0 actionlint"})
}

// base returns a container with the actionlint binary installed.
func base() *dagger.Container {
	install := dag.Container().
		From("golang:alpine").
		WithExec([]string{
			"go", "install",
			"github.com/rhysd/actionlint/cmd/actionlint@latest",
		})

	shellcheck := dag.Container().
		From("koalaman/shellcheck-alpine:stable")

	return dag.Container().
		From("cgr.dev/chainguard/wolfi-base:latest").
		WithFile("/usr/bin/actionlint", install.File("/go/bin/actionlint")).
		WithFile("/usr/bin/shellcheck", shellcheck.File("/bin/shellcheck"))
}
