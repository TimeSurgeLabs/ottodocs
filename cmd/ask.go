/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chand1012/memory"
	"github.com/spf13/cobra"

	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/utils"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask a question about a file or repo",
	Long: `Uses full text search to find relevant code and ask questions about said code.
Requires a path to a repository or file as a positional argument.`,
	Args: cobra.PositionalArgs(func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires a path to a repository or file")
		}
		return nil
	}),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath := args[0]
		var fileName string
		conf, err := config.Load()
		if err != nil || conf.APIKey == "" {
			// if the API key is not set, prompt the user to login
			log.Error("Please login first.")
			log.Error("Run `ottodocs login` to login.")
			os.Exit(1)
		}

		if chatPrompt == "" {
			fmt.Println("Please enter a question: ")
			fmt.Scanln(&chatPrompt)
			// strip the newline character
			chatPrompt = strings.TrimRight(chatPrompt, " \n")
		}

		info, err := os.Stat(repoPath)
		if err != nil {
			log.Errorf("Error getting file info: %s", err)
			os.Exit(1)
		}
		// check if the first arg is a file or a directory
		// if it's a file, ask a question about that file directly
		if info.IsDir() {
			// Define a temporary path for the index file
			testIndexPath := filepath.Join(args[0], ".index.memory")

			// Create a new memory index
			m, _, err := memory.New(testIndexPath)
			if err != nil {
				log.Errorf("Failed to create memory index: %s", err)
				os.Exit(1)
			}

			// index the repo
			repo, err := utils.GetRepo(repoPath, ignoreFilePath, ignoreGitignore)
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

			// search the memory index
			results, err := m.Search(chatPrompt)
			if err != nil {
				log.Errorf("Failed to search memory index: %s", err)
				os.Exit(1)
			}

			m.Destroy()

			// get the top fragment by average score
			top := memory.TopFragmentAvg(results)
			fileName = top.ID
		} else {
			fileName = repoPath
		}

		fmt.Println("Asking question about " + fileName + "...")

		content, err := utils.LoadFile(fileName)
		if err != nil {
			log.Errorf("Error loading file: %s", err)
			os.Exit(1)
		}

		resp, err := ai.Question(fileName, content, chatPrompt, conf.APIKey, conf.Model)

		if err != nil {
			log.Errorf("Error asking question: %s", err)
			os.Exit(1)
		}

		fmt.Println(resp)
	},
}

func init() {
	RootCmd.AddCommand(askCmd)
	askCmd.Flags().StringVarP(&chatPrompt, "question", "q", "", "The question to ask")
	askCmd.Flags().BoolVarP(&ignoreGitignore, "ignore-gitignore", "g", false, "ignore .gitignore file")
	askCmd.Flags().StringVarP(&ignoreFilePath, "ignore", "n", "", "path to .gptignore file")
}
