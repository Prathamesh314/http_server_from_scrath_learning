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