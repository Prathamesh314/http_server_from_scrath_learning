package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

var SEPARATOR = []byte("\r\n")

func NewHeaders() (Headers){
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	header_parts := bytes.Split(data, SEPARATOR)
	headers := header_parts[0] // This will look like: "Host: localhost:42069")

	idx := bytes.Index(headers, []byte(":"))
	if idx == -1 {
		return 0, false, fmt.Errorf("cannot find colon ':'")
	}

	field_name := string(headers[:idx])

	// this is an invalid field name:= HOST : localhost:42069
	if strings.HasSuffix(field_name, " ") {
		return 0, false, fmt.Errorf("invalid field_name. contains OWS in between colon")
	}

	field_name = strings.Trim(field_name, " ")

	h[field_name] = strings.Trim(string(headers)[idx+1:], " ")

	return len(string(headers)) + len(SEPARATOR), true, nil
}