package middleware

import (
	"cmd/internal/config"
	"net/http"
)

type HttpHandlerFunc func(http.ResponseWriter, *http.Request)

type Middleware struct {
	Configuration config.Config
}
