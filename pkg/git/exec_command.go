package git

import (
	"os/exec"
	"strings"
)

func git(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	// output
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
