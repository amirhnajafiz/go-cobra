package handler

import (
	"cmd/config"
	"cmd/internal/handler"
	"cmd/internal/middleware"
	"cmd/pkg/encrypt"
	"cmd/server"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

func HandleRequests(configuration config.Config, db *gorm.DB) {
	var httpsSrv *http.Server
	var httpSrv *http.Server
	var m *autocert.Manager

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tasks", middleware.CheckSecurity(configuration, handler.AllTasksHandler(db))).Methods("GET")

	if configuration.SSLMode == "development" {
		// Generate ca.crt and ca.key if not found
		caFile, err := os.Open("certs/ca.crt")
		if err != nil {
			encrypt.GenerateCertificateAuthority()
		}
		defer caFile.Close()
		// Generate cert.pem and key.pem for https://localhost
		encrypt.GenerateCert()

		// Launch HTTPS server
		fmt.Println("Starting server https://" + configuration.Host + ":" + configuration.Port)
		log.Fatal(http.ListenAndServeTLS(":"+configuration.Port, "certs/cert.pem", "certs/key.pem", handlers.LoggingHandler(os.Stdout, router)))
	}

	if configuration.SSLMode == "production" {

		// Manage Let's Encrypt SSL

		// Note: use a sensible value for data directory
		// this is where cached certificates are stored

		httpsSrv = &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      router,
		}

		dataDir := "certs/"
		hostPolicy := func(ctx context.Context, host string) error {
			// Note: change to your real domain
			allowedHost := configuration.Host
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}

		m = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}

		httpsSrv.Addr = configuration.Host + ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		// Spin up web server on port 80 to listen for auto-cert HTTP challenge
		httpSrv = server.MakeHTTPServer()
		httpSrv.Addr = ":80"

		// allow auto-cert handle Let's Encrypt auth callbacks over HTTP.
		if m != nil {
			// https://github.com/golang/go/issues/21890
			httpSrv.Handler = m.HTTPHandler(httpSrv.Handler)
		}

		// Launch HTTP server
		go func() {
			fmt.Println("Starting server http://localhost")

			err := httpSrv.ListenAndServe()
			if err != nil {
				log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
			}
		}()

		// Launch HTTPS server
		fmt.Println("Starting server https://" + configuration.Host + ":" + configuration.Port)
		log.Fatal(httpsSrv.ListenAndServeTLS("", ""))
	}
}
