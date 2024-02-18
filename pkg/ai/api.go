package ai

import (
	"encoding/json"

	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/constants"
)

func APIDocs(files []string, conf *config.Config) (string, error) {
	// join all the files into a single string

	fileStr := ""
	for _, file := range files {
		fileStr += file
	}

	resp, err := request(constants.API_DOCS_PROMPT, fileStr, conf)
	if err != nil {
		return "", err
	}

	return resp, nil
}

type endpointsResp struct {
	Endpoints []string `json:"endpoints"`
}

func APIEndpoints(files []string, conf *config.Config) ([]string, error) {
	// ask the AI to list all the endpoints in the files
	resp, err := request(constants.API_ENDPOINTS_PROMPT, "", conf)
	if err != nil {
		return nil, err
	}

	// parse the response
	var endpoints endpointsResp
	err = json.Unmarshal([]byte(resp), &endpoints)
	if err != nil {
		return nil, err
	}

	return endpoints.Endpoints, nil
}
