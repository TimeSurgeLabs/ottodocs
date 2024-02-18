/*
Copyright © 2024 TimeSurgeLabs <chandler@timesurgelabs.com>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/TimeSurgeLabs/ottodocs/pkg/calc"
	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/history"
	"github.com/TimeSurgeLabs/ottodocs/pkg/utils"
	l "github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

// this should be improved with a proper vector db
// so that proper context can be added

// for now it just removes the oldest message
// when it reaches the max context size

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Talk with Otto from the command line!",
	Long: `Talk with Otto from the command line!
No code context is passed in this mode. Emulates web chat.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(l.DebugLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Load()
		if err != nil || conf.APIKey == "" {
			// if the API key is not set, prompt the user to config
			log.Error("Please config first.")
			log.Error("Run `ottodocs config -h` to learn how to config.")
			os.Exit(1)
		}

		// error if read only mode set without a history file
		if readOnly && loadHistory == "" {
			log.Error("Read only mode requires a history file. Use --history to specify a history file.")
			os.Exit(1)
		}

		// all messages in the conversation
		var messages []openai.ChatCompletionMessage
		var fileName string // history file name

		client := openai.NewClient(conf.APIKey)

		if displayHistory {
			files, err := history.ListHistoryFiles()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			if len(files) == 0 {
				fmt.Println("No chat history found.")
				os.Exit(0)
			}

			for i, file := range files {
				// load each file
				historyFile, err := history.LoadHistoryFile(file)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
				fmt.Printf("%d. %s\n", i, historyFile.DisplayName)
			}
			os.Exit(0)
		}

		if deleteHistory != "" {
			if strings.HasSuffix(deleteHistory, ".json") {
				// load the file
				_, err = history.LoadHistoryFile(deleteHistory)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else {
				i, err := strconv.Atoi(deleteHistory)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
				files, err := history.ListHistoryFiles()
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}

				if i >= len(files) {
					log.Error("Index out of bounds")
					os.Exit(1)
				}

				// load the file
				_, err = history.LoadHistoryFile(files[i])
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}

			confirm, err := utils.Input("Are you sure you want to delete this history file? (y/N): ")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			confirm = strings.ToLower(confirm)
			if confirm != "y" {
				os.Exit(0)
			}

			err = history.DeleteHistoryFile(deleteHistory)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			fmt.Println("Deleted history file.")
			os.Exit(0)
		}

		if clearHistory {
			// make sure the user wants to clear the history
			log.Warn("This will clear all chat history. This operation cannot be undone.")
			confirm, err := utils.Input("Are you sure you want to clear the chat history? (y/N): ")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			confirm = strings.ToLower(confirm)

			if confirm != "y" {
				os.Exit(0)
			}

			err = history.ClearHistory()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			fmt.Println("Cleared chat history.")
			os.Exit(0)
		}

		utils.PrintColoredText("Otto: ", conf.OttoColor)
		fmt.Println("Hello! I am Otto. Use Ctrl+C to exit at any time.")

		if loadHistory != "" {
			var historyFile *history.HistoryFile
			if strings.HasSuffix(loadHistory, ".json") {
				// load the file
				historyFile, err = history.LoadHistoryFile(loadHistory)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else {
				i, err := strconv.Atoi(loadHistory)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
				files, err := history.ListHistoryFiles()
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}

				if i >= len(files) {
					log.Error("Index out of bounds")
					os.Exit(1)
				}

				// load the file
				historyFile, err = history.LoadHistoryFile(files[i])
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}

			messages = historyFile.Messages
			fileName = loadHistory

			for _, message := range messages {
				if message.Role == openai.ChatMessageRoleUser {
					utils.PrintColoredText("You: ", conf.UserColor)
				} else {
					utils.PrintColoredText("Otto: ", conf.OttoColor)
				}
				fmt.Println(message.Content)
			}

			if readOnly {
				os.Exit(0)
			}
		}

		for {
			question, err := utils.InputWithColor("You: ", conf.UserColor)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			utils.PrintColoredText("Otto: ", conf.OttoColor)

			log.Debugf("Adding question '%s' to context...", question)

			messages = append(messages, openai.ChatCompletionMessage{
				Content: question,
				Role:    openai.ChatMessageRoleUser,
			})

			// get the length of all the messages
			messageStrings := make([]string, len(messages))
			for i, message := range messages {
				messageStrings[i] = message.Content
			}

			tokens, err := calc.PreciseTokens(messageStrings...)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			if tokens > calc.GetMaxTokens(conf.Model) {
				for tokens > calc.GetMaxTokens(conf.Model) {
					messages = messages[1:]
					messageStrings = messageStrings[1:]
					tokens, err = calc.PreciseTokens(messageStrings...)
					if err != nil {
						log.Error(err)
						os.Exit(1)
					}
				}
			}

			stream, err := client.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
				Model:    conf.Model,
				Messages: messages,
			})

			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			completeStream, err := utils.PrintChatCompletionStream(stream)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			messages = append(messages, openai.ChatCompletionMessage{
				Content: completeStream,
				Role:    openai.ChatMessageRoleAssistant,
			})

			if fileName == "" {
				fileName, err = history.SaveInitialHistory(messages, conf)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			} else {
				err = history.UpdateHistoryFile(messages, fileName)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(chatCmd)

	chatCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	chatCmd.Flags().BoolVarP(&displayHistory, "history", "H", false, "Display chat history")
	chatCmd.Flags().BoolVarP(&readOnly, "read", "r", false, "Read the history file and exit")
	chatCmd.Flags().BoolVar(&clearHistory, "clear", false, "Clear chat history")
	chatCmd.Flags().StringVarP(&loadHistory, "load", "l", "", "Load chat history from file. Can either be a file path or an index of the chat history")
	chatCmd.Flags().StringVarP(&deleteHistory, "delete", "d", "", "Delete chat history from file. Can either be a file path or an index of the chat history")
}
