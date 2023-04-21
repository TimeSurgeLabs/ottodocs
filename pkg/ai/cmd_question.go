package ai

import (
	"fmt"

	gopenai "github.com/CasualCodersProjects/gopenai"
	ai_types "github.com/CasualCodersProjects/gopenai/types"
	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/constants"
)

func CmdQuestion(history []string, chatPrompt, APIKey, model string) (string, error) {
	openai := gopenai.NewOpenAI(&gopenai.OpenAIOpts{
		APIKey: APIKey,
	})

	questionNoHistory := "\nQuestion: " + chatPrompt + "\n\nAnswer:"
	historyQuestion := "Shell History:\n"

	qTokens := calc.EstimateTokens(questionNoHistory)
	commandPromptTokens := calc.EstimateTokens(constants.COMMAND_QUESTION_PROMPT)

	// loop backwards through history to find the most recent question
	for i := len(history) - 1; i >= 0; i-- {
		newHistory := history[i] + "\n"
		tokens := calc.EstimateTokens(newHistory) + qTokens + calc.EstimateTokens(historyQuestion) + commandPromptTokens
		if tokens < constants.OPENAI_MAX_TOKENS {
			historyQuestion += newHistory
		} else {
			break
		}
	}

	question := historyQuestion + questionNoHistory

	// fmt.Println(question)

	messages := []ai_types.ChatMessage{
		{
			Content: constants.COMMAND_QUESTION_PROMPT,
			Role:    "system",
		},
		{
			Content: question,
			Role:    "user",
		},
	}

	tokens := calc.EstimateTokens(messages[0].Content, messages[1].Content)

	maxTokens := constants.OPENAI_MAX_TOKENS - tokens

	if maxTokens < 0 {
		return "", fmt.Errorf("the prompt is too long. max length is %d. Got %d", constants.OPENAI_MAX_TOKENS, tokens)
	}

	req := ai_types.NewDefaultChatRequest("")
	req.Messages = messages
	req.MaxTokens = maxTokens
	req.Model = model
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
