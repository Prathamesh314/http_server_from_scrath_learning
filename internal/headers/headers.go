package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

var SEPARATOR = []byte("\r\n")
var VALID_SPECIAL_CHARACTERS = strings.Split("!, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~", ",")

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) GET(field_name string) string {
	return h[strings.ToLower(field_name)]
}

func (h Headers) SET(field_name, field_value string) {
	key := strings.ToLower(field_name)
	value := h[key]
	if len(value) == 0 {
		value = strings.TrimSpace(field_value)
	} else {
		value = value + ", " + strings.TrimSpace(field_value)
	}
	h[key] = value
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

// parseHeaderLine parses a single header line from the head of data.
// Returns (bytesConsumed, done, err).
// If no full line is present yet, returns (0, false, nil).
// If empty line (end of headers) is found, returns (bytesConsumed, true, nil).
func (h Headers) parseHeaderLine(data []byte) (int, bool, error) {
	i := bytes.Index(data, SEPARATOR)
	if i == -1 {
		// Need more data - no complete line found
		return 0, false, nil
	}

	line := data[:i]
	
	// Empty line means end of headers
	if len(line) == 0 {
		return i + len(SEPARATOR), true, nil
	}

	// Find the colon separator
	colonIdx := bytes.Index(line, []byte(":"))
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("malformed header: missing colon")
	}

	fieldName := string(line[:colonIdx])
	fieldValue := string(line[colonIdx+1:])

	// Check for invalid spacing around colon
	if strings.HasSuffix(fieldName, " ") {
		return 0, false, fmt.Errorf("invalid field_name. contains OWS in between colon")
	}

	fieldName = strings.TrimSpace(fieldName)
	if !h.isValidFieldName(fieldName) {
		return 0, false, fmt.Errorf("invalid character")
	}

	h.SET(fieldName, fieldValue)

	return i + len(SEPARATOR), false, nil
}

// Parse parses headers from streaming data.
// Returns (bytesConsumed, done, err).
// done=true when all headers are parsed (empty line found).
// If it needs more data, returns (0, false, nil).
func (h Headers) Parse(data []byte) (int, bool, error) {
	totalConsumed := 0
	
	for {
		consumed, done, err := h.parseHeaderLine(data[totalConsumed:])
		if err != nil {
			return 0, false, err
		}
		if consumed == 0 {
			// Need more data
			return totalConsumed, false, nil
		}
		
		totalConsumed += consumed
		
		if done {
			// Found empty line - end of headers
			return totalConsumed, true, nil
		}
		
		// If we've consumed all available data, break
		if totalConsumed >= len(data) {
			break
		}
	}
	
	return totalConsumed, false, nil
}