package utils

import (
	"net/http"
)

func RenderJsonMiddleware(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
