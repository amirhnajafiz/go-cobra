package main

import (
	"cmd/pkg/encrypt"
	"fmt"
	"os"
)

func main() {
	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		encrypt.GenerateCertificateAuthority()
		encrypt.GenerateCert()
	}
	fmt.Println("Lets go")
}
