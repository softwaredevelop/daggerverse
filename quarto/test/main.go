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

	"github.com/sourcegraph/conc/pool"
)

// Quartotest is a Dagger module that provides functions for running Quarto
type Quartotest struct{}

// All runs all tests.
func (m *Quartotest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CliVersion)

	return p.Wait()
}

// CliVersion runs the quarto --version command.
func (m *Quartotest) CliVersion(ctx context.Context) error {

	_, err := dag.Quarto().Cli([]string{"quarto", "--version"}).Stderr(ctx)

	if err != nil {
		return err
	}

	return nil
}
