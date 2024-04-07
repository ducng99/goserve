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
	StatusCode        chan int
	statusCodeWritten bool
	Returned          *WriteReturn
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	if !w.statusCodeWritten {
		w.statusCodeWritten = true
		w.StatusCode <- statusCode
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	if !w.statusCodeWritten {
		w.statusCodeWritten = true
		w.StatusCode <- http.StatusOK
	}

	bytesWritten, err := w.ResponseWriter.Write(b)

	w.Returned = &WriteReturn{bytesWritten, err}

	return bytesWritten, err
}
