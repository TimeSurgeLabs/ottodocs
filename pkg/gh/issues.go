package gh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
)

func GetIssue(owner, repo string, issueNumber int, conf *config.Config) (*IssueWithComments, error) {
	if conf.GHToken == "" {
		return nil, fmt.Errorf("no GitHub token found")
	}

	issueURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", owner, repo, issueNumber)
	commentsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments", owner, repo, issueNumber)

	issueReq, err := http.NewRequest("GET", issueURL, nil)
	if err != nil {
		return nil, err
	}

	commentsReq, err := http.NewRequest("GET", commentsURL, nil)
	if err != nil {
		return nil, err
	}

	issueReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.GHToken))
	commentsReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.GHToken))

	client := &http.Client{}

	issueResp, err := client.Do(issueReq)
	if err != nil {
		return nil, err
	}
	defer issueResp.Body.Close()

	commentsResp, err := client.Do(commentsReq)
	if err != nil {
		return nil, err
	}
	defer commentsResp.Body.Close()

	if issueResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve issue: %s", issueResp.Status)
	}

	if commentsResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve issue comments: %s", commentsResp.Status)
	}

	var issue Issue
	err = json.NewDecoder(issueResp.Body).Decode(&issue)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	err = json.NewDecoder(commentsResp.Body).Decode(&comments)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		comment := &comments[i]
		if comment.Body != "" {
			// Extract the username from the comment URL
			if strings.HasPrefix(comment.Body, "https://github.com/") {
				usernameStart := strings.Index(comment.Body, "/")
				if usernameStart != -1 {
					usernameEnd := strings.Index(comment.Body[usernameStart+1:], "/")
					if usernameEnd != -1 {
						username := comment.Body[usernameStart+1 : usernameStart+1+usernameEnd]
						comment.Username = username
					}
				}
			}
		}
	}

	issueWithComments := &IssueWithComments{
		Issue:    issue,
		Comments: comments,
	}

	return issueWithComments, nil
}

type Issue struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type IssueWithComments struct {
	Issue    Issue     `json:"issue"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	// The username of the user who made the comment. This is optional.
	Username string `json:"username,omitempty"`
}
