package git

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func GetRemote(remote string) (string, error) {
	return git("remote", "get-url", remote)
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

func ExtractOriginInfo(gitURL string) (string, string, error) {
	// Regular expression to match both SSH and HTTPS URLs
	// This pattern supports optional ".git" at the end and optional "https://" or "git@" at the beginning
	regexPattern := `^(?:(?:https:\/\/|git@)github\.com(?:\/|:))?(\w+)\/(\w+)(?:\.git)?$`
	regex := regexp.MustCompile(regexPattern)

	matches := regex.FindStringSubmatch(gitURL)

	if len(matches) == 3 {
		owner := matches[1]
		repo := matches[2]
		return owner, repo, nil
	}

	return "", "", fmt.Errorf("unable to parse owner and repo from git URL")
}
