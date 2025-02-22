// A generated module for Revivetest functions
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
	"regexp"

	"github.com/sourcegraph/conc/pool"
)

// Revivetest is a Dagger module that provides functions for running Revive linter
type Revivetest struct{}

// All runs all tests.
func (m *Revivetest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.Check)

	return p.Wait()
}

// Check runs the revive command.
func (m *Revivetest) Check(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Revive().Check(dir).Sync(ctx)

	if err != nil {
		re := regexp.MustCompile("exit code: 1")
		if re.MatchString(err.Error()) {
			return nil
		}
	}

	return err
}
