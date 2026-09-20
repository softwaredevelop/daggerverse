// A Dagger module for linting and formatting Python code and Jupyter Notebooks using Ruff.
//
// Ruff is an extremely fast Python linter and code formatter written in Rust.
// This module provides functions to run lint checks and formatting validations.
package main

import (
	"dagger/ruff/internal/dagger"
)

const (
	defaultImageRepository = "ghcr.io/astral-sh/ruff:latest"
)

// Ruff provides functions for running the Ruff linter and formatter.
type Ruff struct {
	// +private
	Image string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Ruff struct.
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Ruff {
	return &Ruff{
		Image: image,
	}
}

// container returns the underlying Dagger container, lazily initialized.
func (m *Ruff) container() *dagger.Container {
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

// Check runs ruff check on the target directory.
// It supports standalone checks or using a custom ruff.toml / pyproject.toml configuration.
func (m *Ruff) Check(
	// Source directory containing Python files or Jupyter notebooks.
	source *dagger.Directory,
	// Optional configuration file (ruff.toml or pyproject.toml).
	// +optional
	config *dagger.File,
) *dagger.Container {
	ctr := m.container().
		WithMountedDirectory("/work", source).
		WithWorkdir("/work")

	args := []string{"/ruff", "check"}

	if config != nil {
		ctr = ctr.WithFile("/etc/ruff.toml", config)
		args = append(args, "--config", "/etc/ruff.toml")
	}

	return ctr.WithExec(args)
}

// CheckWithConfig runs ruff check with a configuration file.
// Maintained for explicit backwards compatibility.
func (m *Ruff) CheckWithConfig(
	// Source directory containing Python files or Jupyter notebooks.
	source *dagger.Directory,
	// Configuration file for ruff.
	file *dagger.File,
) *dagger.Container {
	return m.Check(source, file)
}

// FormatCheck runs ruff format --check to verify code formatting without modifying files.
func (m *Ruff) FormatCheck(
	// Source directory containing Python files or Jupyter notebooks.
	source *dagger.Directory,
	// Optional configuration file (ruff.toml or pyproject.toml).
	// +optional
	config *dagger.File,
) *dagger.Container {
	ctr := m.container().
		WithMountedDirectory("/work", source).
		WithWorkdir("/work")

	args := []string{"/ruff", "format", "--check"}

	if config != nil {
		ctr = ctr.WithFile("/etc/ruff.toml", config)
		args = append(args, "--config", "/etc/ruff.toml")
	}

	return ctr.WithExec(args)
}
