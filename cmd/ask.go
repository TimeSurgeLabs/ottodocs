/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"os"

	"github.com/chand1012/git2gpt/prompt"
	"github.com/chand1012/memory"
	l "github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"

	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/git"
	"github.com/chand1012/ottodocs/pkg/utils"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask a question about a file or repo",
	Long: `Uses full text search to find relevant code and ask questions about said code.
Requires a path to a repository or file as a positional argument.`,
	Aliases: []string{"a"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(l.DebugLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var stream *openai.ChatCompletionStream
		var repoPath string
		var fileName string

		if len(args) > 0 {
			repoPath = args[0]
		} else {
			repoPath = "."
		}

		conf, err := config.Load()
		if err != nil || conf.APIKey == "" {
			// if the API key is not set, prompt the user to config
			log.Error("Please config first.")
			log.Error("Run `ottodocs config -h` to learn how to config.")
			os.Exit(1)
		}

		if chatPrompt == "" {
			log.Debug("User did not enter a question. Prompting for one...")
			chatPrompt, err = utils.Input("You: ")
			if err != nil {
				log.Errorf("Error prompting for question: %s", err)
				os.Exit(1)
			}
		}

		log.Debug("Getting file contents...")
		info, err := os.Stat(repoPath)
		if err != nil {
			log.Errorf("Error getting file info: %s", err)
			os.Exit(1)
		}
		// check if the first arg is a file or a directory
		// if it's a file, ask a question about that file directly
		if info.IsDir() {
			if !git.IsGitRepo(repoPath) {
				log.Error("Not a git repo.")
				os.Exit(1)
			}
			log.Debug("Constructing repo memory...")

			log.Debug("Creating memory index...")
			// Create a new memory index
			m, _, err := memory.New(":memory:")
			if err != nil {
				log.Errorf("Failed to create memory index: %s", err)
				os.Exit(1)
			}

			log.Debug("Indexing repo...")
			// index the repo
			repo, err := git.GetRepo(repoPath, ignoreFilePath, ignoreGitignore)
			if err != nil {
				log.Errorf("Error processing repo: %s", err)
				os.Exit(1)
			}

			// index the files
			for _, file := range repo.Files {
				err = m.Add(file.Path, file.Contents)
				if err != nil {
					log.Errorf("Error indexing file: %s", err)
					os.Exit(1)
				}
			}

			log.Debug("Searching memory index...")
			// search the memory index
			results, err := m.Search(chatPrompt)
			if err != nil {
				log.Errorf("Failed to search memory index: %s", err)
				os.Exit(1)
			}

			log.Debug("Results extracted. Destroying memory index...")
			// close the memory index
			m.Destroy()

			log.Debug("Sorting results...")
			sortedFragments := utils.SortByAverage(results)

			log.Debug("Getting file contents...")
			var files []prompt.GitFile
			for _, result := range sortedFragments {
				for _, file := range repo.Files {
					if file.Path == result.ID {
						files = append(files, file)
						log.Debugf("Found file: %s Score: %f Average: %f", file.Path, result.Score, result.Avg)
					}
				}
			}

			if len(files) == 0 {
				log.Error("No results found.")
				os.Exit(1)
			}

			log.Debug("Asking chatGPT question...")
			stream, err = ai.Question(files, chatPrompt, conf, verbose)
			if err != nil {
				log.Errorf("Error asking question: %s", err)
				os.Exit(1)
			}
		} else {
			log.Debug("Getting file contents...")
			fileName = repoPath
			content, err := utils.LoadFile(fileName)
			if err != nil {
				log.Errorf("Error loading file: %s", err)
				os.Exit(1)
			}

			log.Debug("Calculating tokens and constructing file...")
			tokens, err := calc.PreciseTokens(content)
			if err != nil {
				log.Errorf("Error calculating tokens: %s", err)
				os.Exit(1)
			}

			files := []prompt.GitFile{
				{
					Path:     fileName,
					Contents: content,
					Tokens:   int64(tokens),
				},
			}

			log.Debug("Asking chatGPT question...")
			stream, err = ai.Question(files, chatPrompt, conf, verbose)
			if err != nil {
				log.Errorf("Error asking question: %s", err)
				os.Exit(1)
			}
		}

		utils.PrintColoredText("Otto: ", conf.OttoColor)
		_, err = utils.PrintChatCompletionStream(stream)
		if err != nil {
			log.Errorf("Error printing chat completion stream: %s", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(askCmd)
	askCmd.Flags().StringVarP(&chatPrompt, "question", "q", "", "The question to ask")
	askCmd.Flags().BoolVarP(&ignoreGitignore, "ignore-gitignore", "g", false, "ignore .gitignore file")
	askCmd.Flags().StringVarP(&ignoreFilePath, "ignore", "n", "", "path to .gptignore file")
	askCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
