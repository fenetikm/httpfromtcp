package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := Server{
		listener: l,
	}
	go s.listen()

	return &s, nil
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

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	resp := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		// "Content-Length: 13\r\n\r\n" +
		"Hello World!\n"
	conn.Write([]byte(resp))
}
