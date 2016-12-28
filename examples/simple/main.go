package main

import (
	"flag"
	"github.com/stretchkennedy/go-smtp-server"
	"log"
)

func main() {
	addrPtr := flag.String("addr", ":2500", "a TCP address to bind to")
	flag.Parse()

	server := &smtp.Server{
		Addr:     *addrPtr,
		Hostname: "smtp.example.org",
		Handler: func(message smtp.Message) error {
			log.Print(message.Content.Body)
			return nil
		},
	}

	log.Fatal(server.Listen())
}
