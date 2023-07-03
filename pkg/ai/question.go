package ai

import (
	"fmt"

	"github.com/chand1012/git2gpt/prompt"
	"github.com/sashabaranov/go-openai"

	"github.com/TimeSurgeLabs/ottodocs/pkg/calc"
	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/constants"
)

func Question(files []prompt.GitFile, chatPrompt string, conf *config.Config, verbose bool) (*openai.ChatCompletionStream, error) {
	question := "\nGiven the context of the above code, answer the following question.\nQuestion: " + chatPrompt + "\nAnswer:"
	t, err := calc.PreciseTokens(question)
	if err != nil {
		return nil, err
	}

	tokens := int64(t)

	var prompt string
	for _, file := range files {
		if file.Tokens+tokens > int64(calc.GetMaxTokens(conf.Model)) {
			break
		}
		if verbose {
			fmt.Println("Adding file: " + file.Path)
		}
		prompt += "Filename: " + file.Path + "\n" + file.Contents + "\n"
		tokens += file.Tokens
	}

	return requestStream(constants.QUESTION_PROMPT, prompt+question, conf)
}
