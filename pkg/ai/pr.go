package ai

import (
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/constants"
	"github.com/sashabaranov/go-openai"
)

func PRTitle(gitLog string, conf *config.Config) (*openai.ChatCompletionStream, error) {
	return requestStream(constants.PR_TITLE_PROMPT, gitLog, conf)
}

func PRBody(info string, conf *config.Config) (*openai.ChatCompletionStream, error) {
	return requestStream(constants.PR_BODY_PROMPT, info, conf)
}

func CompressDiff(diff string, conf *config.Config) (string, error) {
	return request(constants.COMPRESS_DIFF_PROMPT, diff, conf)
}
