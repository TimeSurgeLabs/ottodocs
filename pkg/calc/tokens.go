package calc

import "github.com/pandodao/tokenizer-go"

// Precise token count
func PreciseTokens(inputs ...string) (int, error) {
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

// Estimate token count. Much faster than CalcTokens, but less accurate.
func EstimateTokens(inputs ...string) int {
	total := int(0)
	for _, input := range inputs {
		tokens := len(input) / 4
		total += tokens
	}

	return total
}