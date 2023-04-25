package ai

import (
	"github.com/chand1012/git2gpt/prompt"

	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/constants"
)

func Question(files []prompt.GitFile, chatPrompt string, conf *config.Config) (string, error) {
	question := "\nGiven the context of the above code, answer the following question.\nQuestion: " + chatPrompt + "\nAnswer:"
	t, err := calc.PreciseTokens(question)
	if err != nil {
		return "", err
	}

	tokens := int64(t)

	var prompt string
	for _, file := range files {
		if file.Tokens+tokens > int64(calc.GetMaxTokens(conf.Model)) {
			break
		}
		prompt += "Filename: " + file.Path + "\n" + file.Contents + "\n"
		tokens += file.Tokens
	}

	return request(constants.QUESTION_PROMPT, prompt+question, conf)
}
