// Package main provides tests for the Hadolint Dagger module.
package main

import (
	"context"
	"dagger/hadolint/test/internal/dagger"
	"errors"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Hadolinttest provides test functions for the hadolint module.
type Hadolinttest struct{}

// All runs all tests concurrently.
func (m *Hadolinttest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.Check)
	p.Go(m.CheckWithOptionalConfig)
	p.Go(m.CheckWithConfig)

	return p.Wait()
}

// Check tests hadolint on invalid testdata without configuration.
func (m *Hadolinttest) Check(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Hadolint().Check(dir).Stderr(ctx)
	if err == nil {
		return errors.New("expected hadolint to fail on invalid Dockerfile, but it succeeded")
	}

	if strings.Contains(err.Error(), "exit code: 123") || strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// CheckWithOptionalConfig tests the new optional config parameter of the Check function.
func (m *Hadolinttest) CheckWithOptionalConfig(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")
	file := dag.CurrentModule().Source().File("./testdata/.config/.hadolint.yaml")

	// Calling Check with the newly introduced optional argument
	_, err := dag.Hadolint().Check(dir, dagger.HadolintCheckOpts{
		Config: file,
	}).Stderr(ctx)

	if err == nil {
		return errors.New("expected hadolint to fail on invalid Dockerfile, but it succeeded")
	}

	if strings.Contains(err.Error(), "exit code: 123") || strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// CheckWithConfig tests the backwards-compatible CheckWithConfig wrapper function.
func (m *Hadolinttest) CheckWithConfig(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")
	file := dag.CurrentModule().Source().File("./testdata/.config/.hadolint.yaml")

	_, err := dag.Hadolint().CheckWithConfig(dir, file).Stderr(ctx)
	if err == nil {
		return errors.New("expected hadolint to fail on invalid Dockerfile, but it succeeded")
	}

	if strings.Contains(err.Error(), "exit code: 123") || strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}
