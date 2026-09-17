// Package main provides test suites for the Editorconfig Dagger module.
package main

import (
	"context"
	"dagger/editorconfig/test/internal/dagger"
	"errors"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Editorconfigtest provides test functions for the editorconfig module.
type Editorconfigtest struct{}

// All runs all test cases concurrently.
func (m *Editorconfigtest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckFailsOnInvalidFiles)
	p.Go(m.CheckExcludeDirectory)

	return p.Wait()
}

// CheckExcludeDirectory tests that excluding a directory prevents failures.
func (m *Editorconfigtest) CheckExcludeDirectory(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Editorconfig().Check(dir, dagger.EditorconfigCheckOpts{
		ExcludeDirectoryPattern: ".testdata",
	}).Sync(ctx)

	return err
}

// CheckFailsOnInvalidFiles tests that editorconfig-checker properly reports formatting errors.
func (m *Editorconfigtest) CheckFailsOnInvalidFiles(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Editorconfig().Check(dir).Sync(ctx)
	if err == nil {
		// If there is no error, the test MUST fail because testdata is intentionally invalid!
		return errors.New("expected editorconfig-checker to fail on invalid files, but it succeeded")
	}

	// Verify that it failed due to linting violations (exit code 1) and not a container/runtime failure
	if !strings.Contains(err.Error(), "exit code: 1") {
		return err
	}

	return nil
}
