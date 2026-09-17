// A Dagger module to lint Dockerfiles using hadolint.
//
// This module inspects Dockerfiles for best practices, security issues,
// and style violations using hadolint.
package main

import (
	"dagger/hadolint/internal/dagger"
)

const (
	defaultImageRepository = "hadolint/hadolint:latest-alpine"
)

// Hadolint provides functions for linting Dockerfiles.
type Hadolint struct {
	// +private
	Image string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Hadolint struct.
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Hadolint {
	return &Hadolint{
		Image: image,
	}
}

// container returns the underlying Dagger container, lazily initialized.
func (m *Hadolint) container() *dagger.Container {
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

// Check runs hadolint on Dockerfiles found in the directory.
// It supports Dockerfile, Dockerfile.*, *.dockerfile, and Containerfile naming conventions.
func (m *Hadolint) Check(
	// Source directory containing Dockerfiles.
	source *dagger.Directory,
	// Optional hadolint configuration file (e.g. .hadolint.yaml).
	// +optional
	config *dagger.File,
) *dagger.Container {
	ctr := m.container().
		WithMountedDirectory("/work", source).
		WithWorkdir("/work")

	configArg := ""
	if config != nil {
		ctr = ctr.WithFile("/etc/hadolint.yaml", config)
		configArg = "--config /etc/hadolint.yaml "
	}

	// Finds all standard Dockerfile and Containerfile variants safely
	cmd := "find . -type f \\( -name 'Dockerfile' -o -name 'Dockerfile.*' -o -name '*.dockerfile' -o -name 'Containerfile*' \\) -print0 | xargs -0 -r hadolint " + configArg

	return ctr.WithExec([]string{"sh", "-c", cmd})
}

// CheckWithConfig runs hadolint with a specific configuration file.
// Maintained for explicit backwards compatibility.
func (m *Hadolint) CheckWithConfig(
	// Source directory containing Dockerfiles.
	source *dagger.Directory,
	// Configuration file for hadolint.
	file *dagger.File,
) *dagger.Container {
	return m.Check(source, file)
}
