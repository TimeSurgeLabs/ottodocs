package ai

import (
	"fmt"

	gopenai "github.com/CasualCodersProjects/gopenai"
	ai_types "github.com/CasualCodersProjects/gopenai/types"
	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/constants"
)

func CommitMessage(diff string, conventional bool, conf *config.Config) (string, error) {
	openai := gopenai.NewOpenAI(&gopenai.OpenAIOpts{
		APIKey: conf.APIKey,
	})

	sysMessage := constants.GIT_DIFF_PROMPT_STD
	if conventional {
		sysMessage = constants.GIT_DIFF_PROMPT_CONVENTIONAL
	}

	messages := []ai_types.ChatMessage{
		{
			Content: sysMessage,
			Role:    "system",
		},
		{
			Content: diff,
			Role:    "user",
		},
	}

	tokens, err := calc.PreciseTokens(messages[0].Content, messages[1].Content)
	if err != nil {
		return "", fmt.Errorf("could not calculate tokens: %s", err)
	}

	maxTokens := calc.GetMaxTokens(conf.Model) - tokens

	if maxTokens < 0 {
		return "", fmt.Errorf("the prompt is too long. max length is %d. Got %d", calc.GetMaxTokens(conf.Model), tokens)
	}

	req := ai_types.NewDefaultChatRequest("")
	req.Messages = messages
	req.MaxTokens = maxTokens
	req.Model = conf.Model

	resp, err := openai.CreateChat(req)
	if err != nil {
		return "", err
	}

	message := resp.Choices[0].Message.Content

	return message, nil
}
