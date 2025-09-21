## Added multiple values in headers
```
package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

var SEPARATOR = []byte("\r\n")
var VALID_SPECIAL_CHARACTERS = strings.Split("!, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~", ",")

func NewHeaders() (Headers){
	return Headers{}
}

func (h Headers) GET(field_name string) string {
	return h[strings.ToLower(field_name)]
}

func (h Headers) SET(field_name, field_value string) {
	value := h.GET(field_name)
	if  len(value) == 0 {
		value = strings.Trim(field_value, " ")
	}else {
		value = value + ", " + strings.Trim(field_value, " ")
	}
	
	h[strings.ToLower(field_name)] = value
}

func (h Headers) isValidFieldName(field_name string) bool {
	for _, c := range field_name {
		switch {
		case c >= 'A' && c <= 'Z':
			continue
		case c >= 'a' && c <= 'z':
			continue
		case c >= '0' && c <= '9':
			continue
		default:
			// Check if c is in the list of valid special characters
			found := false
			for _, s := range VALID_SPECIAL_CHARACTERS {
				s = strings.TrimSpace(s)
				if len(s) == 1 && rune(s[0]) == c {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	header_parts := bytes.Split(data, SEPARATOR)
	len_of_headers := 0

	for i:=0; i<int(len(header_parts)); i = i+1 {
		headers := header_parts[i]
		if len(headers) == 0 {
			// because we have reached the end of string.
			break
		}
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
		if !h.isValidFieldName(field_name) {
			return 0, false, fmt.Errorf("invalid character")
		}

		h.SET(field_name, string(headers)[idx+1:])

		len_of_headers += len(headers)
	}

	return len_of_headers + len(SEPARATOR), true, nil
}
```
