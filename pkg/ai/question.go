package ai

import (
	"strings"

	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/constants"
)

func Question(filePath, fileContent, chatPrompt string, conf *config.Config) (string, error) {
	question := "File Name: " + filePath + "\nQuestion: " + chatPrompt + "\n\n" + strings.TrimRight(string(fileContent), " \n") + "\nAnswer:"

	return request(constants.QUESTION_PROMPT, question, conf)
}
