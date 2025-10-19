package headers

import (
	"fmt"
	"strconv"
	"strings"
)

type Headers map[string]string

const validKeyChars = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890!#$%&'*+-.^_`|~"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if !strings.Contains(string(data), "\r\n") {
		return 0, false, nil
	}

	// End of headers
	if strings.HasPrefix(string(data), "\r\n") {
		return 2, true, nil
	}

	sh := strings.Split(string(data), "\r\n")
	s := strings.Trim(sh[0], " ")
	parts := strings.Split(s, ": ")
	if len(parts) != 2 {
		return 0, false, fmt.Errorf("Bad header")
	}
	key := []byte(parts[0])
	if string(key[len(key)-1]) == " " {
		return 0, false, fmt.Errorf("Bad header, string left of colon")
	}
	for _, b := range key {
		if !strings.Contains(validKeyChars, string(b)) {
			return 0, false, fmt.Errorf("Bad header key, bad character")
		}
	}

	val := parts[1]

	found := false
	for k, v := range h {
		if strings.EqualFold(k, string(key)) {
			h[k] = v + "," + val
			found = true
			break
		}
	}
	if !found {
		h[string(key)] = val
	}

	return len(sh[0]) + 2, false, nil
}

func (h Headers) Get(key string) string {
	for k, v := range h {
		if strings.EqualFold(k, key) {
			return v
		}
	}

	return ""
}

func (h Headers) ContentLength() (int, error) {
	cl := h.Get("Content-Length")
	if cl == "" {
		return 0, nil
	}

	cli, err := strconv.Atoi(cl)
	if err != nil {
		return 0, fmt.Errorf("Non numeric Content-Length")
	}

	return cli, nil
}
