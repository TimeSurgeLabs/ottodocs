package calc

import (
	"fmt"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

// Precise token count
func PreciseTokens(inputs ...string) (int, error) {
	tke, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return -1, err
	}
	total := int(0)
	for _, input := range inputs {
		encodedTokens := tke.Encode(input, nil, nil)
		tokens := len(encodedTokens)
		total += tokens
	}

	return total, nil
}

// Estimate token count. Much faster than CalcTokens, but less accurate.
func EstimateTokens(inputs ...string) int {
	total := int(0)
	for _, input := range inputs {
		tokens := len(input) / 4
		total += tokens
	}

	return total
}

func GetMaxTokens(model string) int {
	if strings.Contains(model, "32k") {
		return 32768
	}
	if strings.Contains(model, "4") {
		return 8192
	}

	return 4096
}

func PreciseTokensFromModel(messages []openai.ChatCompletionMessage, model string) (num_tokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("EncodingForModel: %v", err)
		fmt.Println(err)
		return
	}

	var tokens_per_message int
	var tokens_per_name int
	if model == "gpt-3.5-turbo-0301" || model == "gpt-3.5-turbo" {
		tokens_per_message = 4
		tokens_per_name = -1
	} else if strings.Contains(model, "4") {
		tokens_per_message = 3
		tokens_per_name = 1
	} else {
		tokens_per_message = 3
		tokens_per_name = 1
	}

	for _, message := range messages {
		num_tokens += tokens_per_message
		num_tokens += len(tkm.Encode(message.Content, nil, nil))
		num_tokens += len(tkm.Encode(message.Role, nil, nil))
		if message.Name != "" {
			num_tokens += tokens_per_name
		}
	}
	num_tokens += 3
	return num_tokens
}
