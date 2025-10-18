package server

import (
	"io"
	"log"
	"net"
	"strings"
)

type Server struct {
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		return nil, err
	}

	return &Server{
		listener: l,
	}, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) listen() {
	for {
		// Wait for a connection.
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			s.handle(c)
			c.Close()
		}(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	resp := `HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 13

Hello World!`
	reader := strings.NewReader(resp)
	io.Copy(conn, reader)
}
