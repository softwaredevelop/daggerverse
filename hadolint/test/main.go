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
	"regexp"

	"github.com/sourcegraph/conc/pool"
)

// Hadolinttest is a module for checking hadolint files.
type Hadolinttest struct{}

// All runs all tests.
func (m *Hadolinttest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckWithConfig)
	p.Go(m.Check)

	return p.Wait()
}

// CheckWithConfig runs the hadolint-checker command with a configuration file.
func (m *Hadolinttest) CheckWithConfig(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	file := dag.CurrentModule().Source().File("./testdata/.config/.hadolint.yaml")
	_, err := dag.Hadolint().CheckWithConfig(dir, file).Stderr(ctx)

	if err != nil {
		re := regexp.MustCompile("exit code: 123")
		if re.MatchString(err.Error()) {
			return nil
		}
	}

	return err
}

// Check runs the hadolint-checker command.
func (m *Hadolinttest) Check(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Hadolint().Check(dir).Stderr(ctx)

	if err != nil {
		re := regexp.MustCompile("exit code: 123")
		if re.MatchString(err.Error()) {
			return nil
		}
	}

	return err
}
