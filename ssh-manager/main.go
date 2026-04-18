package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"dagger/ssh-manager/internal/dagger"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type SshManager struct{}

type KeyResult struct {
	PublicKey  *dagger.File      `json:"publicKey"`
	PrivateKey *dagger.Secret    `json:"privateKey"`
	Config     *dagger.File      `json:"config"`
	Files      *dagger.Directory `json:"files"`
}

// func (m *SshManager) GenerateRsa(
// 	ctx context.Context,
// 	remoteHost string,
// 	remoteUser string,
// 	// +optional
// 	// +default="dagger"
// 	localHostname string,
// 	// +optional
// 	// +default=4096
// 	bits int,
// 	// +optional
// 	autoPassphrase bool,
// 	// +optional
// 	existingConfig *dagger.File,
// ) (*KeyResult, error) {
// 	priv, err := rsa.GenerateKey(rand.Reader, bits)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m.finalize(ctx, priv, priv.Public(), "rsa", remoteHost, remoteUser, localHostname, autoPassphrase, existingConfig)
// }

/* func (m *SshManager) GenerateEd25519(
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
} */

// func (m *SshManager) GenerateEcdsa(
// 	ctx context.Context,
// 	remoteHost string,
// 	remoteUser string,
// 	// +optional
// 	// +default="dagger"
// 	localHostname string,
// 	// +optional
// 	// +default=521
// 	bits int,
// 	// +optional
// 	autoPassphrase bool,
// 	// +optional
// 	existingConfig *dagger.File,
// ) (*KeyResult, error) {
// 	var curve elliptic.Curve
// 	switch bits {
// 	case 256:
// 		curve = elliptic.P256()
// 	case 384:
// 		curve = elliptic.P384()
// 	case 521:
// 		curve = elliptic.P521()
// 	default:
// 		return nil, fmt.Errorf("invalid bits for ecdsa")
// 	}
// 	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m.finalize(ctx, priv, priv.Public(), "ecdsa", remoteHost, remoteUser, localHostname, autoPassphrase, existingConfig)
// }

func (m *SshManager) finalize(
	ctx context.Context,
	priv any,
	pub any,
	algoName string,
	remoteHost string,
	remoteUser string,
	localHostname string,
	autoPassphrase bool,
	existingConfig *dagger.File,
) (*KeyResult, error) {

	hash := sha256.Sum256([]byte(time.Now().String()))
	shortHash := hex.EncodeToString(hash[:])[:6]

	keyBaseName := fmt.Sprintf("%s.%s_%s_%s", remoteHost, remoteUser, localHostname, shortHash)
	fullKeyName := fmt.Sprintf("id_%s_%s", algoName, keyBaseName)
	comment := fmt.Sprintf("%s@%s", remoteUser, remoteHost)

	passStr := ""
	if autoPassphrase {
		passStr = m.randomString(32)
	}

	var pemBlock *pem.Block
	var err error
	if passStr != "" {
		pemBlock, err = ssh.MarshalPrivateKeyWithPassphrase(priv, comment, []byte(passStr))
	} else {
		pemBlock, err = ssh.MarshalPrivateKey(priv, comment)
	}
	if err != nil {
		return nil, err
	}
	privContents := string(pem.EncodeToMemory(pemBlock))

	sshPubKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		return nil, err
	}
	pubString := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(sshPubKey))) + " " + comment + "\n"

	configEntry := fmt.Sprintf("\nHost                 %s.%s\n"+
		"Hostname             %s\n"+
		"IdentitiesOnly       yes\n"+
		"IdentityFile         ~/.ssh/%s\n"+
		"User                 git\n"+
		"ProxyCommand         nc -X 5 -x 127.0.0.1:9050 %%h %%p\n",
		remoteHost, remoteUser, remoteHost, fullKeyName)

	var configContent string
	if existingConfig != nil {
		old, _ := existingConfig.Contents(ctx)
		configContent = old + configEntry
	} else {
		configContent = configEntry
	}

	outputDir := dag.Directory().
		WithNewFile(fullKeyName, privContents).
		WithNewFile(fullKeyName+".pub", pubString).
		WithNewFile("config", configContent)

	return &KeyResult{
		PublicKey:  outputDir.File(fullKeyName + ".pub"),
		PrivateKey: dag.SetSecret(fullKeyName, privContents),
		Config:     outputDir.File("config"),
		Files:      outputDir,
	}, nil
}

func (m *SshManager) randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		ret[i] = charset[num.Int64()]
	}
	return string(ret)
}
