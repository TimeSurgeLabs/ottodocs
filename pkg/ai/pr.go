package ai

import (
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/constants"
)

func PRTitle(gitLog string, conf *config.Config) (string, error) {
	return request(constants.PR_TITLE_PROMPT, gitLog, conf)
}

func PRBody(info string, conf *config.Config) (string, error) {
	return request(constants.PR_BODY_PROMPT, info, conf)
}
