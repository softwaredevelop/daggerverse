package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"dagger/ssh-manager/internal/dagger"
)

// GenerateEd25519 creates a new Ed25519 SSH key pair. This is the recommended modern algorithm.
func (m *SshManager) GenerateEd25519(
	ctx context.Context,
	// The remote host address (e.g., github.com).
	remoteHost string,
	// The user name on the remote host.
	remoteUser string,
	// The hostname of the local machine used for key naming.
	// +optional
	// +default="dagger"
	localHostname string,
	// Automatically generate a 32-character passphrase for the private key.
	// +optional
	autoPassphrase bool,
	// An existing SSH config file to append the new configuration to.
	// +optional
	existingConfig *dagger.File,
) (*KeyResult, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return m.finalize(ctx, priv, pub, "ed25519", remoteHost, remoteUser, localHostname, autoPassphrase, existingConfig)
}
