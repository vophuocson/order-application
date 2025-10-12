package middleware

import (
	"net/http"
	"time"
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
			logger.WithContext(r.Context()).Info("Start request", domainoutport.LogFields{
				"method":      r.Method,
				"request_url": r.URL,
			})
			wrapped := newResponseWriterWrapper(w)
			defer func(startTime time.Time) {
				logger.WithContext(r.Context()).Info("end request", domainoutport.LogFields{
					"status":   wrapped.Status(),
					"duration": time.Since(startTime).Microseconds(),
				})
			}(time.Now())

			h.ServeHTTP(wrapped, r)
		})
	}
}
