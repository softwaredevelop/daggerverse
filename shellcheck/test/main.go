// Package main provides test suites for the Shellcheck Dagger module.
package main

import (
	"context"
	"dagger/shellcheck/test/internal/dagger"
	"errors"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Shellchecktest provides test functions for the Shellcheck module.
type Shellchecktest struct{}

// All runs all tests concurrently.
func (m *Shellchecktest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckDirectoryFailsOnInvalid)
	p.Go(m.CheckBashFiles)
	p.Go(m.CheckWithOptionalConfigPasses)
	p.Go(m.CheckWithConfigPasses)

	return p.Wait()
}

// CheckDirectoryFailsOnInvalid tests that shellcheck detects errors in standard .sh scripts.
func (m *Shellchecktest) CheckDirectoryFailsOnInvalid(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata")

	_, err := dag.Shellcheck().Check(dir).Stderr(ctx)
	if err == nil {
		return errors.New("expected shellcheck to fail on invalid shell script, but it succeeded")
	}

	if strings.Contains(err.Error(), "exit code: 123") || strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// CheckBashFiles tests that shellcheck properly discovers and checks .bash files.
func (m *Shellchecktest) CheckBashFiles(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata/bash")

	_, err := dag.Shellcheck().Check(dir).Stderr(ctx)
	if err == nil {
		return errors.New("expected shellcheck to detect errors in .bash file, but it succeeded or skipped it")
	}

	if strings.Contains(err.Error(), "exit code: 123") || strings.Contains(err.Error(), "exit code: 1") {
		return nil
	}

	return err
}

// CheckWithOptionalConfigPasses tests that passing an optional config to Check disables warnings and passes.
func (m *Shellchecktest) CheckWithOptionalConfigPasses(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata/config")
	configFile := dag.CurrentModule().Source().File("./testdata/config/custom.shellcheckrc")

	_, err := dag.Shellcheck().Check(dir, dagger.ShellcheckCheckOpts{
		Config: configFile,
	}).Stderr(ctx)

	if err != nil {
		return errors.New("expected shellcheck to pass with optional config, but it failed: " + err.Error())
	}

	return nil
}

// CheckWithConfigPasses tests the dedicated CheckWithConfig function with --rcfile.
func (m *Shellchecktest) CheckWithConfigPasses(ctx context.Context) error {
	dir := dag.CurrentModule().Source().Directory("./testdata/config")
	configFile := dag.CurrentModule().Source().File("./testdata/config/custom.shellcheckrc")

	// Call CheckWithConfig where configFile is mandatory
	_, err := dag.Shellcheck().CheckWithConfig(dir, configFile).Stderr(ctx)

	if err != nil {
		return errors.New("expected CheckWithConfig (--rcfile) to pass, but it failed: " + err.Error())
	}

	return nil
}
