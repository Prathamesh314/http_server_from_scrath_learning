package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/Prathamesh314/http_server_from_scrath_learning/internal/headers"
)

var crlf = []byte("\r\n")

type ParseState int

const (
	StateInit ParseState = iota
	StateHeaders
	StateBody
	StateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers // Direct access for easier testing
	State       ParseState
	Body        []byte
}

func newRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(), // Initialize headers
	}
}

func (r *Request) isDone() bool {
	return r.State == StateDone
}

func (r *Request) validateBody() (bool, error) {
	contentLengthStr := r.Headers.GET("Content-Length")
	if contentLengthStr == "" {
		return false, fmt.Errorf("content length is not present in headers")
	}
	var contentLength int
	_, err := fmt.Sscanf(contentLengthStr, "%d", &contentLength)
	if err != nil {
		return false, fmt.Errorf("invalid Content-Length: %v", err)
	}
	if r.isDone() {
		if len(r.Body) == contentLength {
			return true, nil
		}
		return false, fmt.Errorf("body length does not match Content-Length")
	}
	return false, nil
}

func (r *Request) parseBody(data []byte) (int, bool, error) {
	// Get the expected content length
	contentLengthStr := r.Headers.GET("Content-Length")
	var contentLength int
	_, err := fmt.Sscanf(contentLengthStr, "%d", &contentLength)
	if err != nil {
		return 0, false, fmt.Errorf("invalid Content-Length: %v", err)
	}

	// Calculate how much more body data we need
	remaining := contentLength - len(r.Body)

	// If we already have all the data we need
	if remaining == 0 {
		return 0, true, nil
	}

	// Take only what we need from the available data
	toConsume := len(data)
	if toConsume > remaining {
		toConsume = remaining
	}

	// Append the body data
	r.Body = append(r.Body, data[:toConsume]...)

	// Check if we're done or if we've exceeded the limit
	if len(r.Body) > contentLength {
		return 0, false, fmt.Errorf("body length (%d) exceeds Content-Length (%d)", len(r.Body), contentLength)
	}

	// Check if we've received all the body data
	done := len(r.Body) == contentLength

	return toConsume, done, nil
}
// parseRequestLine parses a single request-line from the head of data.
// Returns (requestLine, bytesConsumedIncludingCRLF, err).
// If no full line (no CRLF) is present yet, returns (zeroValue, 0, nil).
func parseRequestLine(data []byte) (RequestLine, int, error) {
	i := bytes.Index(data, crlf)
	if i == -1 { // need more data
		return RequestLine{}, 0, nil
	}
	line := string(data[:i]) // exclude the CRLF

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, fmt.Errorf("malformed request line: %q", line)
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]
	if !strings.HasPrefix(version, "HTTP/") || len(version) < len("HTTP/1.0") {
		return RequestLine{}, 0, fmt.Errorf("invalid http version: %q", version)
	}

	rl := RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   strings.TrimPrefix(version, "HTTP/"),
	}
	// bytes consumed includes the trailing CRLF
	return rl, i + len(crlf), nil
}

// parseSingle handles parsing for a single state transition
// Returns (bytesConsumed, err). If it needs more data, returns (0, nil).
func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case StateInit:
		rl, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil // need more bytes
		}
		r.RequestLine = rl
		r.State = StateHeaders
		return n, nil

	case StateHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil // need more data
		}
		if done {
			if found := r.Headers.GET("Content-Length"); len(found) == 0 {
				// If there is no content lenght then we directly move on to StateDone
				fmt.Println("Moving to StateDonoe")
				r.State = StateDone
			}else {
				// We will parse the body
				fmt.Println("Moving to StateBody")
				r.State = StateBody
			}
		}
		return n, nil

	case StateBody:
		// append all the data to .Body field
		n, done, err := r.parseBody(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			// We need more data
			return 0, nil
		}

		if done {
			r.State = StateDone
		}

		return n, nil

	case StateDone:
		return 0, nil

	default:
		return 0, fmt.Errorf("unknown parser state")
	}
}

// parse consumes bytes into the Request based on the current state.
// Returns (bytesConsumed, err). If it needs more data, returns (0, nil).
func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.State != StateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, err
		}
		if n == 0 {
			// Need more data to continue parsing
			return totalBytesParsed, nil
		}
		totalBytesParsed += n

		// If we've consumed all available data, break
		if totalBytesParsed >= len(data) {
			break
		}
	}

	return totalBytesParsed, nil
}

// RequestFromReader streams bytes from reader, parses incrementally,
// and returns once the request is fully parsed.
func RequestFromReader(reader io.Reader) (*Request, error) {
	req := newRequest()

	// rolling buffer we append to and trim from the front
	buf := make([]byte, 0, 8)
	tmp := make([]byte, 8)

	for {
		// Try to parse whatever we already have buffered
		for {
			if len(buf) == 0 {
				break
			}
			consumed, err := req.parse(buf)
			if err != nil {
				return nil, err
			}
			if consumed == 0 {
				break // need more data to progress
			}
			// drop consumed prefix
			buf = buf[consumed:]

			if req.isDone() {
				return req, nil
			}
		}

		// Need more data; read a chunk
		n, err := reader.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
			continue
		}
		if err == io.EOF {
			// Reader ended; if not done, we didn't get a full request
			if req.isDone() {
				return req, nil
			}
			return nil, io.ErrUnexpectedEOF
		}
		if err != nil {
			return nil, err
		}
	}
}
