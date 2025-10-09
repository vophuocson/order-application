package middleware

import (
	"net/http"
	applicationoutbound "user-domain/internal/application/outbound"
	domainoutport "user-domain/internal/domain/outport"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	code int
}

func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
	}
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Status() int {
	return w.code
}

func LoggingMiddleware(logger applicationoutbound.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := newResponseWriterWrapper(w)

			defer func() {
				logger.WithContext(r.Context()).Info("in", domainoutport.LogFields{})
			}()

			h.ServeHTTP(wrapped, r)
		})
	}
}
