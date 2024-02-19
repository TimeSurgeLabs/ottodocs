package utils

import "regexp"

// uses a regex to remove any quotes from a string
// and returns the modified string. Does nothing if the
// string doesn't contain any quotes.
func RemoveQuotes(s string) string {
	// regex to match any quotes
	regex := `["']`
	// replace any quotes with an empty string
	return regexp.MustCompile(regex).ReplaceAllString(s, "")
}
