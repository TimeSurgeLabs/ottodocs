/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/utils"
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

		// all messages in the conversation
		var messages []openai.ChatCompletionMessage

		client := openai.NewClient(conf.APIKey)

		utils.PrintColoredText("Otto: ", "#008080")
		fmt.Println("Hello! I am Otto. Use Ctrl+C to exit at any time.")

		for {
			question, err := utils.InputWithColor("You: ", "#007bff")
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			utils.PrintColoredText("Otto: ", "#008080")

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
		}
	},
}

func init() {
	RootCmd.AddCommand(chatCmd)

	chatCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}
