package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"dagger/ssh-manager/internal/dagger"
	"fmt"
)

func (m *SshManager) GenerateEcdsa(
	ctx context.Context,
	remoteHost string,
	remoteUser string,
	// +optional
	// +default="dagger"
	localHostname string,
	// +optional
	// +default=521
	bits int,
	// +optional
	autoPassphrase bool,
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
		return nil, fmt.Errorf("invalid bits for ecdsa")
	}
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return m.finalize(ctx, priv, priv.Public(), "ecdsa", remoteHost, remoteUser, localHostname, autoPassphrase, existingConfig)
}
