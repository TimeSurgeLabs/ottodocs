/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/git"
	"github.com/spf13/cobra"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Generate a pull request",
	Long: `The "pr" command generates a pull request by combining commit messages, a title, and the git diff between branches.
Requires Git to be installed on the system. If a title is not provided, one will be generated.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Load()
		if err != nil {
			log.Errorf("Error loading config: %s", err)
			os.Exit(1)
		}

		currentBranch, err := git.GetBranch()
		if err != nil {
			log.Errorf("Error getting current branch: %s", err)
			os.Exit(1)
		}

		if base == "" {
			// ask them for the base branch
			fmt.Print("Please provide a base branch: ")
			fmt.Scanln(&base)
		}

		logs, err := git.LogBetween(base, currentBranch)
		if err != nil {
			log.Errorf("Error getting logs: %s", err)
			os.Exit(1)
		}

		if title == "" {
			// generate the title
			title, err = ai.PRTitle(logs, c)
			if err != nil {
				log.Errorf("Error generating title: %s", err)
				os.Exit(1)
			}
		}

		// get the diff
		diff, err := git.Diff()
		if err != nil {
			log.Errorf("Error getting diff: %s", err)
			os.Exit(1)
		}

		// count the diff tokens
		diffTokens, err := calc.PreciseTokens(diff)
		if err != nil {
			log.Errorf("Error counting diff tokens: %s", err)
			os.Exit(1)
		}

		titleTokens, err := calc.PreciseTokens(title)
		if err != nil {
			log.Errorf("Error counting title tokens: %s", err)
			os.Exit(1)
		}

		var prompt string
		if diffTokens+titleTokens > calc.GetMaxTokens(c.Model) {
			// if the diff is too large, just use the logs and the title
			prompt = "Title: " + title + "\n\nGit logs: " + logs
		} else {
			prompt = "Title: " + title + "\n\nGit logs: " + logs + "\n\nGit diff: " + diff
		}

		body, err := ai.PRBody(prompt, c)
		if err != nil {
			log.Errorf("Error generating PR body: %s", err)
			os.Exit(1)
		}

		fmt.Println("Title: ", title)
		fmt.Println("Body: ", body)
	},
}

func init() {
	RootCmd.AddCommand(prCmd)

	prCmd.Flags().StringVarP(&base, "base", "b", "", "Base branch to create the pull request against")
	prCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the pull request")
}
