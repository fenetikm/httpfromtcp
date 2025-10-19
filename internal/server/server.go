package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/fenetikm/httpfromtcp/internal/request"
	"github.com/fenetikm/httpfromtcp/internal/response"
)

type Server struct {
	listener    net.Listener
	closed      atomic.Bool
	handlerFunc Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he *HandlerError) Write(w io.Writer) {
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		log.Fatalf("Couldn't handle writing status line")
	}
	body := he.Message
	cl := len(body)
	heads := response.GetDefaultHeaders(cl)
	response.WriteHeaders(w, heads)
	w.Write([]byte(body))
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := Server{
		listener:    l,
		handlerFunc: handler,
	}
	go s.listen()

	return &s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) listen() {
	for {
		// Wait for a connection.
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		he := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		he.Write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	he := s.handlerFunc(buf, req)
	if he != nil {
		he.Write(conn)
		return
	}

	he = &HandlerError{
		StatusCode: response.StatusCodeOK,
		Message:    buf.String(),
	}
	he.Write(conn)
}
