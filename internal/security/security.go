package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SafePath verifies that the requested subpath stays strictly inside baseDir,
// preventing path traversal attacks (e.g. "../../../etc/passwd").
func SafePath(baseDir, untrustedPath string) (string, error) {
	cleanBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("invalid base path: %w", err)
	}

	joined := filepath.Join(cleanBase, untrustedPath)
	cleanTarget, err := filepath.Abs(joined)
	if err != nil {
		return "", fmt.Errorf("invalid target path: %w", err)
	}

	if !strings.HasPrefix(cleanTarget, cleanBase) {
		return "", fmt.Errorf("access denied: path traversal attempt detected outside workspace")
	}

	return cleanTarget, nil
}

// EncryptSecret encrypts sensitive data with AES-256-GCM using a hashed master key.
func EncryptSecret(plaintext []byte, masterKey string) (string, error) {
	key := sha256.Sum256([]byte(masterKey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

// DecryptSecret decrypts a previously encrypted AES-256-GCM hex string.
func DecryptSecret(cipherHex string, masterKey string) ([]byte, error) {
	data, err := hex.DecodeString(cipherHex)
	if err != nil {
		return nil, err
	}

	key := sha256.Sum256([]byte(masterKey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("malformed ciphertext: length less than nonce")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// RedactSecrets scans a string or log line and replaces matches with [REDACTED].
func RedactSecrets(input string, secrets []string) string {
	res := input
	for _, sec := range secrets {
		if len(sec) > 3 {
			res = strings.ReplaceAll(res, sec, "[REDACTED]")
		}
	}
	return res
}

// WriteAtomic safely writes data to a temporary file before renaming to target,
// preventing file corruption on crashes or power loss.
func WriteAtomic(filename string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(filename)
	tmpFile, err := os.CreateTemp(dir, "atomic-*")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()

	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpName)
	}()

	if _, err := tmpFile.Write(data); err != nil {
		return err
	}
	if err := tmpFile.Chmod(perm); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, filename)
}
