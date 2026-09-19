// Package main provides test suites for the Revive Dagger module.
package main

import (
	"context"
	"dagger/revive/test/internal/dagger"
	"errors"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Revivetest provides test functions for the Revive module.
type Revivetest struct{}

// All runs all tests concurrently.
func (m *Revivetest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckFailsByDefault)
	p.Go(m.CheckWithConfig)

	return p.Wait()
}

// CheckFailsByDefault asserts that Revive fails on invalid code when using default rules.
func (m *Revivetest) CheckFailsByDefault(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Revive().Check(dir).Sync(ctx)
	if err == nil {
		return errors.New("expected revive to fail on invalid Go code by default, but it succeeded")
	}

	// Must fail with exit code 1 due to lint violations
	if strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// CheckWithConfig asserts that providing a custom revive.toml overrides rules and passes.
func (m *Revivetest) CheckWithConfig(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")
	// Read the renamed custom config file
	configFile := dag.CurrentModule().Source().File("./testdata/custom-revive.toml")

	// Pass the custom configuration which disables the failing rules
	_, err := dag.Revive().Check(dir, dagger.ReviveCheckOpts{
		Config: configFile,
	}).Sync(ctx)

	// Here we expect success (err == nil) because the custom config disabled the failing rules
	if err != nil {
		return errors.New("expected revive to pass with custom config, but it failed: " + err.Error())
	}

	return nil
}
