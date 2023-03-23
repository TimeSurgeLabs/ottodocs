package utils

import (
	"path/filepath"

	"github.com/chand1012/git2gpt/prompt"
)

// GetRepo returns a GitRepo object for the given repoPath
func GetRepo(repoPath, ignoreFilePath string, ignoreGitIgnore bool) (*prompt.GitRepo, error) {
	ignoreList := prompt.GenerateIgnoreList(repoPath, ignoreFilePath, !ignoreGitIgnore)
	ignoreList = append(ignoreList, filepath.Join(repoPath, ".gptignore"))
	repo, err := prompt.ProcessGitRepo(repoPath, ignoreList)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
