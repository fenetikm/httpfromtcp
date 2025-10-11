package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type requestState int

const (
	requestStateInitialised requestState = iota
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	State       requestState
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == requestStateInitialised {
		rline, n, err := parseRequestLine(string(data))
		if err != nil {
			return 0, fmt.Errorf("Error trying to parse line")
		}
		// Needs more data
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = rline
		r.State = requestStateDone
		return n, nil
	}
	if r.State == requestStateDone {
		return 0, fmt.Errorf("Error state is done")
	}

	return 0, fmt.Errorf("Error unknown state")
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const regnurse = "\r\n"
const bufsize = 8

func parseRequestLine(line string) (RequestLine, int, error) {
	if !strings.Contains(line, regnurse) {
		return RequestLine{}, 0, nil
	}

	rline := strings.Split(line, regnurse)

	parts := strings.Split(rline[0], " ")
	if len(parts) != 3 {
		fmt.Printf("Error parsing request line: %s\n", line)
		return RequestLine{}, 0, errors.New("Request line error")
	}

	if strings.ToUpper(parts[0]) != parts[0] {
		fmt.Println("Error, request line method is not all uppercase")
		return RequestLine{}, 0, errors.New("Request line error")
	}

	version := strings.Split(parts[2], "/")
	if version[1] != "1.1" {
		fmt.Println("Error, request line version is not 1.1")
		return RequestLine{}, 0, errors.New("Request line error")
	}

	return RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   version[1],
	}, len(rline[0]) + 2, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufsize, bufsize)
	// Where we have read up to in buffer
	readUpTo := 0
	req := &Request{
		State: requestStateInitialised,
	}

	// What this does:
	// - if buffer doesn't have enough space for next read
	//   - make a newbuf, double the size
	//   - copy buf into newbuf
	//   - set buf to point to newbuf
	// - read into the buf, offset by readUpTo
	// - if EOF, set to done
	// - inc readUpTo by the number of bytes read
	// - try to parse the buf, sliced up to readUpTo
	// - if we parsed then num of bytes is non-zero
	// - copy just those bytes into the the buf (if will just be the request then, nothing else)
	// - set readUpTo back the num of bytes parsed (?seems redundant?)
	for req.State != requestStateDone {
		if readUpTo >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readUpTo:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.State = requestStateDone
				break
			}
			return nil, err
		}
		readUpTo += numBytesRead

		numBytesParsed, err := req.parse(buf[:readUpTo])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readUpTo -= numBytesParsed
	}

	return req, nil
}
