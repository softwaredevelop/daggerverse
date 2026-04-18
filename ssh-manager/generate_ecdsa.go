package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"dagger/ssh-manager/internal/dagger"
	"fmt"
)

// GenerateEcdsa creates a new ECDSA SSH key pair.
func (m *SshManager) GenerateEcdsa(
	ctx context.Context,
	// The remote host address (e.g., github.com).
	remoteHost string,
	// The user name on the remote host.
	remoteUser string,
	// The hostname of the local machine used for key naming.
	// +optional
	// +default="dagger"
	localHostname string,
	// The bit size for the ECDSA curve (256, 384, or 521).
	// +optional
	// +default=521
	bits int,
	// Automatically generate a 32-character passphrase for the private key.
	// +optional
	autoPassphrase bool,
	// An existing SSH config file to append the new configuration to.
	// +optional
	existingConfig *dagger.File,
) (*KeyResult, error) {
	var curve elliptic.Curve
	switch bits {
	case 256:
		curve = elliptic.P256()
	case 384:
		curve = elliptic.P384()
	case 521:
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("invalid bits for ecdsa: must be 256, 384, or 521")
	}
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return m.finalize(ctx, priv, priv.Public(), "ecdsa", remoteHost, remoteUser, localHostname, autoPassphrase, existingConfig)
}
