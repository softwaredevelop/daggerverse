package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"dagger/ssh-manager/internal/dagger"
)

func (m *SshManager) GenerateRsa(
	ctx context.Context,
	remoteHost string,
	remoteUser string,
	// +optional
	// +default="dagger"
	localHostname string,
	// +optional
	// +default=4096
	bits int,
	// +optional
	autoPassphrase bool,
	// +optional
	existingConfig *dagger.File,
) (*KeyResult, error) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return m.finalize(ctx, priv, priv.Public(), "rsa", remoteHost, remoteUser, localHostname, autoPassphrase, existingConfig)
}
