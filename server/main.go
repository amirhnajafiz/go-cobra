package server

import (
	"cmd/config"
	"fmt"
)

func Exec() {
	configuration := config.GetConfig()

	// Launch HTTP server
	go func() {
		fmt.Println("Starting server http://localhost")

		//err := httpSrv.ListenAndServe()
		//if err != nil {
		//	log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
		//}

	}()

	// Launch HTTPS server
	fmt.Println("Starting server https://" + configuration.Host + ":" + configuration.Port)
	// log.Fatal(httpsSrv.ListenAndServeTLS("", ""))
}
