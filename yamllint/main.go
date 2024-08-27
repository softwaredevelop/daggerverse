// A generated module for Yamllint functions
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
	"dagger/yamllint/internal/dagger"
)

// Yamllint is a module for checking YAML files.
type Yamllint struct{}

// Check runs yamllint on the provided source directory.
func (m *Yamllint) Check(
	// source is an optional argument that specifies a directory.
	source *dagger.Directory,
) *dagger.Container {
	return base().
		WithMountedDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithExec([]string{"yamllint",
			"--config-data",
			"{extends: default, rules: {line-length: {level: warning}}}",
			"--no-warnings",
			"."})
}

// base returns a container with the yamllint image and no entrypoint.
func base() *dagger.Container {
	return dag.
		Container().
		From("pipelinecomponents/yamllint").
		WithoutEntrypoint()
}
