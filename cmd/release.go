/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/constants"
	"github.com/chand1012/ottodocs/pkg/gh"
	"github.com/chand1012/ottodocs/pkg/git"
	"github.com/chand1012/ottodocs/pkg/utils"
	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Generate GitHub release notes from git commit logs",
	Long: `This command generates GitHub release notes from git commit logs.
It will create a new release given a tag and post it to GitHub as a draft.`,
	Aliases: []string{"r"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(l.DebugLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Load()
		if err != nil {
			log.Errorf("Error loading config: %s", err)
			os.Exit(1)
		}

		if previousTag == "" {
			previousTag, err = utils.Input("Previous tag: ")
			if err != nil {
				log.Errorf("Error getting previous tag: %s", err)
				os.Exit(1)
			}
		}

		if currentTag == "" {
			currentTag, err = utils.Input("Current tag: ")
			if err != nil {
				log.Errorf("Error getting current tag: %s", err)
				os.Exit(1)
			}
		}

		fmt.Print("Release notes: ")

		// get the log between the tags
		gitLog, err := git.LogBetween(previousTag, currentTag)
		if err != nil {
			log.Errorf("Error getting log between tags: %s", err)
			os.Exit(1)
		}

		prompt := constants.RELEASE_PROMPT + gitLog

		stream, err := ai.SimpleStreamRequest(prompt, c)
		if err != nil {
			log.Errorf("Error getting response: %s", err)
			os.Exit(1)
		}

		releaseNotes, err := utils.PrintChatCompletionStream(stream)
		if err != nil {
			log.Errorf("Error printing completion stream: %s", err)
			os.Exit(1)
		}

		if !force {
			confirm, err := utils.Input("Create release? (y/n): ")
			if err != nil {
				log.Errorf("Error getting confirmation: %s", err)
				os.Exit(1)
			}
			confirm = strings.ToLower(confirm)
			if confirm != "y" {
				os.Exit(0)
			}
		}

		origin, err := git.GetRemote("origin")
		if err != nil {
			log.Errorf("Error getting remote: %s", err)
			os.Exit(1)
		}

		owner, repo, err := git.ExtractOriginInfo(origin)
		if err != nil {
			log.Errorf("Error extracting origin info: %s", err)
			os.Exit(1)
		}

		err = gh.CreateDraftRelease(owner, repo, currentTag, releaseNotes, currentTag, c)
		if err != nil {
			log.Errorf("Error creating release: %s", err)
			os.Exit(1)
		}

		fmt.Println("Release created successfully!")
	},
}

func init() {
	RootCmd.AddCommand(releaseCmd)

	releaseCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	releaseCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite of existing file")
	releaseCmd.Flags().StringVarP(&previousTag, "prev-tag", "p", "", "Previous tag")
	releaseCmd.Flags().StringVarP(&currentTag, "tag", "t", "", "Current tag")
}
