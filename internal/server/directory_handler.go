package server

import (
	"net/http"

	"r.tomng.dev/goserve/internal/files"
	"r.tomng.dev/goserve/internal/tmpl"
)

// Handler for directory requests.
// Display an indexing page of contents in the directory
func (c *ServerConfig) directoryHandler(w http.ResponseWriter, r *http.Request, path string) {
	relativePath := files.RelativeRoot(c.RootDir, path)

	entries, err := files.GetEntries(path)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set some security headers
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; frame-ancestors 'self'; form-action 'self';")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Feature-Policy", "geolocation 'none'; microphone 'none'; camera 'none'; speaker 'none'; vibrate 'none'; payment 'none'; usb 'none';")

	tmpl.RenderDirectoryView(w, r, relativePath, entries)
}
