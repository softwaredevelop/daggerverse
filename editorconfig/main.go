// A Dagger module to validate files against an .editorconfig specification.
//
// This module wraps the editorconfig-checker tool to ensure codebases adhere
// to formatting rules defined in their .editorconfig files. It can be easily
// integrated into local developer workflows or CI/CD pipelines (such as GitHub Actions).
package main

import (
	"dagger/editorconfig/internal/dagger"
)

const (
	defaultImageRepository = "mstruebing/editorconfig-checker:latest"
)

// Editorconfig provides functions to run editorconfig checks.
type Editorconfig struct {
	// +private
	Image string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Editorconfig module.
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Editorconfig {
	return &Editorconfig{
		Image: image,
	}
}

// container returns the underlying Dagger container, lazily initialized.
func (m *Editorconfig) container() *dagger.Container {
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

// Check runs the editorconfig-checker command on the target directory.
// It returns the container allowing callers to chain methods such as stderr, stdout, or sync.
func (m *Editorconfig) Check(
	// Source directory to validate.
	source *dagger.Directory,
	// A relative directory path to exclude from the check (defaults to ".git").
	// +default=".git"
	// +optional
	excludeDirectoryPattern string,
) *dagger.Container {

	src := source
	// Only exclude the path if a non-empty string is provided to prevent runtime errors
	if excludeDirectoryPattern != "" {
		src = src.WithoutDirectory(excludeDirectoryPattern)
	}

	return m.container().
		WithMountedDirectory("/work", src).
		WithWorkdir("/work").
		WithExec([]string{"editorconfig-checker"})
}
