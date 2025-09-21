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
	assert.Equal(t, "localhost:42069", headers.GET("HOST"))
	assert.Equal(t, 23, n)
	assert.True(t, done)

	// // Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// // Test: Case insensitivity - capital letters in header keys
	headers = NewHeaders()
	data = []byte("HOST: localhost:42069\r\nContent-TYPE: application/json\r\nUSER-Agent: curl/7.81.0\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	// Header keys should be normalized to lowercase
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "application/json", headers["content-type"])
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	// Original case should NOT be found
	assert.Empty(t, headers["HOST"])
	assert.Empty(t, headers["Content-TYPE"])
	assert.Empty(t, headers["USER-Agent"])
	assert.True(t, done)

	// // Test: Invalid character in header key
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	assert.Contains(t, err.Error(), "invalid character") // or whatever error message your parser uses
}