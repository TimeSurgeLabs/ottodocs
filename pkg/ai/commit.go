package ai

import (
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/constants"
)

func CommitMessage(diff string, conventional bool, conf *config.Config) (string, error) {
	sysMessage := constants.GIT_DIFF_PROMPT_STD
	if conventional {
		sysMessage = constants.GIT_DIFF_PROMPT_CONVENTIONAL
	}

	return request(sysMessage, diff, conf)
}
