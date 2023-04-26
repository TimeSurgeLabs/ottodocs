/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	g "github.com/chand1012/git2gpt/prompt"
	"github.com/chand1012/memory"
	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/gh"
	"github.com/chand1012/ottodocs/pkg/git"
	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:     "issue",
	Short:   "Get a prompt for or ask ChatGPT about a GitHub Issue.",
	Long:    `Get a prompt for or ask ChatGPT about a GitHub Issue. Uses the current working directory's repo by default.`,
	Aliases: []string{"i"},
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

		if issuePRNumber == 0 {
			// prompt for issue number
			fmt.Print("Please provide an issue number: ")
			fmt.Scanln(&issuePRNumber)
		}

		if question == "" {
			// prompt for question
			fmt.Print("Please provide a question: ")
			fmt.Scanln(&question)
		}

		remote, err := git.GetRemote("origin")
		if err != nil {
			log.Errorf("Error getting remote: %s", err)
			os.Exit(1)
		}

		// get repo and owner
		owner, repo, err := git.ExtractOriginInfo(remote)
		if err != nil {
			log.Errorf("Error extracting origin info: %s", err)
			os.Exit(1)
		}

		// get issue
		issue, err := gh.GetIssue(owner, repo, issuePRNumber, c)
		if err != nil {
			log.Errorf("Error getting issue: %s", err)
			os.Exit(1)
		}

		body := issue.Issue.Body
		title := issue.Issue.Title

		var prompt string
		// generate a prompt
		suffix := fmt.Sprintf("\n\nGiven the following context, answer the following question: %s", question)
		prompt = fmt.Sprintf("Issue #%d: Title: %s\n\n%s", issuePRNumber, title, body)
		tokens, err := calc.PreciseTokens(prompt)
		if err != nil {
			log.Errorf("Error getting tokens: %s", err)
			os.Exit(1)
		}
		suffixTokens, err := calc.PreciseTokens(suffix)
		if err != nil {
			log.Errorf("Error getting tokens: %s", err)
			os.Exit(1)
		}
		tokens += suffixTokens

		if tokens > calc.GetMaxTokens(c.Model) {
			log.Errorf("Error: prompt is too long. Max tokens: %d, prompt tokens: %d", calc.GetMaxTokens(c.Model), tokens)
			os.Exit(1)
		}

		if useComments {
			// make the prompt just the issue body to start
			for _, comment := range issue.Comments {
				comment := fmt.Sprintf("\n\nComment:\nAuthor: %s\nBody: %s", comment.Username, comment.Body)
				commentTokens, err := calc.PreciseTokens(comment)
				if err != nil {
					log.Errorf("Error getting tokens: %s", err)
					os.Exit(1)
				}
				if tokens+commentTokens > calc.GetMaxTokens(c.Model) {
					break
				}
				prompt += comment
				tokens += commentTokens
			}
			prompt += suffix
		} else {
			// get the repo context here
			repoFiles, err := git.GetRepo(".", ".gptignore", false) // gptignore isn't working AGAIN
			if err != nil {
				log.Errorf("Error getting repo: %s", err)
				os.Exit(1)
			}

			m, _, err := memory.New(filepath.Join(".", ".index.memory"))
			if err != nil {
				log.Errorf("Error creating memory: %s", err)
				os.Exit(1)
			}

			for _, file := range repoFiles.Files {
				err = m.Add(file.Path, file.Contents)
				if err != nil {
					log.Warnf("Error adding file to memory: %s", err)
					continue
				}
			}

			results, err := m.Search(fmt.Sprintf("%s\n%s\n%s", title, body, question))
			if err != nil {
				log.Errorf("Error searching memory: %s", err)
				os.Exit(1)
			}

			m.Destroy()

			sorted := sortByScore(results)
			log.Debug("Getting file contents...")
			var files []g.GitFile
			for _, result := range sorted {
				for _, file := range repoFiles.Files {
					if file.Path == result.ID {
						files = append(files, file)
					}
				}
			}

			if len(files) == 0 {
				log.Error("No results found.")
				os.Exit(1)
			}

			// keep adding files until we hit the max tokens
			for _, file := range files {
				fileTokens, err := calc.PreciseTokens(file.Contents)
				if err != nil {
					log.Errorf("Error getting tokens: %s", err)
					os.Exit(1)
				}
				if tokens+fileTokens > calc.GetMaxTokens(c.Model) {
					break
				}
				prompt += fmt.Sprintf("\n\nFile: %s\n%s", file.Path, file.Contents)
				tokens += fileTokens
			}
		}

		if promptOnly {
			fmt.Println(prompt)
			os.Exit(0)
		}

		// ask ChatGPT
		resp, err := ai.SimpleRequest(prompt, c)
		if err != nil {
			log.Errorf("Error getting response: %s", err)
			os.Exit(1)
		}

		fmt.Println(resp)
	},
}

func init() {
	RootCmd.AddCommand(issueCmd)

	issueCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	issueCmd.Flags().IntVarP(&issuePRNumber, "number", "n", 0, "the number of the issue to get")
	issueCmd.Flags().StringVarP(&question, "question", "q", "", "the question to ask ChatGPT")
	issueCmd.Flags().BoolVarP(&useComments, "comments", "c", false, "use comments instead of git repo for context")
	issueCmd.Flags().BoolVarP(&promptOnly, "prompt-only", "p", false, "only generate a prompt, don't ask ChatGPT")
}
