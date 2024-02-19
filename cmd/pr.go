/*
Copyright © 2024 TimeSurgeLabs <chandler@timesurgelabs.com>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/TimeSurgeLabs/ottodocs/pkg/ai"
	"github.com/TimeSurgeLabs/ottodocs/pkg/calc"
	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/gh"
	"github.com/TimeSurgeLabs/ottodocs/pkg/git"
	"github.com/TimeSurgeLabs/ottodocs/pkg/utils"
	g "github.com/chand1012/git2gpt/prompt"
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

		utils.PrintColoredText("Title: ", c.OttoColor)
		if title == "" {
			// generate the title
			log.Debug("Generating title...")
			stream, err := ai.PRTitle(logs, c)
			if err != nil {
				log.Errorf("Error generating title: %s", err)
				os.Exit(1)
			}
			title, err = utils.PrintChatCompletionStream(stream)
			if err != nil {
				log.Errorf("Error printing chat completion stream: %s", err)
				os.Exit(1)
			}
		} else {
			fmt.Println(title)
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
				// compress the diff with Otto
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

		if issuePRNumber != 0 {
			remote, err := git.GetRemote("origin")
			if err != nil {
				log.Errorf("Error getting remote: %s", err)
				os.Exit(1)
			}

			log.Debug("Getting repo info...")
			// get repo and owner
			owner, repo, err := git.ExtractOriginInfo(remote)
			if err != nil {
				log.Errorf("Error extracting origin info: %s", err)
				os.Exit(1)
			}

			log.Debugf("Owner: %s, Repo: %s", owner, repo)
			log.Debug("Getting issue...")

			// get issue
			issue, err := gh.GetIssue(owner, repo, issuePRNumber, c)
			if err != nil {
				log.Errorf("Error getting issue: %s", err)
				os.Exit(1)
			}

			body := issue.Issue.Body
			title := issue.Issue.Title

			log.Debug("Constructing prompt...")
			prompt += "\n\nRelated Issue Title: " + title + "\nRelated Issue Body: " + body
		}

		utils.PrintColoredText("Body: ", c.OttoColor)
		stream, err := ai.PRBody(prompt, c)
		if err != nil {
			log.Errorf("Error generating PR body: %s", err)
			os.Exit(1)
		}

		body, err := utils.PrintChatCompletionStream(stream)
		if err != nil {
			log.Errorf("Error printing chat completion stream: %s", err)
			os.Exit(1)
		}

		utils.PrintColoredText("Branch: ", c.OttoColor)
		fmt.Println(base)

		if !push {
			os.Exit(0)
		}

		if !force {
			confirm, err := utils.Input("Publish PR? (y/n): ")
			if err != nil {
				log.Errorf("Error getting input: %s", err)
				os.Exit(1)
			}
			confirm = strings.ToLower(confirm)
			if confirm != "y" {
				fmt.Println("Exiting...")
				os.Exit(0)
			}
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

		fmt.Println("Opening pull request...")
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
	prCmd.Flags().BoolVarP(&force, "force", "f", false, "Force the creation of the pull request")
	prCmd.Flags().IntVarP(&issuePRNumber, "issue", "i", 0, "Issue number to associate with the pull request")
}
