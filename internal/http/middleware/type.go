package middleware

import (
	"github.com/amirhnajafiz/go-cobra/internal/config"
	"net/http"
)

type HttpHandlerFunc func(http.ResponseWriter, *http.Request)

type Middleware struct {
	Configuration config.Config
}
