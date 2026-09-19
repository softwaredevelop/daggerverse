// Package main provides test suites for the Quarto Dagger module.
package main

import (
	"context"
	"dagger/quarto/test/internal/dagger"
	"errors"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Quartotest provides tests for the Quarto module.
type Quartotest struct{}

// All runs all tests concurrently.
func (m *Quartotest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.Extensions)
	p.Go(m.Build)
	p.Go(m.Render)
	p.Go(m.FullVersion)
	p.Go(m.Version)

	return p.Wait()
}

// Extensions tests installing Quarto extensions.
func (m *Quartotest) Extensions(ctx context.Context) error {
	_, err := dag.Quarto(
		dagger.QuartoOpts{
			Extensions: []string{
				"quarto-ext/latex-environment",
				"quarto-ext/include-code-files",
			},
		},
	).Cli("quarto --version").Sync(ctx)

	return err
}

// Build tests compiling a document to PDF and verifies the output artifact.
func (m *Quartotest) Build(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	outputDir, err := dag.Quarto(
		dagger.QuartoOpts{
			Image: "ghcr.io/quarto-dev/quarto-full",
		},
	).Build(dir).Sync(ctx)
	if err != nil {
		return err
	}

	// Verify generated files directly inside the Dagger Directory (no host disk writes needed)
	entries, err := outputDir.Entries(ctx)
	if err != nil {
		return err
	}

	for _, name := range entries {
		if strings.HasSuffix(name, ".pdf") {
			return nil // PDF successfully generated!
		}
	}

	return errors.New("quarto build succeeded, but no .pdf file was found in output directory")
}

// Render tests running the render process directly on the container.
func (m *Quartotest) Render(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Quarto(
		dagger.QuartoOpts{
			Image: "ghcr.io/quarto-dev/quarto-full",
		},
	).Render(dir).Sync(ctx)

	return err
}

// FullVersion tests the CLI against the quarto-full image.
func (m *Quartotest) FullVersion(ctx context.Context) error {
	_, err := dag.Quarto(
		dagger.QuartoOpts{
			Image: "ghcr.io/quarto-dev/quarto-full",
		},
	).Cli("quarto --version").Sync(ctx)

	return err
}

// Version tests the CLI against the default lightweight image.
func (m *Quartotest) Version(ctx context.Context) error {
	_, err := dag.Quarto().Cli("quarto --version").Sync(ctx)

	return err
}
