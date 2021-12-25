package main

import (
	"cmd/pkg/encrypt"
	"fmt"
)

func main() {
	encrypt.GenerateCertificateAuthority()
	encrypt.GenerateCert()
	fmt.Println("Lets go")
}
