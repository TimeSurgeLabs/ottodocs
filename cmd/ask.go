/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	gopenai "github.com/CasualCodersProjects/gopenai"
	"github.com/chand1012/ottodocs/config"
	"github.com/spf13/cobra"
)

var question string

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask ChatGPT a question from the command line.",
	Long: `Ask ChatGPT a question from the command line.
	
If '-q' is not specified, the user will be prompted to enter a question.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Load()
		if err != nil || conf.APIKey == "" {
			// if the API key is not set, prompt the user to login
			fmt.Println("Please login first.")
			fmt.Println("Run `ottodocs login` to login.")
			os.Exit(1)
		}

		openai := gopenai.NewOpenAI(&gopenai.OpenAIOpts{
			APIKey: conf.APIKey,
		})

		// if the question is not provided, prompt the user for it
		if question == "" {
			fmt.Print("What would you like to ask ChatGPT?\n> ")
			fmt.Scanln(&question)
		}

		// ask ChatGPT the question
		resp, err := openai.CreateChatSimple(question, 0) // 0 sets the max tokens to 1024
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}

		message := resp.Choices[0].Message.Content

		fmt.Println(message)
	},
}

func init() {
	rootCmd.AddCommand(askCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// askCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// askCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	askCmd.Flags().StringVarP(&question, "question", "q", "", "Question to ask ChatGPT")
}
