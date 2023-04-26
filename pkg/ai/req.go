package ai

import (
	"context"
	"errors"

	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/sashabaranov/go-openai"
)

func request(systemMsg, userMsg string, conf *config.Config) (string, error) {

	c := openai.NewClient(conf.APIKey)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model: conf.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Content: systemMsg,
				Role:    openai.ChatMessageRoleSystem,
			},
			{
				Content: userMsg,
				Role:    openai.ChatMessageRoleUser,
			},
		},
	}

	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

func SimpleRequest(prompt string, conf *config.Config) (string, error) {
	c := openai.NewClient(conf.APIKey)
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model: conf.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Content: prompt,
				Role:    openai.ChatMessageRoleUser,
			},
		},
	}

	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}
