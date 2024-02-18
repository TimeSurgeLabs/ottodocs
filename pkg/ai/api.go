package ai

import (
	"context"
	"encoding/json"

	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/constants"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
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

	fileStr := ""
	for _, file := range files {
		fileStr += file
	}

	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"endpoints": {
				Type:        jsonschema.Array,
				Description: "A string array of endpoints found in the files",
				Items: &jsonschema.Definition{
					Type:        jsonschema.String,
					Description: "The endpoint. For example GET /api/v1/users, POST /api/v1/users",
				},
			},
		},
	}
	f := openai.FunctionDefinition{
		Name:        "get_endpoints",
		Description: "Get all the endpoints in the files",
		Parameters:  params,
	}
	t := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: f,
	}

	messages := []openai.ChatCompletionMessage{
		{
			Content: constants.API_ENDPOINTS_PROMPT,
			Role:    openai.ChatMessageRoleSystem,
		},
		{
			Content: fileStr,
			Role:    openai.ChatMessageRoleUser,
		},
	}

	c := makeClient(conf)
	ctx := context.Background()

	r, err := c.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    conf.Model,
		Messages: messages,
		Tools:    []openai.Tool{t},
	})
	if err != nil {
		return nil, err
	}

	resp := r.Choices[0].Message.Content

	// parse the response
	var endpoints endpointsResp
	err = json.Unmarshal([]byte(resp), &endpoints)
	if err != nil {
		return nil, err
	}

	return endpoints.Endpoints, nil
}

func APIDocumentEndpoint(endpoint string, files []string, conf *config.Config) (string, error) {
	// join all the files into a single string
	fileStr := ""
	for _, file := range files {
		fileStr += file
	}

	// ask the AI to document the endpoint
	resp, err := request(constants.API_DOCUMENT_ENDPOINT_PROMPT, fileStr, conf)
	if err != nil {
		return "", err
	}

	return resp, nil
}
