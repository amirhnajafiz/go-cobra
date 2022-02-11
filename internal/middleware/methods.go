package middleware

import (
	"net/http"
)

func (m Middleware) Auth(next HttpHandlerFunc) HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("token")
		if header != m.Configuration.Token {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("401 - Unauthorized"))
			return
		}
		next(w, r)
	}
}
