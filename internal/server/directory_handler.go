package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/ducng99/goserve/internal/files"
	"github.com/ducng99/goserve/internal/logger"
	"github.com/ducng99/goserve/internal/tmpl/dirview"
)

var ErrNonceGeneration = errors.New("failed to generate nonce")

// Handler for directory requests.
// Display an indexing page of contents in the directory
func (c *ServerConfig) directoryHandler(w http.ResponseWriter, r *http.Request, dirPath string) {
	// Get files in the provided directory
	relativePath := files.RelativeRoot(c.RootDir, dirPath)

	entries, err := files.GetEntries(dirPath)
	if err != nil {
		http.Error(w, "Cannot get entries in the provided directory", http.StatusInternalServerError)
		logger.Printf(logger.LogError, "%v\n", err)
		return
	}

	// Generate nonce for CSP
	nonce, err := generateNonce()
	if err != nil {
		http.Error(w, "Cannot generate nonce for CSP", http.StatusInternalServerError)
		logger.Printf(logger.LogError, "%v\n", err)
		return
	}

	// Set some security headers
	w.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src 'none'; script-src 'nonce-%s'; connect-src 'self'; img-src 'self'; style-src 'nonce-%s'; frame-ancestors 'self'; form-action 'self';", nonce, nonce))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Permissions-Policy", "accelerometer=(),ambient-light-sensor=(),autoplay=(),battery=(),camera=(),display-capture=(),document-domain=(),encrypted-media=(),fullscreen=(),gamepad=(),geolocation=(),gyroscope=(),magnetometer=(),microphone=(),midi=(),payment=(),picture-in-picture=(),publickey-credentials-get=(),speaker-selection=(),sync-xhr=(self),usb=(),screen-wake-lock=(),web-share=(),xr-spatial-tracking=()")

	dirview.Render(w, r, relativePath, entries, nonce, c.DirViewTheme)
}

func generateNonce() (string, error) {
	const length = 16

	nonce := make([]byte, length)
	genLength, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}

	if genLength != length {
		return "", ErrNonceGeneration
	}

	return hex.EncodeToString(nonce), nil
}
