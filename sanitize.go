package pgqb

import "strings"

// Sanitize escapes dangerous characters
func Sanitize(source string) string {
	return strings.Replace(source, "'", "''", -1)
}

// Quote sanitize and quote the given string
func Quote(source string) string {
	return "'" + Sanitize(source) + "'"
}

// DoubleQuote sanitize and wrap the given string in double quotes
func DoubleQuote(source string) string {
	return `"` + Sanitize(source) + `"`
}
