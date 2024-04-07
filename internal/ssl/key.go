/*
Original source:
https://github.com/suyashkumar/ssl-proxy/blob/169fda92ebf3ce91bc05b691252124410b96e3cb/gen/gen.go
*/
package ssl

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"r.tomng.dev/goserve/internal/logger"
)

type KeyPair struct {
	Cert        *bytes.Buffer
	Key         *bytes.Buffer
	Fingerprint [32]byte
}

const (
	CertFileName = "goserve_cert.crt"
	KeyFileName  = "goserve_privatekey.key"
)

// Generates a new P256 ECDSA public private key pair for TLS.
// It returns a bytes buffer for the PEM encoded private key and certificate.
func NewKeys(validFor time.Duration) (*KeyPair, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Fatalf("failed to generate private key: %v", err)
		return nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logger.Fatalf("failed to generate serial number: %v", err)
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"goserve"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		logger.Fatalf("Failed to create certificate: %v", err)
		return nil, err
	}

	// Encode and write certificate and key to bytes.Buffer
	cert := bytes.NewBuffer([]byte{})
	pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	key := bytes.NewBuffer([]byte{})
	pem.Encode(key, pemBlockForKey(privKey))

	fingerprint := sha256.Sum256(derBytes)

	keyPair := &KeyPair{
		Cert:        cert,
		Key:         key,
		Fingerprint: fingerprint,
	}

	return keyPair, nil
}

func pemBlockForKey(key *ecdsa.PrivateKey) *pem.Block {
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
		os.Exit(2)
	}
	return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
}

// Saves certificate and private key to the given directory.
// Uses previously generated key pair if they exist.
func (k *KeyPair) Save(dir string) (string, string, error) {
	// Files don't exist, create them
	err := os.MkdirAll(dir, fs.ModeDir|0700)
	if err != nil {
		return "", "", fmt.Errorf("Error creating temp dir for self-signed SSL: %v", err)
	}

	certPath := filepath.Join(dir, CertFileName)
	f, err := os.Create(certPath)
	if err != nil {
		return "", "", fmt.Errorf("Error creating cert file: %v", err)
	}
	f.Write(k.Cert.Bytes())
	f.Close()

	privKeyPath := filepath.Join(dir, KeyFileName)
	f, err = os.OpenFile(privKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", "", fmt.Errorf("Error creating private key file: %v", err)
	}
	f.Write(k.Key.Bytes())
	f.Close()

	return certPath, privKeyPath, nil
}

// Checks if the certificate and private key files already exist in the given directory.
// Returns their paths and a boolean indicating if they exist.
func KeysExist(dir string) (string, string, bool) {
	certPath := filepath.Join(dir, CertFileName)
	privKeyPath := filepath.Join(dir, KeyFileName)

	// Return cert and key files if they already exists
	if _, err := os.Stat(certPath); err == nil {
		if _, err := os.Stat(privKeyPath); err == nil {
			return certPath, privKeyPath, true
		}
	}

	return "", "", false
}
