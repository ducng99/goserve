package responsewriter

import (
	"net/http"
)

type WriteReturn struct {
	BytesWritten int
	Err          error
}

type CustomResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Returned   *WriteReturn
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	if w.StatusCode == 0 {
		w.StatusCode = statusCode
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	if w.StatusCode == 0 {
		w.StatusCode = http.StatusOK
	}

	bytesWritten, err := w.ResponseWriter.Write(b)

	w.Returned = &WriteReturn{bytesWritten, err}

	return bytesWritten, err
}
