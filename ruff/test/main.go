// A generated module for Rufftest functions
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
	"strings"

	"github.com/sourcegraph/conc/pool"
)

type Rufftest struct{}

// All runs all tests.
func (m *Rufftest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.Check)
	p.Go(m.CheckWithConfig)

	return p.Wait()
}

// CheckWithConfig runs the ruff check command with a configuration file.
func (m *Rufftest) CheckWithConfig(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	file := dag.CurrentModule().Source().File("./testdata/.config/.ruff.toml")
	_, err := dag.Ruff().CheckWithConfig(dir, file).Stdout(ctx)

	if err != nil {
		errorIDs := []string{"F821"}
		for _, id := range errorIDs {
			if !strings.Contains(err.Error(), id) {
				return err
			}
		}
	}

	return nil
}

// Check runs the ruff check command.
func (m *Rufftest) Check(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Ruff().Check(dir).Stdout(ctx)

	if err != nil {
		errorIDs := []string{"F401"}
		for _, id := range errorIDs {
			if !strings.Contains(err.Error(), id) {
				return err
			}
		}
	}

	return nil
}
