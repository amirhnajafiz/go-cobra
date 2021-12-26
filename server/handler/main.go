package handler

import (
	"cmd/config"
	"cmd/pkg/encrypt"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
	"time"
)

func HandleRequests(configuration config.Config) {
	router := mux.NewRouter().StrictSlash(true)

	if configuration.SSLMode == "development" {
		// Generate ca.crt and ca.key if not found
		caFile, err := os.Open("certs/ca.crt")
		if err != nil {
			encrypt.GenerateCertificateAuthority()
		}
		defer caFile.Close()
		// Generate cert.pem and key.pem for https://localhost
		encrypt.GenerateCert()
	}

	if configuration.SSLMode == "production" {

		// Manage Let's Encrypt SSL

		// Note: use a sensible value for data directory
		// this is where cached certificates are stored

		httpsSrv := &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      router,
		}

		//  handlers.LoggingHandler(os.Stdout, router

		dataDir := "certs/"
		hostPolicy := func(ctx context.Context, host string) error {
			// Note: change to your real domain
			allowedHost := configuration.Host
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}

		_ = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}
	}
}
