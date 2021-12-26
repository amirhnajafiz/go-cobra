package handler

import (
	"cmd/config"
	"cmd/pkg/encrypt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func HandleRequests(configuration config.Config) {
	var httpsSrv *http.Server
	var httpSrv *http.Server
	var m *autocert.Manager

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
}
