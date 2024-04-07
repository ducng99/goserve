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
	"math/big"
	"os"
	"time"

	"r.tomng.dev/goserve/internal/logger"
)

// Keys generates a new P256 ECDSA public private key pair for TLS.
// It returns a bytes buffer for the PEM encoded private key and certificate.
func NewKeys(validFor time.Duration) (cert, key *bytes.Buffer, fingerprint [32]byte, err error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Fatalf("failed to generate private key: %v", err)
		return nil, nil, fingerprint, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logger.Fatalf("failed to generate serial number: %v", err)
		return nil, nil, fingerprint, err
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
		return nil, nil, fingerprint, err
	}

	// Encode and write certificate and key to bytes.Buffer
	cert = bytes.NewBuffer([]byte{})
	pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	key = bytes.NewBuffer([]byte{})
	pem.Encode(key, pemBlockForKey(privKey))

	fingerprint = sha256.Sum256(derBytes)

	return cert, key, fingerprint, nil
}

func pemBlockForKey(key *ecdsa.PrivateKey) *pem.Block {
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
		os.Exit(2)
	}
	return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
}
