package utils

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/chand1012/git2gpt/prompt"
	"github.com/go-git/go-git/v5"
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

func GetOrigin(repoPath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}
	origin, err := repo.Remote("origin")
	if err != nil {
		return "", err
	}
	return origin.Config().URLs[0], nil
}

func OriginToGitHub(origin string) (string, error) {
	parsedURL, err := url.Parse(origin)
	if err != nil {
		return "", fmt.Errorf("failed to parse origin URL: %v", err)
	}

	// Remove .git from the path if it exists
	path := strings.TrimSuffix(parsedURL.Path, ".git")

	// Construct the GitHub link
	if parsedURL.Scheme == "https" {
		return fmt.Sprintf("https://github.com%s", path), nil
	} else if parsedURL.Scheme == "http" {
		return fmt.Sprintf("http://github.com%s", path), nil
	} else if parsedURL.Scheme == "git" || parsedURL.Scheme == "ssh" {
		host := strings.TrimPrefix(parsedURL.Host, "git@")
		return fmt.Sprintf("https://%s%s", host, path), nil
	}

	return "", fmt.Errorf("unsupported protocol: %s", parsedURL.Scheme)
}

func ExtractOriginInfo(origin string) (string, string, error) {
	// Remove .git from the path if it exists
	path := strings.TrimSuffix(origin, ".git")

	// Extract the owner and repo name
	splitPath := strings.Split(path, "/")
	if len(splitPath) < 2 {
		return "", "", fmt.Errorf("invalid origin URL: %s", origin)
	}
	owner := splitPath[len(splitPath)-2]
	repo := splitPath[len(splitPath)-1]

	return owner, repo, nil
}
