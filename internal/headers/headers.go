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

	if strings.HasPrefix(string(data), "\r\n") {
		return 0, true, nil
	}

	s := strings.Trim(string(data), "\r\n ")
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

	h[strings.ToLower(string(key))] = parts[1]

	return len(data) - 2, false, nil
}
