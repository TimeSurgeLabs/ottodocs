package utils

import (
	"errors"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
)

func PrintChatCompletionStream(stream *openai.ChatCompletionStream) (string, error) {
	var completeStream string

	defer stream.Close()

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("")
			return completeStream, nil
		} else if err != nil {
			return "", err
		}

		if len(msg.Choices) == 1 {
			fmt.Print(msg.Choices[0].Delta.Content)
			completeStream += msg.Choices[0].Delta.Content
		} else {
			return "", errors.New("received multiple choices from ChatGPT")
		}
	}
}
