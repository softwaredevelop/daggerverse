// SSH Manager module for generating secure SSH key pairs and configurations.
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

// SshManager provides services to generate and manage SSH identities.
type SshManager struct{} //revive:disable-line

// KeyResult represents the collection of generated SSH components.
type KeyResult struct {
	// The OpenSSH formatted public key file.
	PublicKey *dagger.File `json:"publicKey"`
	// The OpenSSH formatted private key, stored as a secret.
	PrivateKey *dagger.Secret `json:"privateKey"`
	// The SSH configuration file containing host and proxy settings.
	Config *dagger.File `json:"config"`
	// A directory containing the private key, public key, and config for easy export.
	Files *dagger.Directory `json:"files"`
}

// finalize is an internal helper to process raw keys into formatted files and secrets.
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

// randomString generates a cryptographically secure random string for passphrases.
func (m *SshManager) randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		ret[i] = charset[num.Int64()]
	}
	return string(ret)
}
