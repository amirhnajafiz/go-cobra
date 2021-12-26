package server

import (
	"io"
	"net/http"
	"time"
)

const (
	htmlIndex = `<html><body>Welcome!</body></html>`
)

func handleIndex(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, htmlIndex)
}

func makeServerFromMux(mux *http.ServeMux) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func MakeHTTPServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleIndex)
	return makeServerFromMux(mux)
}
