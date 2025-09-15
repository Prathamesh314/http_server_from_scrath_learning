package request

import (
	"fmt"
	"io"
	"strings"
)

const SEPARATOR = "\r\n"
const HTTP_VERSION = "1.1"

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

func parseRequestLine(request_line string) (*Request, error) {
	// for now we have to discard everything after the requestline
	// we know that the input string will look like this
	/*
		(request-lines) CRLF
		*(field lines) CRLf
		(body)
	*/

	parts := strings.Split(request_line, SEPARATOR)[0]
	/*
	In request line we have three parts
	1. method
	2. target
	3. http/httpversion
	*/
	
	request_line_parts := strings.Split(parts, " ")
	method := request_line_parts[0]
	target := request_line_parts[1]
	http_parts := strings.Split(string(request_line_parts[2]), "/")
	if len(http_parts) != 2 {
		return &Request{}, fmt.Errorf("Wrong http version: %v\n", http_parts)
	}

	http_version := http_parts[1]

	if http_version != HTTP_VERSION {
		return nil, fmt.Errorf("Wrong http version. Expecting: %s", HTTP_VERSION)
	}

	return &Request{
		RequestLine{
			Method: string(method),
			HttpVersion: http_version,
			RequestTarget: string(target),
		},
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	request_string := string(request)
	fmt.Printf("Request: %s\n", request_string)
	return parseRequestLine(request_string)
}