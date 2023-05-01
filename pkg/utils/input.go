package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input(prompt string) (string, error) {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}

func InputWithColor(prompt, cssColorCode string) (string, error) {
	PrintColoredText(prompt, cssColorCode)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}
