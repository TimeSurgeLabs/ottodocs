package ai

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	gopenai "github.com/CasualCodersProjects/gopenai"
	ai_types "github.com/CasualCodersProjects/gopenai/types"
	"github.com/chand1012/ottodocs/constants"
	"github.com/pandodao/tokenizer-go"
)

func Markdown(filePath, chatPrompt, APIKey string) (string, error) {
	openai := gopenai.NewOpenAI(&gopenai.OpenAIOpts{
		APIKey: APIKey,
	})

	// load the file
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	if !utf8.Valid(contents) {
		return "", errors.New("the file is not valid UTF-8")
	}

	question := chatPrompt + "\n\n" + strings.TrimRight(string(contents), " \n")

	messages := []ai_types.ChatMessage{
		{
			Content: constants.DOCUMENT_MARKDOWN_PROMPT,
			Role:    "system",
		},
		{
			Content: question,
			Role:    "user",
		},
	}

	tokens := tokenizer.MustCalToken(messages[0].Content) + tokenizer.MustCalToken(messages[1].Content)

	maxTokens := constants.OPENAI_MAX_TOKENS - tokens

	if maxTokens < 0 {
		return "", fmt.Errorf("the prompt is too long. max length is %d. Got %d", constants.OPENAI_MAX_TOKENS, tokens)
	}

	req := ai_types.NewDefaultChatRequest("")
	req.Messages = messages
	req.MaxTokens = maxTokens
	// lower the temperature to make the model more deterministic
	// req.Temperature = 0.3

	// ask ChatGPT the question
	resp, err := openai.CreateChat(req)
	if err != nil {
		return "", err
	}

	message := resp.Choices[0].Message.Content

	return message, nil
}
