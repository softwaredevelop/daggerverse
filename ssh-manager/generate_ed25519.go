package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"dagger/ssh-manager/internal/dagger"
)

func (m *SshManager) GenerateEd25519(
	ctx context.Context,
	remoteHost string,
	remoteUser string,
	// +optional
	// +default="dagger"
	localHostname string,
	// +optional
	autoPassphrase bool,
	// +optional
	existingConfig *dagger.File,
) (*KeyResult, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return m.finalize(ctx, priv, pub, "ed25519", remoteHost, remoteUser, localHostname, autoPassphrase, existingConfig)
}
