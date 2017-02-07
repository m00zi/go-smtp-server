package smtp

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"strings"

	"github.com/mailhog/data"
	"github.com/mailhog/smtp"
)

// session represents a SMTP session using net.TCPConn
type session struct {
	conn          io.ReadWriteCloser
	proto         *smtp.Protocol
	handler       Handler
	remoteAddress string
	isTLS         bool
	line          string
	tlsConfig     *tls.Config
}

// Accept starts a new SMTP session using io.ReadWriteCloser
func accept(
	remoteAddress string,
	conn io.ReadWriteCloser,
	handler Handler,
	hostname string,
	tlsConfig *tls.Config,
) {
	defer conn.Close()

	proto := smtp.NewProtocol()
	proto.Hostname = hostname

	s := &session{conn, proto, handler, remoteAddress, false, "", tlsConfig}
	proto.MessageReceivedHandler = s.acceptMessage
	proto.GetAuthenticationMechanismsHandler = func() []string { return []string{"PLAIN"} }
	if tlsConfig != nil {
		proto.TLSHandler = s.tlsHandler
	}

	s.logf("Starting session")
	s.write(proto.Start())
	for s.read() == true {
	}
	s.logf("Session ended")
}

func (s *session) tlsHandler(done func(ok bool)) (errorReply *smtp.Reply, callback func(), ok bool) {
	s.logf("Returning TLS handler")
	return nil, func() {
		s.logf("Upgrading session to TLS")
		tConn := tls.Server(s.conn.(net.Conn), s.tlsConfig)
		err := tConn.Handshake()
		if err != nil {
			s.logf("handshake error in TLS connection: %s", err)
			done(false)
			return
		}
		s.conn = tConn
		s.isTLS = true
		s.logf("Session upgrade complete")
		done(true)
	}, true
}

func (s *session) acceptMessage(msg *data.SMTPMessage) (id string, err error) {
	m := msg.Parse(s.proto.Hostname)
	s.logf("Storing message %s", m.ID)
	return string(m.ID), s.handler(m)
}

func (s *session) logf(message string, args ...interface{}) {
	message = strings.Join([]string{"[SMTP %s]", message}, " ")
	args = append([]interface{}{s.remoteAddress}, args...)
	log.Printf(message, args...)
}

// read reads from the underlying net.TCPConn
func (s *session) read() bool {
	buf := make([]byte, 1024)
	n, err := s.conn.Read(buf)

	if n == 0 {
		s.logf("Connection closed by remote host\n")
		s.conn.Close() // not sure this is necessary?
		return false
	}

	if err != nil {
		s.logf("Error reading from socket: %s\n", err)
		return false
	}

	text := string(buf[0:n])
	logText := strings.Replace(text, "\n", "\\n", -1)
	logText = strings.Replace(logText, "\r", "\\r", -1)
	s.logf("Received %d bytes: '%s'\n", n, logText)

	s.line += text

	for strings.Contains(s.line, "\r\n") {
		line, reply := s.proto.Parse(s.line)
		s.line = line

		if reply != nil {
			s.write(reply)
			if reply.Status == 221 {
				io.Closer(s.conn).Close()
				return false
			}
		}
	}

	return true
}

// Write writes a reply to the underlying net.TCPConn
func (s *session) write(reply *smtp.Reply) {
	lines := reply.Lines()
	for _, l := range lines {
		logText := strings.Replace(l, "\n", "\\n", -1)
		logText = strings.Replace(logText, "\r", "\\r", -1)
		s.logf("Sent %d bytes: '%s'", len(l), logText)
		s.conn.Write([]byte(l))
	}
	if reply.Done != nil {
		reply.Done()
	}
}
