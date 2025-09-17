package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	return n, nil
}

func Test_RequestLine_SimpleGET(t *testing.T) {
	reqStr := "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"

	for chunk := 1; chunk <= len(reqStr); chunk *= 2 { // 1,2,4,8,...
		reader := &chunkReader{data: reqStr, numBytesPerRead: chunk}
		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	}
}

func Test_RequestLine_PathGET(t *testing.T) {
	reqStr := "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"

	for _, chunk := range []int{1, 3, 5, len(reqStr)} {
		reader := &chunkReader{data: reqStr, numBytesPerRead: chunk}
		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	}
}

// Optional: proves we return "need more data" behavior internally (no panic)
func Test_RequestLine_Incomplete_NoCRLF(t *testing.T) {
	reader := &chunkReader{data: "GE", numBytesPerRead: 1}
	_, err := RequestFromReader(reader)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}
