package ai

import (
	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/constants"
	"github.com/sashabaranov/go-openai"
)

func CmdQuestion(history []string, chatPrompt string, conf *config.Config) (*openai.ChatCompletionStream, error) {
	questionNoHistory := "\nQuestion: " + chatPrompt + "\n\nAnswer:"
	historyQuestion := "Shell History:\n"

	qTokens := calc.EstimateTokens(questionNoHistory)
	commandPromptTokens := calc.EstimateTokens(constants.COMMAND_QUESTION_PROMPT)

	// loop backwards through history to find the most recent question
	for i := len(history) - 1; i >= 0; i-- {
		newHistory := history[i] + "\n"
		tokens := calc.EstimateTokens(newHistory) + qTokens + calc.EstimateTokens(historyQuestion) + commandPromptTokens
		if tokens < calc.GetMaxTokens(conf.Model) {
			historyQuestion += newHistory
		} else {
			break
		}
	}

	question := historyQuestion + questionNoHistory

	return requestStream(constants.COMMAND_QUESTION_PROMPT, question, conf)
}
