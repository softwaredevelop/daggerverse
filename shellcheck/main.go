// A Dagger module for linting shell scripts using ShellCheck.
//
// ShellCheck is a static analysis tool for shell scripts (sh, bash, dash, ksh)
// that finds bugs, syntax issues, and style warnings.
package main

import (
	"dagger/shellcheck/internal/dagger"
)

const (
	defaultImageRepository = "koalaman/shellcheck-alpine:latest"
)

// Shellcheck provides functions for checking shell scripts.
type Shellcheck struct {
	// +private
	Image string
	// +private
	Ctr *dagger.Container
}

// New creates a new instance of the Shellcheck struct.
func New(
	// Custom image reference in "repository:tag" format to use as a base container.
	// +optional
	image string,
) *Shellcheck {
	return &Shellcheck{
		Image: image,
	}
}

// container returns the underlying Dagger container, lazily initialized.
func (m *Shellcheck) container() *dagger.Container {
	if m.Ctr != nil {
		return m.Ctr
	}

	image := m.Image
	if image == "" {
		image = defaultImageRepository
	}

	m.Ctr = dag.Container().From(image)
	return m.Ctr
}

// Check runs shellcheck on all .sh and .bash files in the source directory.
// If an optional config file is provided, it is placed in the working directory as .shellcheckrc.
func (m *Shellcheck) Check(
	// Source directory containing shell scripts.
	source *dagger.Directory,
	// Minimum severity of errors to report (error, warning, info, style).
	// +optional
	severity string,
	// Optional shellcheck configuration file (.shellcheckrc).
	// +optional
	config *dagger.File,
) *dagger.Container {
	ctr := m.container().
		WithMountedDirectory("/work", source).
		WithWorkdir("/work")

	if config != nil {
		ctr = ctr.WithFile("/work/.shellcheckrc", config).
			WithFile("/root/.shellcheckrc", config)
	}

	severityArg := ""
	if severity != "" {
		severityArg = "--severity=" + severity + " "
	}

	cmd := "find . -type f \\( -name '*.sh' -o -name '*.bash' \\) -print0 | xargs -0 -r shellcheck " + severityArg

	return ctr.WithExec([]string{"sh", "-c", cmd})
}

// CheckWithConfig runs shellcheck explicitly using the provided configuration file via --rcfile.
// The configuration file is mandatory for this function.
func (m *Shellcheck) CheckWithConfig(
	// Source directory containing shell scripts.
	source *dagger.Directory,
	// Mandatory configuration file for shellcheck.
	file *dagger.File,
	// Minimum severity of errors to report (error, warning, info, style).
	// +optional
	severity string,
) *dagger.Container {
	ctr := m.container().
		WithMountedDirectory("/work", source).
		WithWorkdir("/work").
		WithFile("/etc/shellcheckrc", file)

	severityArg := ""
	if severity != "" {
		severityArg = "--severity=" + severity + " "
	}

	// Explicitly pass --rcfile to force using the provided configuration file
	cmd := "find . -type f \\( -name '*.sh' -o -name '*.bash' \\) -print0 | xargs -0 -r shellcheck --rcfile=/etc/shellcheckrc " + severityArg

	return ctr.WithExec([]string{"sh", "-c", cmd})
}
