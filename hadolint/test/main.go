// A generated module for Hadolinttest functions
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

// Hadolinttest is a module for checking hadolint files.
type Hadolinttest struct{}

// All runs all tests.
func (m *Hadolinttest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckWithoutConfig)

	return p.Wait()
}

// CheckWithoutConfig runs the hadolint-checker command.
func (m *Hadolinttest) CheckWithoutConfig(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Hadolint().CheckWithoutConfig(dir).Stdout(ctx)

	if err != nil {
		if !strings.Contains(err.Error(), "DL3007") {
			return err
		}
	}

	return nil
}
