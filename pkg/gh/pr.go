package gh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chand1012/ottodocs/pkg/config"
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
