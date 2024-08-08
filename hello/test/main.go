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
	"fmt"

	"github.com/sourcegraph/conc/pool"
)

// Test is a module for running tests.
type Test struct{}

// All runs all tests.
func (m *Test) All(ctx context.Context) error {
	p := pool.New().WithContext(ctx)

	p.Go(m.HelloString)
	p.Go(m.HelloContainer)

	return p.Wait()
}

// HelloContainer tests the HelloContainer function.
func (m *Test) HelloContainer(ctx context.Context) error {
	const expected = "Hello, Daggerverse!\n"

	actual, err := dag.Hello().HelloContainer(ctx)

	if err != nil {
		return err
	}

	if actual != expected {
		return fmt.Errorf("expected %q, got %q", expected, actual)
	}

	return nil
}

// HelloString tests the HelloString function.
func (m *Test) HelloString(ctx context.Context) error {
	const expected = "Hello, Daggerverse!"

	actual, err := dag.Hello().HelloString(ctx)

	if err != nil {
		return err
	}

	if actual != expected {
		return fmt.Errorf("expected %q, got %q", expected, actual)
	}

	return nil
}
