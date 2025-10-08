package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

var Recoverer = func(h http.Handler) http.Handler {
	return middleware.Recoverer(h)
}
