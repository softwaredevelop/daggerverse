// A Dagger module for linting YAML files using yamllint.
//
// yamllint does not only check for syntax validity, but also for weirdnesses
// like key repetition and cosmetic problems such as lines length, trailing spaces, indentation, etc.
package main

import (
	"dagger/yamllint/internal/dagger"
)

const (
	defaultConfigData = "{extends: default, rules: {line-length: {level: warning}}}"
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
	// If not provided, a clean container built from official Alpine Linux with yamllint will be used.
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

	// If the user provided a custom image, use it directly
	if m.Image != "" {
		m.Ctr = dag.Container().From(m.Image)
		return m.Ctr
	}

	// Build a minimal, secure container on-the-fly from official Alpine Linux
	m.Ctr = dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "yamllint"})

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
