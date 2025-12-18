package ssl_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ducng99/goserve/internal/ssl"
)

func TestNewKeys(t *testing.T) {
	keyPair, err := ssl.NewKeys(365 * 24 * time.Hour)
	if err != nil {
		t.Fatalf("NewKeys() returned error: %v", err)
	}

	if keyPair.Cert == nil {
		t.Fatalf("NewKeys() returned nil cert")
	}

	if keyPair.Key == nil {
		t.Fatalf("NewKeys() returned nil key")
	}

	if len(keyPair.Fingerprint) != 32 {
		t.Fatalf("NewKeys() returned fingerprint with length %d, expected 32", len(keyPair.Fingerprint))
	}
}

func TestSaveKeys(t *testing.T) {
	keyPair, err := ssl.NewKeys(5 * time.Minute)
	if err != nil {
		t.Fatalf("NewKeys() returned error: %v", err)
	}

	keysSavePath, err := os.MkdirTemp(".", "goserve_*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(keysSavePath)

	certPath, keyPath, err := keyPair.Save(keysSavePath)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	if certPath != filepath.Join(keysSavePath, ssl.CertFileName) {
		t.Fatalf("Save() returned wrong cert path: %s", certPath)
	}

	if keyPath != filepath.Join(keysSavePath, ssl.KeyFileName) {
		t.Fatalf("Save() returned wrong key path: %s", keyPath)
	}

	if _, err := os.Stat(certPath); err != nil {
		t.Fatalf("Save() failed to save cert: %v", err)
	}

	if _, err := os.Stat(keyPath); err != nil {
		t.Fatalf("Save() failed to save key: %v", err)
	}
}
