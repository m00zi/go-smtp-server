package main

import (
	"flag"
	"github.com/stretchkennedy/go-smtp-server"
	"log"
)

func main() {
	addrPtr := flag.String("addr", ":2500", "a TCP address to bind to")
	certFilePtr := flag.String("cert-file", "", "a certificate")
	keyFilePtr := flag.String("key-file", "", "a private key file")
	flag.Parse()

	server := &smtp.Server{
		Addr:     *addrPtr,
		CertFile: *certFilePtr,
		KeyFile:  *keyFilePtr,
		Hostname: "smtp.example.org",
		Handler: func(message smtp.Message) error {
			log.Print(message.Content.Body)
			return nil
		},
	}

	server.Listen()
}
