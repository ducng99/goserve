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
	mux := http.NewServeMux()

	routeHandlerFunc := middlewares.LogConnectionMiddleware(http.HandlerFunc(c.routeHandler))

	if c.CorsEnabled {
		routeHandlerFunc = middlewares.CorsMiddleware(routeHandlerFunc)
	}

	mux.Handle("/", routeHandlerFunc)

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

// Handler for all requests
func (c *ServerConfig) routeHandler(w http.ResponseWriter, r *http.Request) {
	sanitisedPath, err := files.SanitisePath(c.RootDir, r.URL.Path)
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
