// A generated module for Yamllinttest functions
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

// Yamllinttest is a module for checking YAML files.
type Yamllinttest struct{}

// All runs all tests.
func (m *Yamllinttest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.Check)

	return p.Wait()
}

// Check runs the revive command.
func (m *Yamllinttest) Check(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Yamllint().Check(dir).Stdout(ctx)

	if err != nil {
		expectedErrors := []string{
			"syntax error: expected <block end>, but found '<block mapping start>' (syntax)",
			"too many spaces before colon  (colons)",
		}
		for _, expectedError := range expectedErrors {
			if !strings.Contains(err.Error(), expectedError) {
				return err
			}
		}
	}

	return nil
}
