package ai

import (
	"strings"

	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/constants"
)

func Markdown(filePath, contents, chatPrompt string, conf *config.Config) (string, error) {

	question := chatPrompt + "\n\n" + strings.TrimRight(contents, " \n")

	return request(constants.DOCUMENT_MARKDOWN_PROMPT, question, conf)
}
