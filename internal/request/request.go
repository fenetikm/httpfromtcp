package request

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/fenetikm/httpfromtcp/internal/headers"
)

/*
TODO:
- Handle body longer than content-length
*/

type requestState int

const (
	requestStateInitialised requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       requestState
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}

		// Data too short, doesn't contain something that can be parsed
		if n == 0 {
			break
		}

		// Something was successfully parsed
		totalBytesParsed += n
		if totalBytesParsed > len(data) {
			return 0, fmt.Errorf("Too many bytes?")
		}

		if totalBytesParsed == len(data) {
			break
		}
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	if r.state == requestStateInitialised {
		rline, n, err := parseRequestLine(string(data))
		if err != nil {
			return 0, fmt.Errorf("Error trying to parse line")
		}
		// Needs more data
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = rline
		r.state = requestStateParsingHeaders
		r.Headers = headers.Headers{}
		return n, nil
	}
	if r.state == requestStateParsingHeaders {
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			r.state = requestStateParsingBody
			return n, nil
		}

		// More data please
		if !done && n == 0 {
			return n, nil
		}

		return n, nil
	}
	if r.state == requestStateParsingBody {
		cl, err := r.Headers.ContentLength()
		if err != nil {
			return 0, err
		}
		if cl == 0 {
			r.state = requestStateDone
			return 0, nil
		}

		if len(data) < cl {
			return 0, nil
		}

		r.Body = make([]byte, cl)
		copy(r.Body, data[:cl])
		r.state = requestStateDone
		return cl, nil
	}
	if r.state == requestStateDone {
		return 0, fmt.Errorf("Error state is done")
	}

	return 0, fmt.Errorf("Error unknown state")
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const bufsize = 8

func parseRequestLine(line string) (RequestLine, int, error) {
	if !strings.Contains(line, crlf) {
		return RequestLine{}, 0, nil
	}

	rline := strings.Split(line, crlf)

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
	buf := make([]byte, bufsize)
	// Where we have read up to in buffer
	readUpTo := 0
	req := &Request{
		state: requestStateInitialised,
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
	// - copy just those bytes into the the buf (it will just be the request then, nothing else)
	// - set readUpTo back the num of bytes parsed (?seems redundant?)
	for req.state != requestStateDone {
		if readUpTo >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readUpTo:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				cl, err := req.Headers.ContentLength()
				if err != nil {
					return nil, err
				}
				if req.state == requestStateParsingBody && cl != 0 {
					return nil, fmt.Errorf("Not enough bytes to parse body")
				}

				req.state = requestStateDone
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
