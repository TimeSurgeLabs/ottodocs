package gh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chand1012/ottodocs/pkg/config"
)

func CreateDraftRelease(owner, repo, title, body, tag string, conf *config.Config) error {
	if conf.GHToken == "" {
		return fmt.Errorf("no GitHub token found")
	}

	releaseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	data := map[string]interface{}{
		"title":      title,
		"body":       body,
		"tag_name":   tag,
		"draft":      true,
		"prerelease": false,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	releaseReq, err := http.NewRequest("POST", releaseURL, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}

	releaseReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.GHToken))
	releaseReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	releaseResp, err := client.Do(releaseReq)
	if err != nil {
		return err
	}
	defer releaseResp.Body.Close()

	if releaseResp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create release draft: %s", releaseResp.Status)
	}

	return nil
}
