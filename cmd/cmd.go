/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/shell"
)

// cmdCmd represents the cmd command
var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Browses Bash",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		history, err := shell.GetHistory(100)
		if err != nil {
			log.Error("This command is only supported on MacOS and Linux using Bash or Zsh. Windows and other shells coming soon!")
			log.Error(err)
			os.Exit(1)
		}

		// fmt.Println("History:", history)

		resp, err := ai.CmdQuestion(history, chatPrompt, conf.APIKey, conf.Model)

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		fmt.Println("Answer:", resp)
	},
}

func init() {
	RootCmd.AddCommand(cmdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cmdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cmdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
