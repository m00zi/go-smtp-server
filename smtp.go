package smtp

import (
	"crypto/tls"
	"io"
	"log"
	"net"

	"github.com/mailhog/data"
)

// A server defines the parameters for running an SMTP server.
type Server struct {
	Addr     string  // Address to listen on
	CertFile string  // Certificate file to load
	KeyFile  string  // Key file to load
	Hostname string  // Hostname to report in SMTP, defaults to Addr
	Handler  Handler // Handler
}

type Handler func(Message) error

// Alias of mailhog's Message
type Message *data.Message

// Listen starts a server on the TCP network address server.Addr and then
// reads emails, using server.Hostname within SMTP, and calling
// server.Handler for every email successfully recieved. Listen always
// returns a non-nil error.
func (server Server) Listen() {
	addr := server.Addr
	if addr == "" {
		addr = ":2500"
	}

	hostname := server.Hostname
	if hostname == "" {
		hostname = addr
	}

	handler := server.Handler
	if handler == nil {
		handler = func(m Message) error {
			return nil
		}
	}

	log.Printf("[SMTP] Binding to address: %s\n", addr)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[SMTP] Error listening on socket: %s\n", err)
	}
	defer ln.Close()

	var tlsConfig *tls.Config
	if server.CertFile != "" && server.KeyFile != "" {
		certificate, err := tls.LoadX509KeyPair(server.CertFile, server.KeyFile)
		if err != nil {
			log.Printf("[SMTP] Failed to load cert: %s\n", err)

		} else {
			tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{certificate},
			}
			log.Printf("[SMTP] Loaded cert: %s\n", server.CertFile)
		}
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[SMTP] Error accepting connection: %s\n", err)
			continue
		}

		go Accept(
			conn.(*net.TCPConn).RemoteAddr().String(),
			io.ReadWriteCloser(conn),
			handler,
			hostname,
			tlsConfig,
		)
	}
}
