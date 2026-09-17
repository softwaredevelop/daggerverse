// Package main provides tests for the Actionlint Dagger module.
package main

import (
	"context"
	"errors"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Actionlinttest is a module for testing the actionlint module.
type Actionlinttest struct{}

// All runs all tests concurrently.
func (m *Actionlinttest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckWorkflowFailsOnInvalid)

	return p.Wait()
}

// CheckWorkflowFailsOnInvalid tests that actionlint detects invalid workflow files.
func (m *Actionlinttest) CheckWorkflowFailsOnInvalid(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Actionlint().Check(dir).Sync(ctx)
	if err == nil {
		return errors.New("expected actionlint to fail on invalid workflow, but it succeeded")
	}

	// Exit code 1 means actionlint successfully parsed the workflow and found lint/syntax errors!
	if strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	// If it fails with exit code 3 or anything else, fail the test
	return err
}
