package ai

import (
	"errors"

	gopenai "github.com/CasualCodersProjects/gopenai"
	ai_types "github.com/CasualCodersProjects/gopenai/types"
	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/config"
)

func request(systemMsg, userMsg string, conf *config.Config) (string, error) {

	openai := gopenai.NewOpenAI(&gopenai.OpenAIOpts{
		APIKey: conf.APIKey,
	})

	messages := []ai_types.ChatMessage{
		{
			Content: systemMsg,
			Role:    "system",
		},
		{
			Content: userMsg,
			Role:    "user",
		},
	}

	tokens, err := calc.PreciseTokens(messages[0].Content, messages[1].Content)
	if err != nil {
		return "", err
	}

	req := ai_types.NewDefaultChatRequest("")
	req.Messages = messages
	req.MaxTokens = calc.GetMaxTokens(conf.Model) - tokens
	req.Model = conf.Model

	if req.MaxTokens < 0 {
		return "", errors.New("the prompt is too long")
	}

	resp, err := openai.CreateChat(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned. Check your OpenAI API key")
	}

	message := resp.Choices[0].Message.Content

	return message, nil
}
