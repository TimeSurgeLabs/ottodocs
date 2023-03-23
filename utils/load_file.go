package utils

import (
	"errors"
	"os"
	"unicode/utf8"
)

// Loads a file and return an error if its not valid UTF-8
func LoadFile(filePath string) (string, error) {
	// load the file
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	if !utf8.Valid(contents) {
		return "", errors.New("the file is not valid UTF-8")
	}

	return string(contents), nil
}
