/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
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
			var index bleve.Index
			var err error
			mapping := bleve.NewIndexMapping()

			// check if .index.bleve exists
			// if it does, load it
			// if it doesn't, create it
			// fmt.Println("Indexing repo...")
			index, err = bleve.New(filepath.Join(args[0], ".index.bleve"), mapping)
			if err != nil {
				log.Errorf("Error creating index: %s", err)
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
				err = index.Index(file.Path, file)
				if err != nil {
					log.Errorf("Error indexing file: %s", err)
					os.Exit(1)
				}
			}

			// ask a question about the repo
			query := bleve.NewQueryStringQuery(chatPrompt)
			search := bleve.NewSearchRequest(query)
			searchResults, err := index.Search(search)
			if err != nil {
				log.Errorf("Error searching index: %s", err)
				os.Exit(1)
			}

			// tokenize the question
			tokens := utils.SimpleTokenize(chatPrompt)
			for _, token := range tokens {
				query := bleve.NewQueryStringQuery(token)
				search := bleve.NewSearchRequest(query)
				r, err := index.Search(search)
				if err != nil {
					log.Errorf("Error searching index: %s", err)
					os.Exit(1)
				}
				// combines the searches, but doesn't weight by ID
				searchResults.Merge(r)
			}
			hits := make(map[string]float64)

			for _, results := range searchResults.Hits {
				currentScore := hits[results.ID]
				hits[results.ID] = currentScore + results.Score
			}

			var bestHit string
			var bestScore float64
			for hit, score := range hits {
				if score > bestScore {
					bestScore = score
					bestHit = hit
				}
			}

			fileName = bestHit
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

		// if .index.bleve exists, delete it
		if _, err := os.Stat(filepath.Join(args[0], ".index.bleve")); err == nil {
			err = os.RemoveAll(filepath.Join(args[0], ".index.bleve"))
			if err != nil {
				log.Errorf("Error deleting index: %s", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(askCmd)
	askCmd.Flags().StringVarP(&chatPrompt, "question", "q", "", "The question to ask")
	askCmd.Flags().BoolVarP(&ignoreGitignore, "ignore-gitignore", "g", false, "ignore .gitignore file")
	askCmd.Flags().StringVarP(&ignoreFilePath, "ignore", "n", "", "path to .gptignore file")
}
