package utils

import "os"

func WriteFile(filePath string, contents string) error {
	return os.WriteFile(filePath, []byte(contents), 0644)
}
