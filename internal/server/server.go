package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"r.tomng.dev/goserve/internal/files"
	"r.tomng.dev/goserve/internal/logger"
	"r.tomng.dev/goserve/internal/server/middlewares"
)

// Starts web server
func (c *ServerConfig) StartServer() {
	// Set up routes
	mux := c.newServeMux()

	// Start server
	listenAddr := net.JoinHostPort(c.Host, c.Port)

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

func (c *ServerConfig) newServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	routeHandler := http.Handler(http.HandlerFunc(c.routeHandlerFunc))

	if c.CorsEnabled {
		routeHandler = middlewares.CorsMiddleware(routeHandler)
	}

	routeHandler = middlewares.LogConnectionMiddleware(routeHandler)
	mux.Handle("/", routeHandler)

	return mux
}

// Handler for all requests
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
