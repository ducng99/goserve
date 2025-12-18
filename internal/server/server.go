package server

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

	"github.com/ducng99/goserve/internal/files"
	"github.com/ducng99/goserve/internal/logger"
	"github.com/ducng99/goserve/internal/proxy"
	"github.com/ducng99/goserve/internal/server/assets"
	"github.com/ducng99/goserve/internal/server/middlewares"
	"github.com/ducng99/goserve/internal/ssl"
)

var SelfSignedSSLPath = filepath.Join(os.TempDir(), "goserve")

// Starts web server
func (c *ServerConfig) StartServer() {
	// Set up routes
	mux := c.NewServeMux()

	// Setup HTTPS if enabled
	c.SetupSSL()

	// Start server
	listenAddr := net.JoinHostPort(c.Host, c.Port)

	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	go func() {
		var err error

		if c.HttpsEnabled {
			err = httpServer.ListenAndServeTLS(c.CertPath, c.KeyPath)
		} else {
			err = httpServer.ListenAndServe()
		}

		if !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("HTTP server error: %v\n", err)
		}
		logger.Printf(logger.LogNormal, "Interrupted. Shutting down...\n")
	}()

	protocol := "http"

	if c.HttpsEnabled {
		protocol = "https"
	}

	serverURL := protocol + "://" + listenAddr

	logger.Printf(logger.LogNormal, "Started goserve %s server (%s)\n", strings.ToUpper(protocol), serverURL)

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

// NewServeMux creates a new HTTP ServeMux with configured routes
func (c *ServerConfig) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	var routeHandler http.Handler

	if c.ProxyToAddr != "" {
		proxyHandler, err := proxy.New(c.ProxyToAddr, c.ProxyHeadersEnabled, c.ProxyIgnoreRedirect)
		if err != nil {
			logger.Fatalf("Error creating reverse proxy handler: %v\n", err)
		}

		routeHandler = proxyHandler
	} else {
		routeHandler = http.Handler(http.HandlerFunc(c.routeHandlerFunc))
	}

	if c.CorsEnabled {
		routeHandler = middlewares.CorsMiddleware(routeHandler)
	}

	routeHandler = middlewares.LogConnectionMiddleware(routeHandler)
	mux.Handle("/", routeHandler)
	mux.HandleFunc(assets.PrefixPath+"{asset}", assets.AssetsHandler)

	return mux
}

// Handler for all requests.
// Serves files or display directory index
func (c *ServerConfig) routeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	sanitisedPath, err := files.SanitisePath(c.RootDir, r.URL.Path)
	if err != nil {
		switch {
		case errors.Is(err, files.ErrorSanitiseNotExists):
			http.Error(w, "Path not found", http.StatusNotFound)
		case errors.Is(err, files.ErrorSanitiseUnauthorized):
			http.Error(w, "Not enough permission to read the given path", http.StatusForbidden)
		default:
			http.Error(w, "An unknown error occured", http.StatusInternalServerError)
			logger.Printf(logger.LogError, "%v\n", err)
		}
		return
	}

	pathType, err := files.GetPathType(sanitisedPath)
	if err != nil {
		http.Error(w, "Cannot get path type", http.StatusInternalServerError)
		logger.Printf(logger.LogError, "%v\n", err)
		return
	}

	switch pathType {
	case files.PathTypeFile:
		http.ServeFile(w, r, sanitisedPath)
	case files.PathTypeDirectory:
		c.directoryHandler(w, r, sanitisedPath)
	default:
		http.Error(w, "Path type not handled correctly", http.StatusInternalServerError)
	}
}

// Checks if HTTPS is enabled and sets up SSL keys if necessary
func (c *ServerConfig) SetupSSL() {
	if c.HttpsEnabled {
		if c.CertPath != "" && c.KeyPath != "" {
			f, err := os.Open(c.CertPath)
			if err != nil {
				logger.Fatalf("Cannot read cert file '%s': %v\n", c.CertPath, err)
			}
			f.Close()

			f, err = os.Open(c.KeyPath)
			if err != nil {
				logger.Fatalf("Cannot read key file at '%s': %v\n", c.KeyPath, err)
			}
			f.Close()
		} else if c.CertPath == "" && c.KeyPath == "" {
			certPath, privKeyPath, exists := ssl.KeysExist(SelfSignedSSLPath)

			if exists {
				logger.Printf(logger.LogNormal, "Using previous self-signed SSL certificate\n")
			} else {
				keyPair, err := ssl.NewKeys(365 * 24 * time.Hour)
				if err != nil {
					logger.Fatalf("Error generating SSL keys: %v\n", err)
				}

				certPath, privKeyPath, err = keyPair.Save(SelfSignedSSLPath)
				if err != nil {
					logger.Fatalf("%v\n", err)
				}

				logger.Printf(logger.LogNormal, "Generated SSL key fingerprint:\n% X\n", keyPair.Fingerprint)
			}

			c.CertPath = certPath
			c.KeyPath = privKeyPath
		} else {
			// Cobra already handles this but just in case
			logger.Fatalf("Both cert and key paths must be provided. Or both must be empty to use a self-signed certificate\n")
		}
	}
}
