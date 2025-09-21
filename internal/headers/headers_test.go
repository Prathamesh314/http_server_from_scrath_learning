package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.GET("Host"))
	assert.Equal(t, 23, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Case insensitivity - capital letters in header keys
	headers = NewHeaders()
	data = []byte("HOST: localhost:42069\r\nContent-TYPE: application/json\r\nUSER-Agent: curl/7.81.0\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	// Header keys should be normalized to lowercase
	assert.Equal(t, "localhost:42069", headers.GET("host"))
	assert.Equal(t, "application/json", headers.GET("Content-Type"))
	assert.Equal(t, "curl/7.81.0", headers.GET("User-Agent"))
	assert.True(t, done)

	// Test: Invalid character in header key
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Contains(t, err.Error(), "invalid character") // or whatever error message your parser uses

	// Test: Multiple values for same header key (starting with existing header)
	headers = NewHeaders()
	// First, add an initial header value
	headers["set-person"] = "lane-loves-go"
	
	// Parse additional headers with the same key
	data = []byte("Set-Person: prime-loves-zig\r\nSet-Person: tj-loves-ocaml\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	// Should combine all values with comma separation
	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
	assert.True(t, done)

	// Test: Multiple values for same header key (no existing header)
	headers = NewHeaders()
	data = []byte("Accept: text/html\r\nAccept: application/json\r\nAccept: */*\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	// Should combine values with comma separation
	assert.Equal(t, "text/html, application/json, */*", headers["accept"])
	assert.True(t, done)
}