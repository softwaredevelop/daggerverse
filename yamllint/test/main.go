// Package main provides test suites for the Yamllint Dagger module.
package main

import (
	"context"
	"dagger/yamllint/test/internal/dagger"
	"errors"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Yamllinttest provides test functions for the Yamllint module.
type Yamllinttest struct{}

// All runs all tests concurrently.
func (m *Yamllinttest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.Check)
	p.Go(m.CheckWithOptionalConfig)
	p.Go(m.CheckWithConfig)

	return p.Wait()
}

// Check tests that yamllint detects syntax errors in bad.yaml by default.
func (m *Yamllinttest) Check(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Yamllint().Check(dir).Sync(ctx)
	if err == nil {
		return errors.New("expected yamllint to fail on invalid YAML, but it succeeded")
	}

	// yamllint exits with code 1 when syntax or lint errors are found
	if strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// CheckWithOptionalConfig tests the Check function with the optional config parameter.
func (m *Yamllinttest) CheckWithOptionalConfig(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")
	file := dag.CurrentModule().Source().File("./testdata/.config/.yamllint")

	_, err := dag.Yamllint().Check(dir, dagger.YamllintCheckOpts{
		Config: file,
	}).Sync(ctx)

	if err == nil {
		return errors.New("expected yamllint to fail on invalid YAML with optional config, but it succeeded")
	}

	if strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// CheckWithConfig tests the dedicated CheckWithConfig function with a mandatory config file.
func (m *Yamllinttest) CheckWithConfig(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")
	file := dag.CurrentModule().Source().File("./testdata/.config/.yamllint")

	_, err := dag.Yamllint().CheckWithConfig(dir, file).Sync(ctx)
	if err == nil {
		return errors.New("expected CheckWithConfig to fail on invalid YAML, but it succeeded")
	}

	if strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}
