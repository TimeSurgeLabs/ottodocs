package ai

import "github.com/pandodao/tokenizer-go"

func CalcTokens(inputs ...string) (int, error) {
	total := int(0)
	for _, input := range inputs {
		tokens, err := tokenizer.CalToken(input)
		if err != nil {
			return -1, err
		}
		total += tokens
	}

	return total, nil
}
