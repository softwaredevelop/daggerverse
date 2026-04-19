package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"dagger/ssh-manager/internal/dagger"
)

// GenerateRsa creates a new RSA SSH key pair.
func (m *SshManager) GenerateRsa(
	ctx context.Context,
	// The remote host address (e.g., github.com).
	remoteHost string,
	// The user name on the remote host.
	remoteUser string,
	// The hostname of the local machine used for key naming.
	// +optional
	// +default="dagger"
	localHostname string,
	// The bit length for the RSA key.
	// +optional
	// +default=4096
	bits int,
	// Automatically generate a 32-character passphrase for the private key.
	// +optional
	autoPassphrase bool,
	// An existing SSH config file to append the new configuration to.
	// +optional
	existingConfig *dagger.File,
) (*KeyResult, error) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return m.finalize(ctx, priv, priv.Public(), "rsa", remoteHost, remoteUser, localHostname, autoPassphrase, existingConfig)
}
