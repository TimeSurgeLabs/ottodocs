/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chand1012/git2gpt/prompt"
	"github.com/chand1012/ottodocs/config"
	"github.com/chand1012/ottodocs/document"
	"github.com/spf13/cobra"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Document a repository of files",
	Long: `Document an entire repository of files. Specify the path to the repo as the first positional argument. This command will recursively
search for files in the directory and document them.
	`,
	Args: cobra.PositionalArgs(func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires a path to a repository")
		}
		return nil
	}),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath = args[0]

		if (!markdownMode || inlineMode) && outputFile != "" {
			fmt.Println("Error: cannot specify an output file in inline mode")
			os.Exit(1)
		}

		if markdownMode && overwriteOriginal {
			fmt.Println("Error: cannot overwrite original file in markdown mode")
			os.Exit(1)
		}

		if markdownMode && outputFile == "" {
			fmt.Println("Error: must specify an output file in markdown mode")
			os.Exit(1)
		}

		if outputFile != "" {
			// if output file exists, throw error
			if _, err := os.Stat(outputFile); err == nil {
				fmt.Printf("Error: output file %s already exists!\n", outputFile)
				os.Exit(1)
			}
		}

		conf, err := config.Load()
		if err != nil || conf.APIKey == "" {
			// if the API key is not set, prompt the user to login
			fmt.Println("Please login first.")
			fmt.Println("Run `ottodocs login` to login.")
			os.Exit(1)
		}

		ignoreList := prompt.GenerateIgnoreList(ignoreFilePath, ignoreFilePath, !ignoreGitignore)
		ignoreList = append(ignoreList, filepath.Join(repoPath, ".gptignore"))
		repo, err := prompt.ProcessGitRepo(repoPath, ignoreList)
		if err != nil {
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}

		for _, file := range repo.Files {
			var contents string

			path := filepath.Join(repoPath, file.Path)

			if outputFile != "" {
				fmt.Println("Documenting file", file.Path)
			}

			if chatPrompt == "" {
				chatPrompt = "Write documentation for the following code snippet. The file name is" + file.Path + ":"
			}

			if inlineMode || !markdownMode {
				contents, err = document.SingleFile(path, chatPrompt, conf.APIKey)
			} else {
				contents, err = document.Markdown(path, chatPrompt, conf.APIKey)
			}

			if err != nil {
				fmt.Printf("Error documenting file %s: %s\n", path, err)
				continue
			}

			if outputFile != "" && markdownMode {
				// write the string to the output file
				// append if the file already exists
				file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					os.Exit(1)
				}

				_, err = file.WriteString(contents)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					os.Exit(1)
				}

				file.Close()
			} else if overwriteOriginal {
				// overwrite the original file
				// clear the contents of the file
				file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					os.Exit(1)
				}

				// write the new contents to the file
				_, err = file.WriteString(contents)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					os.Exit(1)
				}

				file.Close()
			} else {
				fmt.Println(contents)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(docsCmd)

	// see cmd/vars for the definition of these flags
	docsCmd.Flags().StringVarP(&chatPrompt, "prompt", "p", "", "Prompt to use for the ChatGPT API")
	docsCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Path to the output file. For use with --markdown")
	docsCmd.Flags().StringVarP(&ignoreFilePath, "ignore", "n", "", "path to .gptignore file")
	docsCmd.Flags().BoolVarP(&markdownMode, "markdown", "m", false, "Output in Markdown format")
	docsCmd.Flags().BoolVarP(&inlineMode, "inline", "i", false, "Output in inline format")
	docsCmd.Flags().BoolVarP(&overwriteOriginal, "overwrite", "w", false, "Overwrite the original file")
	docsCmd.Flags().BoolVarP(&ignoreGitignore, "ignore-gitignore", "g", false, "ignore .gitignore file")
}
