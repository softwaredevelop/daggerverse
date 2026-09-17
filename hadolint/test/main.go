// Package main provides tests for the Hadolint Dagger module.
package main

import (
	"context"
	"errors"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Hadolinttest provides test functions for the hadolint module.
type Hadolinttest struct{}

// All runs all tests concurrently.
func (m *Hadolinttest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckWithConfig)
	p.Go(m.Check)

	return p.Wait()
}

// CheckWithConfig tests hadolint with an explicit configuration file.
func (m *Hadolinttest) CheckWithConfig(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")
	file := dag.CurrentModule().Source().File("./testdata/.config/.hadolint.yaml")

	_, err := dag.Hadolint().CheckWithConfig(dir, file).Stderr(ctx)
	if err == nil {
		return errors.New("expected hadolint to fail on invalid Dockerfile, but it succeeded")
	}

	// xargs exits with 123 when the underlying linter fails
	if strings.Contains(err.Error(), "exit code: 123") || strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// Check tests hadolint on invalid testdata without explicit configuration.
func (m *Hadolinttest) Check(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Hadolint().Check(dir).Stderr(ctx)
	if err == nil {
		return errors.New("expected hadolint to fail on invalid Dockerfile, but it succeeded")
	}

	// xargs exits with 123 when the underlying linter fails
	if strings.Contains(err.Error(), "exit code: 123") || strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}
