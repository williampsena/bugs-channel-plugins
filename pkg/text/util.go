package text

import "unicode/utf8"

// This function truncates string by length
func Truncate(content string, length int) string {
	if len(content) < length {
		return content
	}

	if utf8.ValidString(content[:length]) {
		return content[:length]
	}
	return content[:length+1]
}
