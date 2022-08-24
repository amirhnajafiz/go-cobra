package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

type HttpHandlerFunc func(http.ResponseWriter, *http.Request)

type Middleware struct {
	Logger *zap.Logger
	Token  string
}
