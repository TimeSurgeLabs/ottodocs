package ai

import (
	"fmt"
	"strings"

	gopenai "github.com/CasualCodersProjects/gopenai"
	ai_types "github.com/CasualCodersProjects/gopenai/types"
	"github.com/chand1012/ottodocs/constants"
)

func Question(filePath, fileContent, chatPrompt, APIKey string) (string, error) {
	openai := gopenai.NewOpenAI(&gopenai.OpenAIOpts{
		APIKey: APIKey,
	})

	question := "File Name: " + filePath + "\nQuestion: " + chatPrompt + "\n\n" + strings.TrimRight(string(fileContent), " \n") + "\nAnswer:"

	messages := []ai_types.ChatMessage{
		{
			Content: constants.QUESTION_PROMPT,
			Role:    "system",
		},
		{
			Content: question,
			Role:    "user",
		},
	}

	tokens, err := CalcTokens(messages[0].Content, messages[1].Content)
	if err != nil {
		return "", fmt.Errorf("could not calculate tokens: %s", err)
	}

	maxTokens := constants.OPENAI_MAX_TOKENS - tokens

	if maxTokens < 0 {
		return "", fmt.Errorf("the prompt is too long. max length is %d. Got %d", constants.OPENAI_MAX_TOKENS, tokens)
	}

	req := ai_types.NewDefaultChatRequest("")
	req.Messages = messages
	req.MaxTokens = maxTokens
	// lower the temperature to make the model more deterministic
	// req.Temperature = 0.3

	resp, err := openai.CreateChat(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "", err
	}

	message := resp.Choices[0].Message.Content

	return message, nil
}
