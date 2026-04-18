// A generated module for SshManagertest functions
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
	"dagger/ssh-manager/test/internal/dagger"
	"fmt"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

type SshManagertest struct{}

// All runs all SSH key generation tests in parallel.
func (m *SshManagertest) All(ctx context.Context) error {
	p := pool.New().WithErrors().WithContext(ctx)

	p.Go(m.TestGenerateRsa)
	p.Go(m.TestGenerateEd25519)
	p.Go(m.TestGenerateEcdsa)

	return p.Wait()
}

// TestGenerateRsa tests the RSA key generation logic.
func (m *SshManagertest) TestGenerateRsa(ctx context.Context) error {
	host := "rsa-target"
	user := "rsa-user"
	local := "test-machine"

	result := dag.SSHManager().GenerateRsa(host, user, dagger.SSHManagerGenerateRsaOpts{
		LocalHostname: local,
		Bits:          4096,
	})

	return m.verifyKeyResult(ctx, result, "ssh-rsa", host, user, local)
}

// TestGenerateEd25519 tests the Ed25519 key generation logic.
func (m *SshManagertest) TestGenerateEd25519(ctx context.Context) error {
	host := "ed-target"
	user := "ed-user"
	local := "test-machine"

	result := dag.SSHManager().GenerateEd25519(host, user, dagger.SSHManagerGenerateEd25519Opts{
		LocalHostname: local,
	})

	return m.verifyKeyResult(ctx, result, "ssh-ed25519", host, user, local)
}

// TestGenerateEcdsa tests the ECDSA key generation logic.
func (m *SshManagertest) TestGenerateEcdsa(ctx context.Context) error {
	host := "ecdsa-target"
	user := "ecdsa-user"
	local := "test-machine"

	// Testing with 521 bits as requested in previous logic
	result := dag.SSHManager().GenerateEcdsa(host, user, dagger.SSHManagerGenerateEcdsaOpts{
		LocalHostname: local,
		Bits:          521,
	})

	return m.verifyKeyResult(ctx, result, "ecdsa-sha2-nistp521", host, user, local)
}

// verifyKeyResult is a helper function to validate the output files.
func (m *SshManagertest) verifyKeyResult(
	ctx context.Context,
	res *dagger.SSHManagerKeyResult,
	expectedAlgo string,
	host string,
	user string,
	local string,
) error {
	// 1. Validate Public Key
	pub, err := res.PublicKey().Contents(ctx)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}
	expectedComment := user + "@" + host
	if !strings.HasPrefix(pub, expectedAlgo) {
		return fmt.Errorf("public key algorithm mismatch: expected %s, got %s", expectedAlgo, pub[:15])
	}
	if !strings.Contains(pub, expectedComment) {
		return fmt.Errorf("public key comment missing: expected %s", expectedComment)
	}

	// 2. Validate Private Key (Header Check)
	// Even though it's a Secret, we get the plaintext for validation in the test
	priv, err := res.PrivateKey().Plaintext(ctx)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}
	if !strings.Contains(priv, "BEGIN OPENSSH PRIVATE KEY") {
		return fmt.Errorf("invalid private key format: missing OPENSSH header")
	}

	// 3. Validate Config File
	conf, err := res.Config().Contents(ctx)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	checks := []string{
		"Host                 " + host + "." + user,
		"Hostname             " + host,
		"IdentityFile         ~/.ssh/id_",
		"ProxyCommand         nc -X 5 -x 127.0.0.1:9050",
	}

	for _, check := range checks {
		if !strings.Contains(conf, check) {
			return fmt.Errorf("config check failed: missing [%s]", check)
		}
	}

	// 4. Validate Key Filename contains the local hostname
	privatKeyName, err := res.PrivateKey().Name(ctx)
	if err == nil && !strings.Contains(privatKeyName, local) {
		return fmt.Errorf("privat key name mismatch: expected to contain %s, got %s", local, privatKeyName)
	}

	publicKeyName, err := res.PublicKey().Name(ctx)
	if err == nil && !strings.Contains(publicKeyName, local) {
		return fmt.Errorf("public key name mismatch: expected to contain %s, got %s", local, publicKeyName)
	}

	return nil
}
