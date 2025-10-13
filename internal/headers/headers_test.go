package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewHeaders() Headers {
	return make(Headers)
}

func TestParseHeaders(t *testing.T) {
	// Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("   Host: localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 29, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Valid two headers with existing headers
	headers = NewHeaders()
	headers["bob"] = "hey"
	data = []byte("   Host: localhost:42069   \r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 27, n)
	assert.False(t, done)

	assert.Equal(t, "hey", headers["bob"])

	data = []byte("   Content-Type: application/json\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "application/json", headers["content-type"])
	assert.Equal(t, 33, n)
	assert.False(t, done)

	// Valid done
	headers = NewHeaders()
	data = []byte("\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, 0, n)

	// Invalid character
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Valid multiple headers, same key
	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	data = []byte("Set-Person: prime-loves-zig\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	data = []byte("Set-Person: tj-loves-ocaml\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "lane-loves-go,prime-loves-zig,tj-loves-ocaml", headers["set-person"])
	assert.False(t, done)
}
