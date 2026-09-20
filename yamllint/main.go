// A Dagger module for linting YAML files using yamllint.
//
// yamllint does not only check for syntax validity, but also for weirdnesses
// like key repetition and cosmetic problems such as lines length, trailing spaces, indentation, etc.
package main

import (
	"dagger/yamllint/internal/dagger"
)

const (
	defaultImageRepository = "pipelinecomponents/yamllint:latest"
	defaultConfigData      = "{extends: default, rules: {line-length: {level: warning}}}"
)

// Yamllint provides functions for checking YAML files.
type Yamllint struct {
	// +private
	Image string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Yamllint struct.
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Yamllint {
	return &Yamllint{
		Image: image,
	}
}

// container returns the underlying Dagger container, lazily initialized.
func (m *Yamllint) container() *dagger.Container {
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

// Check runs yamllint on the source directory.
// By default, it ignores line-length warnings. An optional configuration file can be provided.
func (m *Yamllint) Check(
	// Source directory containing YAML files.
	source *dagger.Directory,
	// Optional configuration file (.yamllint).
	// +optional
	config *dagger.File,
) *dagger.Container {
	ctr := m.container().
		WithMountedDirectory("/work", source).
		WithWorkdir("/work")

	if config != nil {
		return ctr.
			WithFile("/etc/yamllint", config).
			WithExec([]string{"yamllint", "-c", "/etc/yamllint", "."})
	}

	return ctr.WithExec([]string{
		"yamllint",
		"--config-data",
		defaultConfigData,
		"--no-warnings",
		".",
	})
}

// CheckWithConfig runs yamllint using an explicitly provided configuration file.
// The configuration file is mandatory for this function.
func (m *Yamllint) CheckWithConfig(
	// Source directory containing YAML files.
	source *dagger.Directory,
	// Mandatory configuration file for yamllint.
	file *dagger.File,
) *dagger.Container {
	return m.container().
		WithMountedDirectory("/work", source).
		WithWorkdir("/work").
		WithFile("/etc/yamllint", file).
		WithExec([]string{"yamllint", "-c", "/etc/yamllint", "."})
}
