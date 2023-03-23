package utils

import "regexp"

func removeEmptyTokens(tokens []string) []string {
	// remove empty tokens
	var filteredTokens []string
	for _, token := range tokens {
		if token != "" {
			filteredTokens = append(filteredTokens, token)
		}
	}
	return filteredTokens
}

func SimpleTokenize(text string) []string {
	// regex to split on whitespace and punctuation
	re := regexp.MustCompile(`[\p{P}\p{Zs}]+`)
	// split the text into tokens
	tokens := re.Split(text, -1)
	// remove empty tokens
	tokens = removeEmptyTokens(tokens)
	return tokens
}
