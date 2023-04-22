package git

// git add -A
func AddAll() (string, error) {
	return git("add", "-A")
}

// git commit -am <message>
func Commit(message string) (string, error) {
	return git("commit", "-am", message)
}

// git push
func Push() (string, error) {
	return git("push")
}
