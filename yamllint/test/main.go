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
	"regexp"

	"github.com/sourcegraph/conc/pool"
)

// Yamllinttest is a module for checking YAML files.
type Yamllinttest struct{}

// All runs all tests.
func (m *Yamllinttest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckWithConfig)
	p.Go(m.Check)

	return p.Wait()
}

// CheckWithConfig runs the yamllint command with a configuration file.
func (m *Yamllinttest) CheckWithConfig(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	file := dag.CurrentModule().Source().File("./testdata/.config/.yamllint")
	_, err := dag.Yamllint().CheckWithConfig(dir, file).Stderr(ctx)

	if err != nil {
		re := regexp.MustCompile("exit code: 123")
		if re.MatchString(err.Error()) {
			return nil
		}
	}

	return err
}

// Check runs the revive command.
func (m *Yamllinttest) Check(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Yamllint().Check(dir).Sync(ctx)

	if err != nil {
		re := regexp.MustCompile("exit code: 1")
		if re.MatchString(err.Error()) {
			return nil
		}
	}

	return err
}
