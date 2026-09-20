// Package main provides test suites for the Ruff Dagger module.
package main

import (
	"context"
	"dagger/ruff/test/internal/dagger"
	"errors"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Rufftest provides test functions for the Ruff module.
type Rufftest struct{}

// All runs all tests concurrently.
func (m *Rufftest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.Check)
	p.Go(m.CheckWithConfig)
	p.Go(m.CheckWithOptionalConfig)
	p.Go(m.FormatCheck)

	return p.Wait()
}

// Check tests that ruff check properly fails on invalid Python files.
func (m *Rufftest) Check(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Ruff().Check(dir).Stderr(ctx)
	if err == nil {
		return errors.New("expected ruff check to fail on invalid python testdata, but it succeeded")
	}

	// Ruff exits with code 1 when lint violations are found
	if strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// CheckWithOptionalConfig tests calling Check with the new optional config argument.
func (m *Rufftest) CheckWithOptionalConfig(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")
	file := dag.CurrentModule().Source().File("./testdata/.config/.ruff.toml")

	_, err := dag.Ruff().Check(dir, dagger.RuffCheckOpts{
		Config: file,
	}).Stderr(ctx)
	if err == nil {
		return errors.New("expected ruff check with optional config to fail, but it succeeded")
	}

	if strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// CheckWithConfig tests the backwards-compatible CheckWithConfig wrapper function.
func (m *Rufftest) CheckWithConfig(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")
	file := dag.CurrentModule().Source().File("./testdata/.config/.ruff.toml")

	_, err := dag.Ruff().CheckWithConfig(dir, file).Stderr(ctx)
	if err == nil {
		return errors.New("expected ruff check with config to fail, but it succeeded")
	}

	if strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// FormatCheck tests that ruff format --check detects unformatted code in testdata.
func (m *Rufftest) FormatCheck(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Ruff().FormatCheck(dir).Stderr(ctx)
	if err == nil {
		return errors.New("expected ruff format --check to fail on unformatted testdata, but it succeeded")
	}

	// Ruff format exits with code 1 if files would be reformatted
	if strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}
