package git

import "strings"

func Diff() (string, error) {
	return git("diff")
}

func Status() (string, error) {
	return git("status")
}

func IsDirty() (bool, error) {
	status, err := Status()
	if err != nil {
		return false, err
	}
	return !strings.Contains(status, "nothing to commit, working tree clean"), nil
}
