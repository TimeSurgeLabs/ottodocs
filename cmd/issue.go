/*
Copyright © 2024 TimeSurgeLabs <chandler@timesurgelabs.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	g "github.com/chand1012/git2gpt/prompt"
	"github.com/chand1012/memory"
	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/TimeSurgeLabs/ottodocs/pkg/ai"
	"github.com/TimeSurgeLabs/ottodocs/pkg/calc"
	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/gh"
	"github.com/TimeSurgeLabs/ottodocs/pkg/git"
	"github.com/TimeSurgeLabs/ottodocs/pkg/utils"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:     "issue",
	Short:   "Get a prompt for or ask Otto about a GitHub Issue.",
	Long:    `Get a prompt for or ask Otto about a GitHub Issue. Uses the current working directory's repo by default.`,
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
			log.Debug("No issue number provided")
			// prompt for issue number
			inputPRNumber, err := utils.Input("Please provide an issue number: ")
			if err != nil {
				log.Errorf("Error getting input: %s", err)
				os.Exit(1)
			}
			issuePRNumber, err = strconv.Atoi(inputPRNumber)
			if err != nil {
				log.Errorf("Error converting input to int: %s", err)
				os.Exit(1)
			}
		}

		if question == "" {
			log.Debug("No question provided")
			// prompt for question
			question, err = utils.Input("Please provide a question: ")
			if err != nil {
				log.Errorf("Error getting input: %s", err)
				os.Exit(1)
			}
		}

		log.Debugf("Question: %s", question)

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

		log.Debug("Checking initial prompt length...")
		if tokens > calc.GetMaxTokens(c.Model) {
			log.Errorf("Error: prompt is too long. Max tokens: %d, prompt tokens: %d", calc.GetMaxTokens(c.Model), tokens)
			os.Exit(1)
		}

		if useComments {
			log.Debug("Using comments...")
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
		} else if len(contextFiles) > 0 {
			log.Debug("Using specified files...")
			var contextContent string
			for _, contextFile := range contextFiles {
				contextContent, err = utils.LoadFile(contextFile)
				if err != nil {
					log.Errorf("Error loading context file: %s", err)
					continue
				}
				prompt += "\n\nFile: " + contextFile + "\n\n" + contextContent
				fileTokens, err := calc.PreciseTokens(contextContent)
				if err != nil {
					log.Errorf("Error getting tokens: %s", err)
					os.Exit(1)
				}
				if tokens+fileTokens > calc.GetMaxTokens(c.Model) {
					break
				}
				tokens += fileTokens
			}
		} else {
			log.Debug("Using repo contents...")
			// get the repo context here
			repoFiles, err := git.GetRepo(".", ".gptignore", false) // gptignore isn't working AGAIN
			if err != nil {
				log.Errorf("Error getting repo: %s", err)
				os.Exit(1)
			}

			log.Debug("Constructing memory...")
			m, _, err := memory.New(filepath.Join(".", ".index.memory"))
			if err != nil {
				log.Errorf("Error creating memory: %s", err)
				os.Exit(1)
			}

			log.Debug("Adding files to memory...")
			for _, file := range repoFiles.Files {
				err = m.Add(file.Path, file.Contents)
				if err != nil {
					log.Warnf("Error adding file to memory: %s", err)
					continue
				}
			}

			log.Debug("Searching memory...")
			results, err := m.Search(fmt.Sprintf("%s\n%s\n%s", title, body, question))
			if err != nil {
				log.Errorf("Error searching memory: %s", err)
				os.Exit(1)
			}

			log.Debug("Destroying memory...")
			m.Destroy()

			log.Debug("Sorting results...")
			sorted := utils.SortByAverage(results)
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

			log.Debug("Adding files to prompt...")
			// keep adding files until we hit the max tokens
			for _, file := range files {
				promptAdd := fmt.Sprintf("\n\nFile: %s\n%s", file.Path, file.Contents)
				fileTokens, err := calc.PreciseTokens(promptAdd)
				if err != nil {
					log.Errorf("Error getting tokens: %s", err)
					os.Exit(1)
				}
				if tokens+fileTokens > calc.GetMaxTokens(c.Model) {
					break
				}
				prompt += promptAdd
				tokens += fileTokens
			}
		}

		prompt += suffix

		if countFinalTokens {
			log.Debug("Counting final prompt tokens...")
			tokens, err := calc.PreciseTokens(prompt)
			if err != nil {
				log.Errorf("Error getting tokens: %s", err)
				os.Exit(1)
			}
			if verbose {
				log.Debugf("Final prompt tokens: %d\n", tokens)
			} else {
				fmt.Printf("Final prompt tokens: %d\n", tokens)
			}
		}

		if promptOnly {
			log.Debug("Prompt only mode, printing prompt and exiting.")
			fmt.Println(prompt)
			os.Exit(0)
		}

		log.Debug("Asking Otto...")
		// ask Otto
		stream, err := ai.SimpleStreamRequest(prompt, c)
		if err != nil {
			log.Errorf("Error getting response: %s", err)
			os.Exit(1)
		}

		fmt.Print("Otto: ")
		_, err = utils.PrintChatCompletionStream(stream)
		if err != nil {
			log.Errorf("Error printing response: %s", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(issueCmd)

	issueCmd.Flags().IntVarP(&issuePRNumber, "number", "n", 0, "the number of the issue to get")
	issueCmd.Flags().StringVarP(&question, "question", "q", "", "the question to ask Otto")
	issueCmd.Flags().StringSliceVarP(&contextFiles, "context", "c", []string{}, "the files to use as context")
	issueCmd.Flags().BoolVarP(&useComments, "comments", "r", false, "use comments instead of git repo for context")
	issueCmd.Flags().BoolVarP(&promptOnly, "prompt-only", "p", false, "only generate a prompt, don't ask Otto")
	issueCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	issueCmd.Flags().BoolVar(&countFinalTokens, "count", false, "count the number of tokens")
}
