package ai

import (
	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/constants"
	"github.com/sashabaranov/go-openai"
)

func CommitMessage(diff string, conventional bool, conf *config.Config) (*openai.ChatCompletionStream, error) {
	sysMessage := constants.GIT_DIFF_PROMPT_STD
	if conventional {
		sysMessage = constants.GIT_DIFF_PROMPT_CONVENTIONAL
	}

	return requestStream(sysMessage, diff, conf)
}
