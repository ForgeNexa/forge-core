package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSecurity_SafePath(t *testing.T) {
	tmpDir := t.TempDir()

	// Safe path
	valid, err := SafePath(tmpDir, "sub/file.txt")
	if err != nil {
		t.Fatalf("expected valid safe path, got error: %v", err)
	}
	expected := filepath.Join(tmpDir, "sub/file.txt")
	if valid != expected {
		t.Fatalf("expected %s, got %s", expected, valid)
	}

	// Path traversal attempt
	_, err = SafePath(tmpDir, "../outside.txt")
	if err == nil {
		t.Fatalf("expected path traversal detection, got nil error")
	}
}

func TestSecurity_AESEncryptionAndDecryption(t *testing.T) {
	masterKey := "super-secure-master-key-12345"
	original := []byte("Authorization: Bearer my-top-secret-token")

	encryptedHex, err := EncryptSecret(original, masterKey)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	decrypted, err := DecryptSecret(encryptedHex, masterKey)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if string(decrypted) != string(original) {
		t.Fatalf("expected %s, got %s", string(original), string(decrypted))
	}
}

func TestSecurity_Redaction(t *testing.T) {
	logLine := "Sending request with Authorization: Bearer sk-1234567890abcdef to https://api.openai.com"
	redacted := RedactSecrets(logLine, []string{"sk-1234567890abcdef"})

	if redacted != "Sending request with Authorization: Bearer [REDACTED] to https://api.openai.com" {
		t.Fatalf("redaction failed, got: %s", redacted)
	}
}

func TestSecurity_WriteAtomic(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "config.json")
	data := []byte(`{"version": "1.0.0"}`)

	if err := WriteAtomic(target, data, 0644); err != nil {
		t.Fatalf("WriteAtomic failed: %v", err)
	}

	read, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read back written file: %v", err)
	}

	if string(read) != string(data) {
		t.Fatalf("file content mismatch: expected %s, got %s", string(data), string(read))
	}
}
