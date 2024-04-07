package ssl_test

import (
	"testing"
	"time"

	"r.tomng.dev/goserve/internal/ssl"
)

func TestNewKeys(t *testing.T) {
	cert, key, fingerprint, err := ssl.NewKeys(365 * 24 * time.Hour)
	if err != nil {
		t.Fatalf("NewKeys() returned error: %v", err)
	}

	if cert == nil {
		t.Fatalf("NewKeys() returned nil cert")
	}

	if key == nil {
		t.Fatalf("NewKeys() returned nil key")
	}

	if len(fingerprint) != 32 {
		t.Fatalf("NewKeys() returned fingerprint with length %d, expected 32", len(fingerprint))
	}
}
