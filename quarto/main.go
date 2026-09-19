// A Dagger module for authoring, rendering, and publishing Quarto projects.
//
// This module provides tools to render Quarto documents, books, websites,
// install custom Quarto extensions, and manage LaTeX packages via TeX Live.
package main

import (
	"dagger/quarto/internal/dagger"
	"strings"
)

const (
	defaultImageRepository = "ghcr.io/quarto-dev/quarto:latest"
	tlmgrUpdateURL         = "https://mirror.ctan.org/systems/texlive/tlnet/update-tlmgr-latest.sh"
)

// Quarto provides functions for running and compiling Quarto projects.
type Quarto struct {
	// +private
	Image string
	// +private
	Extensions []string
	// +private
	LatexPackages []string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Quarto module.
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
	// List of Quarto extensions to install (e.g. "quarto-ext/lightbox").
	// +optional
	extensions []string,
	// List of optional LaTeX packages to install via tlmgr (requires a TeX-enabled image like quarto-full).
	// +optional
	latexPackages []string,
) *Quarto {
	return &Quarto{
		Image:         image,
		Extensions:    extensions,
		LatexPackages: latexPackages,
	}
}

// container returns the underlying Dagger container, lazily initialized.
func (m *Quarto) container() *dagger.Container {
	if m.Ctr != nil {
		return m.Ctr
	}

	image := m.Image
	if image == "" {
		image = defaultImageRepository
	}

	ctr := dag.Container().From(image)

	// Only update tlmgr and install packages if LaTeX packages are explicitly requested
	if len(m.LatexPackages) > 0 {
		ctr = ctr.WithExec([]string{
			"sh", "-c",
			"curl -fsSL " + tlmgrUpdateURL + " -o update-tlmgr-latest.sh && sh update-tlmgr-latest.sh -- --update",
		})

		// Install all packages in a single command layer
		installCmd := append([]string{"tlmgr", "install"}, m.LatexPackages...)
		ctr = ctr.WithExec(installCmd)
	}

	// Install Quarto extensions
	for _, ext := range m.Extensions {
		ctr = ctr.WithExec([]string{"quarto", "add", "--no-prompt", ext})
	}

	m.Ctr = ctr
	return m.Ctr
}

// Build compiles the Quarto project and returns the exported output directory.
func (m *Quarto) Build(
	// Source directory of the Quarto project.
	source *dagger.Directory,
	// Output directory name relative to source (e.g. "_output", "_site", "_book").
	// +default="_output"
	// +optional
	outputDir string,
) *dagger.Directory {
	workdir := "/work"
	outPath := workdir + "/" + strings.TrimPrefix(outputDir, "/")

	return m.container().
		WithMountedDirectory(workdir, source).
		WithWorkdir(workdir).
		WithExec([]string{"quarto", "render"}).
		Directory(outPath)
}

// Render runs the quarto render command and returns the container.
func (m *Quarto) Render(
	// Source directory of the Quarto project.
	source *dagger.Directory,
) *dagger.Container {
	return m.container().
		WithMountedDirectory("/work", source).
		WithWorkdir("/work").
		WithExec([]string{"quarto", "render"})
}

// Cli executes an arbitrary command in the Quarto container.
func (m *Quarto) Cli(
	// Command string to execute (e.g. "quarto --version" or "quarto check").
	args string,
) *dagger.Container {
	return m.container().
		WithExec([]string{"sh", "-c", args})
}
