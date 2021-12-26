package middleware

import (
	"cmd/config"
	"net/http"
)

type HttpHandlerFunc func(http.ResponseWriter, *http.Request)

func Auth(configuration config.Config, next HttpHandlerFunc) HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("token")
		if header != configuration.Token {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("401 - Unauthorized"))
			return
		}
		next(w, r)
	}
}
