package main

import (
	"crypto/x509"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello World!")

	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("could not get system cert pool: %s", err)
	}

	fmt.Printf("Subjects: \n%s\n", certPool.Subjects())
}
