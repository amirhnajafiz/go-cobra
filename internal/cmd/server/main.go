package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/amirhnajafiz/go-cobra/internal/config"
	"github.com/amirhnajafiz/go-cobra/internal/http/handler"
	"github.com/amirhnajafiz/go-cobra/internal/middleware"
	"github.com/amirhnajafiz/go-cobra/pkg/encrypt"
	"github.com/amirhnajafiz/go-cobra/pkg/runner"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

type Setup struct {
	Configuration config.Config
	DB            *gorm.DB
	Logger        *zap.Logger
}

func (s Setup) HandleRequests() {
	var httpsSrv *http.Server
	var httpSrv *http.Server
	var m *autocert.Manager
	hdl := handler.Handler{
		DB: s.DB,
		Runner: runner.Runner{
			DB:     s.DB,
			Logger: s.Logger.Named("runner"),
		},
	}
	mid := middleware.Middleware{
		Configuration: s.Configuration,
	}
	enc := encrypt.Encrypt{
		Logger: s.Logger.Named("encrypt"),
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tasks", mid.Auth(hdl.AllTasksHandler())).Methods("GET")
	router.HandleFunc("/tasks/{page}", mid.Auth(hdl.AllTasksHandler())).Methods("GET")
	router.HandleFunc("/run", mid.Auth(hdl.NewRunHandler())).Methods("POST")
	router.HandleFunc("/tasks", mid.Auth(hdl.NewTaskHandler())).Methods("POST")
	router.HandleFunc("/task/{id}", mid.Auth(hdl.DeleteTaskHandler())).Methods("DELETE")
	router.HandleFunc("/task/{id}", mid.Auth(hdl.ViewTaskHandler())).Methods("GET")
	router.HandleFunc("/task/{id}", mid.Auth(hdl.UpdateTaskHandler())).Methods("PUT")

	if s.Configuration.SSLMode == "development" {
		// Generate ca.crt and ca.key if not found
		caFile, err := os.Open("certs/ca.crt")
		if err != nil {
			enc.GenerateCertificateAuthority()
		}
		defer func(caFile *os.File) {
			_ = caFile.Close()
		}(caFile)
		// Generate cert.pem and key.pem for https://localhost
		enc.GenerateCert()

		// Launch HTTPS server
		s.Logger.Info("Starting server https://" + s.Configuration.Host + ":" + s.Configuration.Port)
		s.Logger.Fatal(http.ListenAndServeTLS(":"+s.Configuration.Port, "certs/cert.pem", "certs/key.pem", handlers.LoggingHandler(os.Stdout, router)).Error())
	}

	if s.Configuration.SSLMode == "production" {

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
			allowedHost := s.Configuration.Host
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

		httpsSrv.Addr = s.Configuration.Host + ":443"
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
			s.Logger.Info("Starting server http://localhost")

			err := httpSrv.ListenAndServe()
			if err != nil {
				s.Logger.Fatal("httpSrv.ListenAndServe() failed with " + err.Error())
			}
		}()

		// Launch HTTPS server
		s.Logger.Info("Starting server https://" + s.Configuration.Host + ":" + s.Configuration.Port)
		s.Logger.Fatal(httpsSrv.ListenAndServeTLS("", "").Error())
	}
}
