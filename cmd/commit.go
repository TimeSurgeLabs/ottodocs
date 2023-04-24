/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/git"
	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generates a commit message from the git diff",
	Long:  `Uses the git diff to generate a commit message. Requires Git to be installed on the system.`,
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

		dirty, err := git.IsDirty()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if !dirty {
			log.Error("No changes to commit.")
			os.Exit(1)
		}

		log.Info("Generating commit message...")
		log.Debug("Getting git diff...")
		diff, err := git.Diff()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Debug("Sending diff to ChatGPT...")
		msg, err := ai.CommitMessage(diff, conventional, conf)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if auto || push {
			log.Info("Adding and committing...")
			output, err := git.AddAll()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			fmt.Println(output)
			output, err = git.Commit(msg)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			fmt.Println(output)
			if push {
				log.Info("Pushing...")
				output, err = git.Push()
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
				fmt.Println(output)
			}
			os.Exit(0)
		}

		if plain {
			fmt.Println(msg)
		} else {
			fmt.Println("Commit message:", msg)
		}
	},
}

func init() {
	RootCmd.AddCommand(commitCmd)

	commitCmd.Flags().BoolVarP(&conventional, "conventional", "c", false, "use conventional commits")
	commitCmd.Flags().BoolVarP(&plain, "plain", "p", false, "no output formatting")
	commitCmd.Flags().BoolVarP(&auto, "auto", "a", false, "automatically add all and commit with the generated message")
	commitCmd.Flags().BoolVar(&push, "push", false, "automatically push to the current branch")
	commitCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
