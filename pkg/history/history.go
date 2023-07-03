package history

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
)

type HistoryFile struct {
	DisplayName string                         `json:"display_name"`
	FileName    string                         `json:"file_name"`
	Messages    []openai.ChatCompletionMessage `json:"messages"`
}

func historyDirExists() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	historyPath := filepath.Join(usr.HomeDir, ".ottodocs", "history")
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		err = os.MkdirAll(historyPath, 0755)
		if err != nil {
			return "", err
		}
	}
	return historyPath, nil
}

func ListHistoryFiles() ([]string, error) {
	historyPath, err := historyDirExists()
	if err != nil {
		return nil, err
	}
	dir, err := os.Open(historyPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	// sort by the numerical value of their file name. Higher number = more recent. More recent = lower index
	sort.Slice(files, func(i, j int) bool {
		// get the numerical value of the file name. Remove .json and convert to int64
		iNum, err := strconv.ParseInt(strings.ReplaceAll(files[i], ".json", ""), 10, 64)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		jNum, err := strconv.ParseInt(strings.ReplaceAll(files[j], ".json", ""), 10, 64)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		return iNum > jNum
	})

	return files, nil
}

func LoadHistoryFile(fileName string) (*HistoryFile, error) {
	historyPath, err := historyDirExists()
	if err != nil {
		return nil, err
	}

	var filePath string
	if strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		filePath = filepath.Clean(fileName)
	} else {
		filePath = filepath.Join(historyPath, fileName)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var historyFile HistoryFile
	err = json.NewDecoder(file).Decode(&historyFile)
	if err != nil {
		return nil, err
	}
	return &historyFile, nil
}

func DeleteHistoryFile(fileName string) error {
	historyPath, err := historyDirExists()
	if err != nil {
		return err
	}

	var filePath string
	if strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		filePath = filepath.Clean(fileName)
	} else {
		filePath = filepath.Join(historyPath, fileName)
	}

	return os.Remove(filePath)
}

func ClearHistory() error {
	files, err := ListHistoryFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		err = DeleteHistoryFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveHistoryFile(historyFile *HistoryFile) error {
	historyPath, err := historyDirExists()
	if err != nil {
		return err
	}
	filePath := filepath.Join(historyPath, historyFile.FileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(historyFile)
	if err != nil {
		return err
	}
	return nil
}

func UpdateHistoryFile(messages []openai.ChatCompletionMessage, fileName string) error {
	historyFile, err := LoadHistoryFile(fileName)
	if err != nil {
		return err
	}

	historyFile.Messages = messages

	return SaveHistoryFile(historyFile)
}

func NewHistoryFile(messages []openai.ChatCompletionMessage, fileName, displayName string) *HistoryFile {
	return &HistoryFile{
		DisplayName: displayName,
		FileName:    fileName,
		Messages:    messages,
	}
}

func genDisplayName(messages []openai.ChatCompletionMessage, conf *config.Config) (string, error) {
	c := openai.NewClient(conf.APIKey)
	ctx := context.Background()

	messages = append(messages, openai.ChatCompletionMessage{
		Content: "Write me a display name for the conversation. It should be no longer than 5 words. It should be relevant to the previous conversation.",
		Role:    openai.ChatMessageRoleUser,
	})

	req := openai.ChatCompletionRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

// designed to run on the first run of the chat command
// runs in the background
func SaveInitialHistory(messages []openai.ChatCompletionMessage, conf *config.Config) (string, error) {
	displayName, err := genDisplayName(messages, conf)
	if err != nil {
		return "", err
	}

	// get the current time in seconds
	// this will be used as the file name
	now := time.Now().Unix()

	fileName := fmt.Sprintf("%d.json", now)

	historyFile := NewHistoryFile(messages, fileName, displayName)

	err = SaveHistoryFile(historyFile)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
