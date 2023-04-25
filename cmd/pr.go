/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	g "github.com/chand1012/git2gpt/prompt"
	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/gh"
	"github.com/chand1012/ottodocs/pkg/git"
	"github.com/chand1012/ottodocs/pkg/utils"
	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Generate a pull request",
	Long: `The "pr" command generates a pull request by combining commit messages, a title, and the git diff between branches.
Requires Git to be installed on the system. If a title is not provided, one will be generated.`,
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

		if !git.IsGitRepo(".") {
			log.Error("Error: not a git repository")
			os.Exit(1)
		}

		log.Info("Generating PR...")

		currentBranch, err := git.GetBranch()
		if err != nil {
			log.Errorf("Error getting current branch: %s", err)
			os.Exit(1)
		}

		log.Debugf("Current branch: %s", currentBranch)

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

		log.Debugf("Got %d logs", len(strings.Split(logs, "\n")))

		if title == "" {
			// generate the title
			log.Debug("Generating title...")
			title, err = ai.PRTitle(logs, c)
			if err != nil {
				log.Errorf("Error generating title: %s", err)
				os.Exit(1)
			}
		}

		log.Debugf("Title: %s", title)
		// get the diff
		diff, err := git.GetBranchDiff(base, currentBranch)
		if err != nil {
			log.Errorf("Error getting diff: %s", err)
			os.Exit(1)
		}

		log.Debug("Calculating Diff Tokens...")
		// count the diff tokens
		diffTokens, err := calc.PreciseTokens(diff)
		if err != nil {
			log.Errorf("Error counting diff tokens: %s", err)
			os.Exit(1)
		}

		log.Debug("Calculating Title Tokens...")
		titleTokens, err := calc.PreciseTokens(title)
		if err != nil {
			log.Errorf("Error counting title tokens: %s", err)
			os.Exit(1)
		}

		if diffTokens == 0 {
			log.Warn("Diff is empty!")
		}

		if titleTokens == 0 {
			log.Warn("Title is empty!")
		}

		log.Debugf("Diff tokens: %d", diffTokens)
		log.Debugf("Title tokens: %d", titleTokens)
		var prompt string
		if diffTokens+titleTokens > calc.GetMaxTokens(c.Model) {
			log.Debug("Diff is large, creating compressed diff and using logs and title")
			prompt = "Title: " + title + "\n\nGit logs: " + logs
			// get a list of the changed files
			files, err := git.GetChangedFilesBranches(base, currentBranch)
			if err != nil {
				log.Errorf("Error getting changed files: %s", err)
				os.Exit(1)
			}
			ignoreFiles := g.GenerateIgnoreList(".", ".gptignore", false)
			for _, file := range files {
				if utils.Contains(ignoreFiles, file) {
					log.Debugf("Ignoring file: %s", file)
					continue
				}
				log.Debug("Compressing diff for file: " + file)
				// get the file's diff
				fileDiff, err := git.GetFileDiffBranches(base, currentBranch, file)
				if err != nil {
					log.Errorf("Error getting file diff: %s", err)
					continue
				}
				// compress the diff with ChatGPT
				compressedDiff, err := ai.CompressDiff(fileDiff, c)
				if err != nil {
					log.Errorf("Error compressing diff: %s", err)
					continue
				}
				prompt += "\n\n" + file + ":\n" + compressedDiff
			}
		} else {
			log.Debug("Diff is small enough, using logs, title, and diff")
			prompt = "Title: " + title + "\n\nGit logs: " + logs + "\n\nGit diff: " + diff
		}

		body, err := ai.PRBody(prompt, c)
		if err != nil {
			log.Errorf("Error generating PR body: %s", err)
			os.Exit(1)
		}

		if !push {
			fmt.Println("Title: ", title)
			fmt.Println("Body: ", body)
			os.Exit(0)
		}

		// get the origin remote
		origin, err := git.GetRemote(remote)
		if err != nil {
			log.Errorf("Error getting remote: %s", err)
			os.Exit(1)
		}

		owner, repo, err := git.ExtractOriginInfo(origin)
		if err != nil {
			log.Errorf("Error extracting origin info: %s", err)
			os.Exit(1)
		}

		// print the origin and repo if debug is enabled
		log.Debugf("Origin: %s", origin)
		log.Debugf("Owner: %s", owner)
		log.Debugf("Repo: %s", repo)

		data := make(map[string]string)
		data["title"] = title
		data["body"] = body + "\n\n" + c.Signature
		data["head"] = currentBranch
		data["base"] = base

		log.Info("Opening pull request...")
		prNumber, err := gh.OpenPullRequest(data, owner, repo, c)
		if err != nil {
			log.Errorf("Error opening pull request: %s", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully opened pull request: %s\n", title)
		// link to the pull request
		fmt.Printf("https://github.com/%s/%s/pull/%d\n", owner, repo, prNumber)
	},
}

func init() {
	RootCmd.AddCommand(prCmd)

	prCmd.Flags().StringVarP(&base, "base", "b", "", "Base branch to create the pull request against")
	prCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the pull request")
	prCmd.Flags().StringVarP(&remote, "remote", "r", "origin", "Remote for creating the pull request. Only works with GitHub.")
	prCmd.Flags().BoolVarP(&push, "publish", "p", false, "Create the pull request. Must have a remote named \"origin\"")
	prCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}
