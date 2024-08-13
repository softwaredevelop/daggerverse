// A generated module for Actionlinttes functions
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

// Actionlinttes is a module for testing actionlint.
type Actionlinttes struct{}

// All runs all tests.
func (m *Actionlinttes) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckWorkflow)

	return p.Wait()
}

// CheckWorkflow runs a test on a directory.
func (m *Actionlinttes) CheckWorkflow(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Actionlint().Check(dir).Stderr(ctx)

	if err != nil {
		if !strings.Contains(err.Error(), "exit code") {
			return err
		}
	}

	return nil
}