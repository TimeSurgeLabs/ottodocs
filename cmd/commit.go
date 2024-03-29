/*
Copyright © 2024 TimeSurgeLabs <chandler@timesurgelabs.com>
*/
package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/TimeSurgeLabs/ottodocs/pkg/ai"
	"github.com/TimeSurgeLabs/ottodocs/pkg/calc"
	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/constants"
	"github.com/TimeSurgeLabs/ottodocs/pkg/git"
	"github.com/TimeSurgeLabs/ottodocs/pkg/utils"
	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type fileDiff struct {
	Diff   string
	File   string
	Tokens int
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:     "commit",
	Short:   "Generates a commit message from the git diff",
	Long:    `Uses the git diff to generate a commit message. Requires Git to be installed on the system.`,
	Aliases: []string{"cm"},
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

		if !git.IsGitRepo(".") {
			log.Error("Error: not a git repository")
			os.Exit(1)
		}

		dirty, err := git.IsDirty()
		if err != nil {
			log.Errorf("Error checking if git repo is dirty: %s", err)
			os.Exit(1)
		}

		if !dirty {
			log.Error("No changes to commit.")
			os.Exit(1)
		}

		log.Debug("Generating commit message...")
		log.Debug("Getting git diff...")
		diff, err := git.Diff()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		if diff == "" {
			log.Error("No changes to commit.")
			os.Exit(1)
		}

		log.Debug("Calculating diff tokens...")
		diffTokens, err := calc.PreciseTokens(diff)
		if err != nil {
			log.Errorf("Error calculating diff tokens: %s", err)
			os.Exit(1)
		}

		var msg string

		if diffTokens > calc.GetMaxTokens(conf.Model) {
			log.Debugf("Diff tokens %d is greater than the model maximum of tokens %d", diffTokens, calc.GetMaxTokens(conf.Model))
			log.Debug("Getting changed files...")
			files, err := git.GetChangedFiles()
			if err != nil {
				log.Errorf("Error getting changed files: %s", err)
				os.Exit(1)
			}

			var diffs []fileDiff

			for _, file := range files {
				if file == "" {
					continue
				}
				log.Debugf("Getting diff for %s...", file)
				diff, err := git.GetFileDiff(file)
				if err != nil {
					log.Warnf("Error getting diff for %s: %s", file, err)
					continue
				}

				tokens, err := calc.PreciseTokens(diff)
				if err != nil {
					log.Warnf("Error calculating tokens for %s: %s", file, err)
					continue
				}

				diffs = append(diffs, fileDiff{
					Diff:   diff,
					File:   file,
					Tokens: tokens,
				})
			}

			log.Debugf("Got %d diffs", len(diffs))
			log.Debug("Sorting diffs...")
			// sort diffs by tokens
			sort.Slice(diffs, func(i, j int) bool {
				return diffs[i].Tokens < diffs[j].Tokens
			})

			log.Debug("Combining diffs...")
			// start combining the diffs until we are under the token limit
			var combinedDiff string
			var tokenTotal int
			maxTokens := calc.GetMaxTokens(conf.Model) - 500
			log.Debugf("Max tokens: %d", maxTokens)
			for _, diff := range diffs {
				if tokenTotal+diff.Tokens > maxTokens {
					break
				}
				combinedDiff += diff.Diff + "\n"
				tokenTotal += diff.Tokens
			}
			diff = combinedDiff
			log.Debugf("Combined diff tokens: %d", tokenTotal)
		} else {
			log.Debugf("Diff tokens %d is less than the model maximum of tokens %d", diffTokens, calc.GetMaxTokens(conf.Model))
		}

		log.Debug("Sending diff to Otto...")
		utils.PrintColoredText("Commit Msg: ", conf.OttoColor)
		stream, err := ai.CommitMessage(diff, conventional, conf)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		msg, err = utils.PrintChatCompletionStream(stream)
		if err != nil {
			log.Errorf("Error printing chat completion stream: %s", err)
			os.Exit(1)
		}

		if len(msg) > 75 {
			// summarize the message
			stream, err := ai.SimpleStreamRequest(constants.SUMMARIZE_PROMPT+msg, conf)
			if err != nil {
				log.Errorf("Error summarizing commit message: %s", err)
				os.Exit(1)
			}
			utils.PrintColoredText("Summarized Commit Msg: ", conf.OttoColor)
			msg, err = utils.PrintChatCompletionStream(stream)
			if err != nil {
				log.Errorf("Error printing chat completion stream: %s", err)
				os.Exit(1)
			}
		}

		if !noCommit || push {
			if !force {
				confirm, err := utils.Input("Is this okay? (y/n): ")
				if err != nil {
					log.Errorf("Error getting input: %s", err)
					os.Exit(1)
				}
				// convert to lowercase
				confirm = strings.ToLower(confirm)
				if confirm != "y" {
					fmt.Println("Exiting...")
					os.Exit(0)
				}
			}
			fmt.Println("Adding and committing...")
			output, err := git.AddAll()
			if err != nil {
				log.Errorf("Error adding files: %s", err)
				os.Exit(1)
			}
			fmt.Println(output)
			output, err = git.Commit(msg)
			if err != nil {
				log.Errorf("Error committing: %s", err)
				os.Exit(1)
			}
			fmt.Println(output)
			if push {
				fmt.Println("Pushing...")
				output, err = git.Push()
				if err != nil {
					log.Errorf("Error pushing: %s", err)
					os.Exit(1)
				}
				fmt.Println(output)
			}
			os.Exit(0)
		}

	},
}

func init() {
	RootCmd.AddCommand(commitCmd)

	commitCmd.Flags().BoolVarP(&conventional, "conventional", "c", false, "use conventional commits")
	commitCmd.Flags().BoolVarP(&noCommit, "no-commit", "n", false, "don't commit, just print the message. Ignored if --push is set")
	commitCmd.Flags().BoolVar(&push, "push", false, "automatically push to the current branch")
	commitCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	commitCmd.Flags().BoolVarP(&force, "force", "f", false, "skip confirmation")
}
