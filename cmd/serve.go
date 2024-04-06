/*
Copyright © 2024 Thomas Nguyen <tom@tomng.dev>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"r.tomng.dev/goserve/internal/files"
	"r.tomng.dev/goserve/internal/logger"
	"r.tomng.dev/goserve/internal/responsewriter"
	"r.tomng.dev/goserve/internal/tmpl"
)

const (
	DefaultListenHost = "0.0.0.0"
	DefaultListenPort = "8080"
)

var RootDir = filepath.Dir(os.Args[0])

// Default run function for root command.
//
// Handles flags and continue to start a server
func handleCommand(cmd *cobra.Command, args []string) {
	host := DefaultListenHost
	port := DefaultListenPort

	if len(args) > 0 {
		hostport := args[0]

		// Cannot split without a colon
		// Add a colon to split then use default port
		if !strings.Contains(hostport, ":") {
			hostport = hostport + ":"
		}

		_host, _port, err := net.SplitHostPort(hostport)
		if err != nil {
			logger.Fatalf("Invalid address (%v)\n", err)
		}

		if _host != "" {
			host = _host
		}

		if _port != "" {
			port = _port
		}
	}

	rootDir, err := cmd.Flags().GetString("dir")
	if err != nil {
		logger.Fatalf("Error getting flag: %v\n", err)
	}

	rootDir, err = filepath.EvalSymlinks(filepath.Clean(rootDir))
	if err != nil {
		logger.Fatalf("Error resolving directory: %v\n", err)
	}

	rootDir, err = filepath.Abs(rootDir)
	if err != nil {
		logger.Fatalf("Error getting absolute path: %v\n", err)
	}

	RootDir = rootDir

	startServer(host, port)
}

// Starts web server
func startServer(host, port string) {
	mux := http.NewServeMux()
	mux.Handle("/", logConnectionMiddleware(http.HandlerFunc(routeHandler)))

	listenAddr := net.JoinHostPort(host, port)

	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("HTTP server error: %v\n", err)
		}
		logger.Printf(logger.LogNormal, "Interrupted. Shutting down...\n")
	}()

	serverURL := "http://" + listenAddr

	logger.Printf(logger.LogNormal, "Started goserve HTTP server (%s)\n", serverURL)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("HTTP shutdown error: %v\n", err)
	}

	logger.Printf(logger.LogNormal, "Server stopped\n")
}

// Handler for all requests
func routeHandler(w http.ResponseWriter, r *http.Request) {
	sanitisedPath, err := files.SanitisePath(RootDir, r.URL.Path)
	if err != nil {
		switch {
		case errors.Is(err, files.ErrorSanitiseNotExists):
			http.Error(w, "File not found", http.StatusNotFound)
		case errors.Is(err, files.ErrorSanitiseUnauthorized):
			http.Error(w, "Unauthorized path", http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	pathType, err := files.GetPathType(sanitisedPath)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	switch pathType {
	case files.PathTypeFile:
		http.ServeFile(w, r, sanitisedPath)
	case files.PathTypeDirectory:
		handleDirectory(w, r, sanitisedPath)
	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func handleDirectory(w http.ResponseWriter, r *http.Request, path string) {
	relativePath := files.RelativeRoot(RootDir, path)

	entries, err := files.GetEntries(path)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set some security headers
	w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; script-src 'none'; style-src 'none'; img-src 'none'; form-action 'none';")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Feature-Policy", "geolocation 'none'; microphone 'none'; camera 'none'; speaker 'none'; vibrate 'none'; payment 'none'; usb 'none';")

	tmpl.RenderDirectoryView(w, r, relativePath, entries)
}

// Middleware to log connection details
func logConnectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf(logger.LogNormal, "%s Accepted\n", r.RemoteAddr)

		// Writer wrapper to capture status code
		start := time.Now()

		customWriter := &responsewriter.CustomResponseWriter{
			ResponseWriter: w,
			StatusCode:     0,
			Returned:       nil,
		}

		// Call the original handler
		next.ServeHTTP(customWriter, r)

		// Log the request once we have the status code
		statusCode := customWriter.StatusCode

		// Log connection with colors
		logResultFormat := "%s [%d]: %s %s - %s\n"
		logResultParams := []interface{}{r.RemoteAddr, statusCode, r.Method, r.URL.Path, time.Since(start)}
		switch {
		case statusCode >= 500:
			logger.Printf(logger.LogError, logResultFormat, logResultParams...)
		case statusCode >= 400:
			logger.Printf(logger.LogWarning, logResultFormat, logResultParams...)
		default:
			logger.Printf(logger.LogSuccess, logResultFormat, logResultParams...)
		}

		// Wait for the response to be written
		writeReturn := customWriter.Returned

		if writeReturn != nil {
			if writeReturn.Err != nil {
				logger.Printf(logger.LogError, "%s Error writing response: %v\n", r.RemoteAddr, writeReturn.Err)
				// } else {
				// logger.Printf(logger.LogNormal, "%s Written %d bytes\n", r.RemoteAddr, writeReturn.BytesWritten)
			}
		}

		logger.Printf(logger.LogNormal, "%s Closing\n", r.RemoteAddr)
	})
}
