// A generated module for Editorconfigtest functions
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

	"github.com/sourcegraph/conc/pool"
)

// Editorconfigtest is a module for checking editorconfig files.
type Editorconfigtest struct{}

// All runs all tests.
func (m *Editorconfigtest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.CheckExcludeDirectory)
	p.Go(m.CheckIncludeDirectory)

	return p.Wait()
}

// CheckIncludeDirectory runs the editorconfig-checker command.
func (m *Editorconfigtest) CheckIncludeDirectory(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Editorconfig().Check(dir, "exclude_directory").Stdout(ctx)

	if err != nil {
		return nil
	}

	return err
}

// CheckExcludeDirectory runs the editorconfig-checker command.
func (m *Editorconfigtest) CheckExcludeDirectory(ctx context.Context) error {

	dir := dag.CurrentModule().Source().Directory("./testdata")
	_, err := dag.Editorconfig().Check(dir, ".testdata").Stdout(ctx)

	if err != nil {
		return err
	}

	return nil
}
