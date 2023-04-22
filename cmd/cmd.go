/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/shell"
)

// cmdCmd represents the cmd command
var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Have ChatGPT suggest a command to run next",
	Long: `Have ChatGPT suggest a command to run next. This command will use your shell history to suggest a command to run next.
This command is only supported on MacOS and Linux using Bash or Zsh. Windows and other shells coming soon!`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(l.DebugLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Load()
		if err != nil || conf.APIKey == "" {
			// if the API key is not set, prompt the user to login
			log.Error("Please login first.")
			log.Error("Run `ottodocs login` to login.")
			os.Exit(1)
		}

		if chatPrompt == "" {
			chatPrompt = "What command do you recommend I use next?"
		}

		log.Info("Thinking....")
		log.Debug("Getting shell history...")
		history, err := shell.GetHistory(100)
		if err != nil {
			log.Warn("This command is only supported on MacOS and Linux using Bash or Zsh. Windows and other shells coming soon!")
			log.Error(err)
			os.Exit(1)
		}

		log.Debug("Asking ChatGPT for a command...")
		resp, err := ai.CmdQuestion(history, chatPrompt, conf)

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		fmt.Println("Answer:", resp)
	},
}

func init() {
	RootCmd.AddCommand(cmdCmd)

	cmdCmd.Flags().StringVarP(&chatPrompt, "question", "q", "", "The prompt to use for the chat session")
	cmdCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
