package git

import (
	"os"
	"path/filepath"
)

// checks to see if a .git directory exists in the path
func IsGitRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}
