/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	gopenai "github.com/CasualCodersProjects/gopenai"
	"github.com/chand1012/ottodocs/pkg/config"
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

		openai := gopenai.NewOpenAI(&gopenai.OpenAIOpts{
			APIKey: conf.APIKey,
		})

		// if the question is not provided, prompt the user for it
		if question == "" {
			fmt.Print("What would you like to chat ChatGPT?\n> ")
			fmt.Scanln(&question)
		}

		// chat ChatGPT the question
		resp, err := openai.CreateChatSimple(question, 0) // 0 sets the max tokens to 1024
		if err != nil {
			log.Errorf("Error: %s", err)
			os.Exit(1)
		}

		message := resp.Choices[0].Message.Content

		fmt.Println(message)
	},
}

func init() {
	RootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&question, "question", "q", "", "Question to chat ChatGPT")
}
