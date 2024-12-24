// A generated module for Test functions
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
	"dagger/hello/test/internal/dagger"

	"github.com/sourcegraph/conc/pool"
)

// Test is a module for running tests.
type Test struct{}

// All runs all tests.
func (m *Test) All(ctx context.Context) error {
	p := pool.New().WithContext(ctx)

	p.Go(m.HelloContainer)
	p.Go(m.HelloStringContainer)
	p.Go(m.HelloString)

	return p.Wait()
}

// HelloContainer tests the HelloContainer function.
func (m *Test) HelloContainer(ctx context.Context) error {

	argString := "Hello, Daggerverse!"
	_, err := dag.Hello().HelloContainer(argString).Sync(ctx)

	if err != nil {
		return err
	}

	return nil

}

// HelloStringContainer tests the HelloString function.
func (m *Test) HelloStringContainer(ctx context.Context) error {

	argString := "Hello, Daggerverse!"
	_, err := dag.Hello(
		dagger.HelloOpts{
			Image: "cgr.dev/chainguard/wolfi-base:latest",
		},
	).HelloString(argString).Sync(ctx)

	if err != nil {
		return err
	}

	return nil

}

// HelloString tests the HelloString function.
func (m *Test) HelloString(ctx context.Context) error {

	argString := "Hello, Daggerverse!"
	_, err := dag.Hello().HelloString(argString).Sync(ctx)

	if err != nil {
		return err
	}

	return nil

}
