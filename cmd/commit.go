/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generates a commit message from the git diff",
	Long:  `Uses the git diff to generate a commit message. Requires Git to be installed on the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.Load()
		if err != nil || conf.APIKey == "" {
			// if the API key is not set, prompt the user to login
			log.Error("Please login first.")
			log.Error("Run `ottodocs login` to login.")
			os.Exit(1)
		}

		c := exec.Command("git", "diff")
		diffBytes, err := c.Output()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		diff := string(diffBytes)

		msg, err := ai.CommitMessage(diff, conventional, conf)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		fmt.Println("Commit message:", msg)
	},
}

func init() {
	RootCmd.AddCommand(commitCmd)

	commitCmd.Flags().BoolVarP(&conventional, "conventional", "c", false, "use conventional commits")
	commitCmd.Flags().BoolVarP(&plain, "plain", "p", false, "no output formatting")
}
