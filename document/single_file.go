package document

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	gopenai "github.com/CasualCodersProjects/gopenai"
	ai_types "github.com/CasualCodersProjects/gopenai/types"
	"github.com/pandodao/tokenizer-go"

	"github.com/chand1012/ottodocs/constants"
	"github.com/chand1012/ottodocs/textfile"
)

func extractLineNumber(line string) (int, error) {
	// if the line does not contain a range, return the line number
	if !strings.Contains(line, "-") {
		lineNumber, err := strconv.Atoi(line)
		if err != nil {
			return -1, fmt.Errorf("could not parse line number: %s", err)
		}

		return lineNumber, nil
	}

	// if the line contains a range, return the first line number
	lineNumber, err := strconv.Atoi(strings.Split(line, "-")[0])
	if err != nil {
		return -1, fmt.Errorf("could not parse line number: %s", err)
	}

	return lineNumber, nil
}

// Document a file using the OpenAI ChatGPT API. Takes a file path, a prompt, and an API key as arguments.
func SingleFile(filePath, chatPrompt, APIKey string) (string, error) {

	openai := gopenai.NewOpenAI(&gopenai.OpenAIOpts{
		APIKey: APIKey,
	})

	// load the file
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	if !utf8.Valid(contents) {
		return "", errors.New("the file is not valid UTF-8")
	}

	fileEnding := filepath.Ext(filePath)

	commentOperator, ok := constants.CommentOperators[fileEnding]
	if !ok {
		return "", fmt.Errorf("the file type %s is not supported", fileEnding)
	}

	question := chatPrompt + "\n\n" + strings.TrimRight(string(contents), " \n")

	messages := []ai_types.ChatMessage{
		{
			Content: constants.DOCUMENT_FILE_PROMPT,
			Role:    "system",
		},
		{
			Content: question,
			Role:    "user",
		},
	}

	tokens := tokenizer.MustCalToken(messages[0].Content) + tokenizer.MustCalToken(messages[1].Content)

	maxTokens := constants.OPENAI_MAX_TOKENS - tokens

	if maxTokens < 0 {
		return "", fmt.Errorf("the prompt is too long. max length is %d. Got %d", constants.OPENAI_MAX_TOKENS, tokens)
	}

	req := ai_types.NewDefaultChatRequest("")
	req.Messages = messages
	req.MaxTokens = maxTokens
	// lower the temperature to make the model more deterministic
	req.Temperature = 0.3

	// ask ChatGPT the question
	resp, err := openai.CreateChat(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	message := resp.Choices[0].Message.Content

	// fmt.Println(message)
	// fmt.Println("------------------------")

	lineNumbers := []int{}
	comments := []string{}

	for _, line := range strings.Split(message, "\n") {
		if line == "" {
			continue
		}
		splits := strings.Split(line, ": ")

		lineNumber, err := extractLineNumber(splits[0])
		if err != nil {
			return "", fmt.Errorf("could not parse line number: %s", err)
		}

		lineNumbers = append(lineNumbers, lineNumber)
		comments = append(comments, commentOperator+" "+splits[1])
	}

	newContents, err := textfile.InsertLinesAtIndices(string(contents), lineNumbers, comments)
	if err != nil {
		return "", fmt.Errorf("could not insert comments: %s", err)
	}

	return newContents, nil
}
