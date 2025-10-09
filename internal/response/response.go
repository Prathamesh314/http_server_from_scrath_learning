package response

import (
	"fmt"
	"io"

	"github.com/Prathamesh314/http_server_from_scrath_learning/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var mapping = map[StatusCode]string{
		StatusCodeOK:                  "HTTP/1.1 200 OK",
		StatusCodeBadRequest:          "HTTP/1.1 400 Bad Request",
		StatusCodeInternalServerError: "HTTP/1.1 500 Internal Server Error",
	}
	statusLine, ok := mapping[statusCode]
	if !ok {
		statusLine = fmt.Sprintf("HTTP/1.1 %d Unknown Status", statusCode)
	}
	_, err := w.Write([]byte(statusLine + "\r\n"))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.SET("Content-Length", fmt.Sprintf("%d\n", contentLen))
	headers.SET("Connection", "close")
	headers.SET("Content-Type", "text/plain")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		headerLine := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Write([]byte(headerLine))
		if err != nil {
			return err
		}
	}
	// Write the blank line that separates headers from body
	_, err := w.Write([]byte("\r\n"))
	return err
}