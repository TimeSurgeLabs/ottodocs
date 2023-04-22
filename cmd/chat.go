/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var question string

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Ask ChatGPT a question from the command line.",
	Long: `Ask ChatGPT a question from the command line.

If '-q' is not specified, the user will be prompted to enter a question.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Load()
		if err != nil || conf.APIKey == "" {
			// if the API key is not set, prompt the user to login
			log.Error("Please login first.")
			log.Error("Run `ottodocs login` to login.")
			os.Exit(1)
		}

		client := openai.NewClient(conf.APIKey)

		// if the question is not provided, prompt the user for it
		if question == "" {
			fmt.Print("What would you like to chat ChatGPT?\n> ")
			fmt.Scanln(&question)
		}

		resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: conf.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Content: question,
					Role:    openai.ChatMessageRoleUser,
				},
			},
		})

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if len(resp.Choices) == 0 {
			log.Error("No choices returned")
			os.Exit(1)
		}

		fmt.Println(resp.Choices[0].Message.Content)
	},
}

func init() {
	RootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&question, "question", "q", "", "Question to chat ChatGPT")
}
