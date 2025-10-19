package response

import (
	"fmt"

	"github.com/fenetikm/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

const CRLF = "\r\n"

type Writer struct {
	Response []byte
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case StatusCodeOK:
		w.Response = append(w.Response, []byte("HTTP/1.1 200 OK"+CRLF)...)
	case StatusCodeBadRequest:
		w.Response = append(w.Response, []byte("HTTP/1.1 400 Bad Request"+CRLF)...)
	case StatusCodeInternalServerError:
		w.Response = append(w.Response, []byte("HTTP/1.1 500 Internal Server Error"+CRLF)...)
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

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		w.Response = append(w.Response, []byte(fmt.Sprintf("%s: %s\r\n", k, v))...)
	}
	w.Response = append(w.Response, []byte("\r\n")...)

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	w.Response = append(w.Response, p...)
	return len(p), nil
}
