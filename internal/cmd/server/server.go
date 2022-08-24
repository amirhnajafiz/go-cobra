package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/amirhnajafiz/go-cobra/internal/config"
	"github.com/amirhnajafiz/go-cobra/internal/database"
	"github.com/amirhnajafiz/go-cobra/internal/http/handler"
	"github.com/amirhnajafiz/go-cobra/internal/http/middleware"
	"github.com/amirhnajafiz/go-cobra/internal/runner"
	"github.com/amirhnajafiz/go-cobra/pkg/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"gorm.io/gorm"
)

type Server struct {
	Configuration config.Config
	DB            *gorm.DB
	Logger        *zap.Logger
}

func main() {
	l := logger.New()

	cfg, err := config.LoadConfiguration()
	if err != nil {
		panic(err)
	}

	conn, err := database.Connect()
	if err != nil {
		panic(err)
	}

	h := handler.Handler{
		DB: conn,
		Runner: runner.Runner{
			DB:     conn,
			Logger: l.Named("runner"),
		},
	}

	mid := middleware.Middleware{
		Token: cfg.Token,
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/tasks", mid.Auth(h.AllTasks())).Methods("GET")
	router.HandleFunc("/tasks/{page}", mid.Auth(h.AllTasks())).Methods("GET")
	router.HandleFunc("/run", mid.Auth(h.NewRun())).Methods("POST")
	router.HandleFunc("/tasks", mid.Auth(h.NewTask())).Methods("POST")
	router.HandleFunc("/task/{id}", mid.Auth(h.DeleteTask())).Methods("DELETE")
	router.HandleFunc("/task/{id}", mid.Auth(h.ViewTask())).Methods("GET")
	router.HandleFunc("/task/{id}", mid.Auth(h.UpdateTask())).Methods("PUT")

	server := &Server{
		Configuration: *cfg,
		DB:            conn,
		Logger:        l.Named("server"),
	}

	server.serve(router)
}

func (s Server) serve(router *mux.Router) {
	if s.Configuration.SSLMode == "development" {
		s.Logger.Info("Starting server https://" + s.Configuration.Host + ":" + s.Configuration.Port)
		s.Logger.Fatal(
			http.ListenAndServeTLS(
				":"+s.Configuration.Port,
				"certs/cert.pem",
				"certs/key.pem",
				handlers.LoggingHandler(os.Stdout, router)).Error(),
		)
	}

	if s.Configuration.SSLMode == "production" {
		httpsSrv := &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      router,
		}

		dataDir := "certs/"
		hostPolicy := func(ctx context.Context, host string) error {
			allowedHost := s.Configuration.Host
			if host == allowedHost {
				return nil
			}

			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}

		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}

		httpsSrv.Addr = s.Configuration.Host + ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		httpSrv := &http.Server{
			Addr:         ":8080",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      &http.ServeMux{},
		}

		if m != nil {
			httpSrv.Handler = m.HTTPHandler(httpSrv.Handler)
		}

		go func() {
			s.Logger.Info("Starting server http://localhost")

			err := httpSrv.ListenAndServe()
			if err != nil {
				s.Logger.Fatal("httpSrv.ListenAndServe() failed with " + err.Error())
			}
		}()

		s.Logger.Info("Starting server https://" + s.Configuration.Host + ":" + s.Configuration.Port)
		s.Logger.Fatal(httpsSrv.ListenAndServeTLS("", "").Error())
	}
}
