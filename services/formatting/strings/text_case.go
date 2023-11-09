package strings

import (
	"bytes"
	"strings"
)

// CamelToSnakeCase converts a CamelCase string to a snake_case string
func CamelToSnakeCase(camel string) string {
	var buf bytes.Buffer
	for _, c := range camel {
		if 'A' <= c && c <= 'Z' {
			// just convert [A-Z] to _[a-z]
			if buf.Len() > 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(c - 'A' + 'a')
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}

// CamelToDashedCase converts a CamelCase string to a dashed-case string
func CamelToDashedCase(camel string) string {
	var buf bytes.Buffer
	for _, c := range camel {
		if 'A' <= c && c <= 'Z' {
			// just convert [A-Z] to -[a-z]
			if buf.Len() > 0 {
				buf.WriteRune('-')
			}
			buf.WriteRune(c - 'A' + 'a')
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}

// SnakeToDashedCase replace all underscores with dashes.
func SnakeToDashedCase(dashed string) string {
	return strings.ReplaceAll(dashed, "_", "-")
}

// ToDashedCase converts any string to a dashed-case string
func ToDashedCase(str string) string {
	var lastIsDash bool
	var buf bytes.Buffer
	for _, c := range str {
		if 'A' <= c && c <= 'Z' {
			// just convert [A-Z] to _[a-z]
			if buf.Len() > 0 && !lastIsDash {
				buf.WriteRune('-')
			}
			buf.WriteRune(c - 'A' + 'a')
			lastIsDash = false
		} else if 'a' <= c && c <= 'z' {
			buf.WriteRune(c)
			lastIsDash = false
		} else {
			if !lastIsDash {
				buf.WriteRune('-')
				lastIsDash = true
			}
		}
	}
	return buf.String()
}
