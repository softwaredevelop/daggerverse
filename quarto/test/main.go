// A generated module for Quartotest functions
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
	"context"
	"dagger/quarto/test/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

// Quartotest is a Dagger module that provides functions for running Quarto
type Quartotest struct{}

// All runs all tests.
func (m *Quartotest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.Render)
	p.Go(m.FullVersion)
	p.Go(m.Version)

	return p.Wait()
}

// Render runs the quarto render command.
func (m *Quartotest) Render(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Quarto(
		dagger.QuartoOpts{
			Image: "ghcr.io/quarto-dev/quarto-full",
		},
	).Render(dir).Stderr(ctx)

	if err != nil {
		return err
	}

	return nil
}

// FullVersion runs the quarto --version command.
func (m *Quartotest) FullVersion(ctx context.Context) error {

	_, err := dag.Quarto(
		dagger.QuartoOpts{
			Image: "ghcr.io/quarto-dev/quarto-full",
		},
	).Cli("quarto --version").Stderr(ctx)

	if err != nil {
		return err
	}

	return nil
}

// CliVersion runs the quarto --version command.
func (m *Quartotest) Version(ctx context.Context) error {

	_, err := dag.Quarto().Cli("quarto --version").Stderr(ctx)

	if err != nil {
		return err
	}

	return nil
}
