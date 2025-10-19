package response

import (
	"fmt"
	"io"

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
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "alive")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			fmt.Printf("err %v", err)
			fmt.Println("write err here")
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	return err
}
