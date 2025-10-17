package headers

import (
	"fmt"
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

	lk := strings.ToLower(string(key))
	val := parts[1]

	if v, ok := h[lk]; ok {
		h[lk] = v + "," + val
	} else {
		h[lk] = val
	}

	return len(sh[0]) + 2, false, nil
}

func (h Headers) Get(key string) string {
	lk := strings.ToLower(string(key))
	if _, ok := h[lk]; ok {
		return h[lk]
	}

	return ""
}
