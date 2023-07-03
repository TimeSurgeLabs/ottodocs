/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/TimeSurgeLabs/ottodocs/pkg/calc"
	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/git"
	"github.com/TimeSurgeLabs/ottodocs/pkg/utils"
	"github.com/spf13/cobra"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Count tokens in given context and prompt",
	Long: `This command calculates the token count of the given context and prompt.
It takes the context and prompt files as input.

Example usage:

otto count -c contextfile.txt -g "prompt"`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Load()
		if err != nil {
			log.Errorf("Error loading config: %s", err)
			os.Exit(1)
		}

		var prompt string

		if chatPrompt == "" && len(contextFiles) == 0 && !repoContext {
			log.Error("Requires a prompt or context file as an argument. Example: otto count -c contextfile.txt -g \"prompt\"")
			os.Exit(1)
		}

		if chatPrompt != "" {
			prompt = "GOAL:" + chatPrompt
		}

		if repoContext {
			repo, err := git.GetRepo(".", "", false)
			if err != nil {
				log.Errorf("Error getting repo: %s", err)
				os.Exit(1)
			}

			for _, file := range repo.Files {
				content := file.Contents
				prompt += "\n\nFILE: " + file.Path + "\n\n" + content + "\n"
			}
		} else {
			for _, contextFile := range contextFiles {
				content, err := utils.LoadFile(contextFile)
				if err != nil {
					log.Errorf("Error loading file: %s", err)
					os.Exit(1)
				}
				prompt += "\n\nFILE: " + contextFile + "\n\n" + content + "\n"
			}
		}

		tokens, err := calc.PreciseTokens(prompt)
		if err != nil {
			log.Errorf("Error calculating tokens: %s", err)
			os.Exit(1)
		}

		utils.PrintColoredText("Token count: ", c.OttoColor)
		fmt.Println(tokens)
	},
}

func init() {
	RootCmd.AddCommand(countCmd)

	countCmd.Flags().StringVarP(&chatPrompt, "goal", "g", "", "Prompt for token count")
	countCmd.Flags().StringSliceVarP(&contextFiles, "context", "c", []string{}, "Context files")
	countCmd.Flags().BoolVarP(&repoContext, "repo", "r", false, "Use the current repository as context")
}
