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

const (
	defaultImageRepository = "pipelinecomponents/yamllint:latest"
)

// Yamllint is a module for checking YAML files.
type Yamllint struct {
	// +private
	Image string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Yamllint struct
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Yamllint {
	return &Yamllint{
		Image: image,
	}
}

// Container returns the underlying Dagger container
func (m *Yamllint) Container() *dagger.Container {
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

// CheckWithConfig runs the yamllint command with a configuration file.
func (m *Yamllint) CheckWithConfig(
	// source is an optional argument that specifies a directory.
	source *dagger.Directory,
	// file is an optional argument that specifies yamllint configuration file.
	file *dagger.File,
) *dagger.Container {
	return m.Container().
		WithMountedDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithFile("/.config/.yamllint", file).
		WithExec([]string{"sh", "-c", "find . -type f \\( -name '*.yaml' -o -name '*.yml' \\) -print0 | xargs -0 yamllint -c /.config/.yamllint"})
}

// Check runs yamllint on the provided source directory.
func (m *Yamllint) Check(
	// source is an optional argument that specifies a directory.
	source *dagger.Directory,
) *dagger.Container {
	return m.Container().
		WithMountedDirectory("/tmp", source).
		WithWorkdir("/tmp").
		WithExec([]string{"yamllint",
			"--config-data",
			"{extends: default, rules: {line-length: {level: warning}}}",
			"--no-warnings",
			"."})
}
