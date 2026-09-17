// A Dagger module to lint GitHub Actions workflow files using actionlint.
//
// This module validates GitHub Actions workflow files for syntax errors,
// invalid expressions, and schema violations.
package main

import (
	"dagger/actionlint/internal/dagger"
)

const (
	defaultImageRepository = "rhysd/actionlint:latest"
)

// Actionlint provides functions for checking GitHub Actions workflows.
type Actionlint struct {
	// +private
	Image string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Actionlint struct.
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Actionlint {
	return &Actionlint{
		Image: image,
	}
}

// container returns the underlying Dagger container, lazily initialized.
func (m *Actionlint) container() *dagger.Container {
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

// Check runs actionlint on workflow files.
// It supports directories mounted at repository root or directly inside .github/workflows.
func (m *Actionlint) Check(
	// Source directory containing workflows or repository root.
	source *dagger.Directory,
) *dagger.Container {

	// Shell script that:
	// 1. If .github/workflows exists, runs plain actionlint (standard repo root mode).
	// 2. Otherwise searches for both .yml and .yaml files in the current folder and passes them to actionlint.
	cmd := `find . -type f \( -name '*.yml' -o -name '*.yaml' \) -print0 | xargs -0 -r actionlint`

	return m.container().
		WithMountedDirectory("/work", source).
		WithWorkdir("/work").
		WithExec([]string{"sh", "-c", cmd})
}
