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
			StatusCode:     make(chan int, 1),
			Returned:       nil,
		}

		waitServing := make(chan bool, 1)

		// Serve the request in a goroutine so we can capture status code earlier
		go func(w http.ResponseWriter, r *http.Request, done chan bool) {
			next.ServeHTTP(w, r)
			done <- true
		}(customWriter, r, waitServing)

		// Log the request once we have the status code
		statusCode := <-customWriter.StatusCode

		// Log connection with colors
		logResultFormat := "%s [%d]: %s %s\n"
		logResultParams := []interface{}{r.RemoteAddr, statusCode, r.Method, r.URL.Path}
		switch {
		case statusCode >= 500:
			logger.Printf(logger.LogError, logResultFormat, logResultParams...)
		case statusCode >= 400:
			logger.Printf(logger.LogWarning, logResultFormat, logResultParams...)
		default:
			logger.Printf(logger.LogSuccess, logResultFormat, logResultParams...)
		}

		// Wait for the request to finish then continue logging
		// We can't wait in the middle of a middleware because the request is still being processed
		// Do this in a goroutine to avoid blocking the request
		go func(w *responsewriter.CustomResponseWriter, r *http.Request, start time.Time) {
			<-r.Context().Done()
			close(customWriter.StatusCode)

			writeReturn := w.Returned
			bytesWritten := 0

			if writeReturn != nil {
				if writeReturn.Err != nil {
					logger.Printf(logger.LogError, "%s Error writing response: %v\n", r.RemoteAddr, writeReturn.Err)
				}

				bytesWritten = writeReturn.BytesWritten
			}

			logger.Printf(logger.LogNormal, "%s Closing - written %d bytes - %s\n", r.RemoteAddr, bytesWritten, time.Since(start))
		}(customWriter, r, start)

		<-waitServing
		close(waitServing)
	})
}
