package ai

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/constants"
	"github.com/TimeSurgeLabs/ottodocs/pkg/textfile"
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

// Document a file using the OpenAI Otto API. Takes a file path, a prompt, and an API key as arguments.
func SingleFile(filePath, contents, chatPrompt string, conf *config.Config) (string, error) {

	fileEnding := filepath.Ext(filePath)

	commentOperator, ok := constants.CommentOperators[fileEnding]
	if !ok {
		return "", fmt.Errorf("the file type %s is not supported", fileEnding)
	}

	question := chatPrompt + "\n\n" + strings.TrimRight(contents, " \n")

	message, err := request(constants.DOCUMENT_FILE_PROMPT, question, conf)
	if err != nil {
		return "", err
	}

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
