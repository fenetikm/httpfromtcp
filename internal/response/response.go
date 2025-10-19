package response

import (
	"fmt"
	"io"
	"log"

	"github.com/fenetikm/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

const CRLF = "\r\n"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error
	switch statusCode {
	case StatusCodeOK:
		_, err = w.Write([]byte("HTTP/1.1 200 OK" + CRLF))
	case StatusCodeBadRequest:
		_, err = w.Write([]byte("HTTP/1.1 400 Bad Request" + CRLF))
	case StatusCodeInternalServerError:
		_, err = w.Write([]byte("HTTP/1.1 500 Internal Server Error" + CRLF))
	}
	if err != nil {
		fmt.Println("err writing status code")
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.Headers{}
	_, _, err := h.Parse(fmt.Appendf([]byte{}, "Content-Length: %d\r\n", contentLen))
	if err != nil {
		log.Fatalf("Couldn't parse content length header.")
	}

	_, _, err = h.Parse([]byte("Connection: alive\r\n"))
	if err != nil {
		log.Fatalf("Couldn't parse content length header.")
	}

	_, _, err = h.Parse([]byte("Content-Type: text/plain\r\n"))
	if err != nil {
		log.Fatalf("Couldn't parse content length header.")
	}

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write(fmt.Appendf([]byte{}, "%s: %s\r\n", k, v))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte(CRLF))
	if err != nil {
		return err
	}

	return nil
}
