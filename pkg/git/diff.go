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

func GetChangedFiles() ([]string, error) {
	resp, err := git("diff", "--name-only")
	if err != nil {
		return nil, err
	}

	return strings.Split(resp, "\n"), nil
}

func GetBranchDiff(base, head string) (string, error) {
	return git("diff", base+".."+head)
}

func GetChangedFilesBranches(base, head string) ([]string, error) {
	resp, err := git("diff", "--name-only", base+".."+head)
	if err != nil {
		return nil, err
	}

	return strings.Split(resp, "\n"), nil
}

func Log() (string, error) {
	return git("log", "--oneline")
}

func LogBetween(base, head string) (string, error) {
	return git("log", "--oneline", base+".."+head)
}
