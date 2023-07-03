package gh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
)

func OpenPullRequest(data map[string]string, owner string, repo string, conf *config.Config) (int, error) {
	if conf.GHToken == "" {
		return -1, fmt.Errorf("no GitHub token found")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)

	payload, err := json.Marshal(data)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.GHToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return -1, fmt.Errorf("failed to create pull request: %s", resp.Status)
	}

	// Extract the pull request number from the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var pr struct {
		Number int `json:"number"`
	}
	err = json.Unmarshal(body, &pr)
	if err != nil {
		return -1, err
	}

	// Return the pull request number and nil error
	return pr.Number, nil
}

func SubmitPullRequestReview(owner, repo string, pullRequestNumber int, review *PullRequestReview, conf *config.Config) error {
	if conf.GHToken == "" {
		return fmt.Errorf("no GitHub token found")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/reviews", owner, repo, pullRequestNumber)

	payload, err := json.Marshal(review)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.GHToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to submit review: %s", resp.Status)
	}

	return nil
}

type PullRequestReview struct {
	// The body of the review. This is optional.
	Body string `json:"body,omitempty"`
	// The event to perform on the pull request. This can be one of:
	// "APPROVE", "REQUEST_CHANGES", "COMMENT", or "DISMISS".
	Event string `json:"event,omitempty"`
	// An array of comments to add to the review. This is optional.
	Comments []*ReviewComment `json:"comments,omitempty"`
	// The commit ID of the pull request. This is required for reviews on multi-commit pull requests.
	CommitID string `json:"commit_id,omitempty"`
	// The path to the file being commented on. This is required for single-file reviews.
	Path string `json:"path,omitempty"`
	// The position in the diff to comment on. This is required for single-file reviews.
	Position int `json:"position,omitempty"`
	// The ID of an existing review to update. This is optional.
	ReviewID int `json:"review_id,omitempty"`
	// Set to true to submit the review and mark the pull request as reviewed. This is optional.
	Submit *bool `json:"submit,omitempty"`
	// Set to true to request changes to the pull request. This is optional.
	RequestChanges *bool `json:"request_changes,omitempty"`
}

type ReviewComment struct {
	// The body of the comment.
	Body string `json:"body"`
	// The relative path to the file that this comment applies to.
	Path string `json:"path"`
	// The line index in the diff to which the comment applies. Note that this is not a line number in the file itself.
	Position int `json:"position"`
}
