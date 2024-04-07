package middlewares

import (
	"net/http"
	"time"

	"r.tomng.dev/goserve/internal/logger"
	"r.tomng.dev/goserve/internal/responsewriter"
)

// Middleware to log connection details
func LogConnectionMiddleware(next http.Handler) http.Handler {
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
