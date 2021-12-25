package server

import (
	"cmd/config"
	"fmt"
	"log"
)

func Exec() {
	config := config.GetConfig()

	// Launch HTTP server
	go func() {
		fmt.Println("Starting server http://localhost")

		err := httpSrv.ListenAndServe()
		if err != nil {
			log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
		}

	}()

	// Launch HTTPS server
	fmt.Println("Starting server https://" + config.Host + ":" + config.Port)
	log.Fatal(httpsSrv.ListenAndServeTLS("", ""))
}
