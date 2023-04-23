package git

func GetBranch() (string, error) {
	return git("rev-parse", "--abbrev-ref", "HEAD")
}
