package server

import (
	"cmd/internal/config"
	handler2 "cmd/internal/http/handler"
	"cmd/internal/middleware"
	"cmd/pkg/encrypt"
	logger "cmd/pkg/zap-logger"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

type Setup struct {
	configuration config.Config
	db            *gorm.DB
}

func (s Setup) HandleRequests() {
	var httpsSrv *http.Server
	var httpSrv *http.Server
	var m *autocert.Manager
	handler := handler2.Handler{
		DB: s.db,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tasks", middleware.Auth(s.configuration, handler.AllTasksHandler())).Methods("GET")
	router.HandleFunc("/tasks/{page}", middleware.Auth(s.configuration, handler.AllTasksHandler())).Methods("GET")
	router.HandleFunc("/run", middleware.Auth(s.configuration, handler.NewRunHandler())).Methods("POST")
	router.HandleFunc("/tasks", middleware.Auth(s.configuration, handler.NewTaskHandler())).Methods("POST")
	router.HandleFunc("/task/{id}", middleware.Auth(s.configuration, handler.DeleteTaskHandler())).Methods("DELETE")
	router.HandleFunc("/task/{id}", middleware.Auth(s.configuration, handler.ViewTaskHandler())).Methods("GET")
	router.HandleFunc("/task/{id}", middleware.Auth(s.configuration, handler.UpdateTaskHandler())).Methods("PUT")

	if s.configuration.SSLMode == "development" {
		// Generate ca.crt and ca.key if not found
		caFile, err := os.Open("certs/ca.crt")
		if err != nil {
			encrypt.GenerateCertificateAuthority()
		}
		defer func(caFile *os.File) {
			_ = caFile.Close()
		}(caFile)
		// Generate cert.pem and key.pem for https://localhost
		encrypt.GenerateCert()

		// Launch HTTPS server
		logger.GetLogger().Info("Starting server https://" + s.configuration.Host + ":" + s.configuration.Port)
		logger.GetLogger().Fatal(http.ListenAndServeTLS(":"+s.configuration.Port, "certs/cert.pem", "certs/key.pem", handlers.LoggingHandler(os.Stdout, router)).Error())
	}

	if s.configuration.SSLMode == "production" {

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
			allowedHost := s.configuration.Host
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

		httpsSrv.Addr = s.configuration.Host + ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		// Spin up web server on port 80 to listen for auto-cert HTTP challenge
		httpSrv = MakeHTTPServer()
		httpSrv.Addr = ":80"

		// allow auto-cert handle Let's Encrypt auth callbacks over HTTP.
		if m != nil {
			// https://github.com/golang/go/issues/21890
			httpSrv.Handler = m.HTTPHandler(httpSrv.Handler)
		}

		// Launch HTTP server
		go func() {
			logger.GetLogger().Info("Starting server http://localhost")

			err := httpSrv.ListenAndServe()
			if err != nil {
				logger.GetLogger().Fatal("httpSrv.ListenAndServe() failed with " + err.Error())
			}
		}()

		// Launch HTTPS server
		logger.GetLogger().Info("Starting server https://" + s.configuration.Host + ":" + s.configuration.Port)
		logger.GetLogger().Fatal(httpsSrv.ListenAndServeTLS("", "").Error())
	}
}
