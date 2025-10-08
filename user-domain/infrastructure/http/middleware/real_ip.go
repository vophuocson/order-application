package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

var RealIP = func(h http.Handler) http.Handler {
	return middleware.RealIP(h)
}
