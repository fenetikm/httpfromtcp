package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(line string) (RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		fmt.Printf("Error parsing request line: %s\n", line)
		return RequestLine{}, errors.New("Request line error")
	}

	if strings.ToUpper(parts[0]) != parts[0] {
		fmt.Println("Error, request line method is not all uppercase")
		return RequestLine{}, errors.New("Request line error")
	}

	version := strings.Split(parts[2], "/")
	if version[1] != "1.1" {
		fmt.Println("Error, request line version is not 1.1")
		return RequestLine{}, errors.New("Request line error")
	}

	return RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   version[1],
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading all from reader")
		return nil, err
	}

	lines := strings.Split(string(req), "\r\n")
	rline, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: rline,
	}, nil
}
