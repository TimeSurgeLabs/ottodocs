/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/chand1012/git2gpt/prompt"
	"github.com/spf13/cobra"
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Generates a ChatGPT prompt from a given Git repo",
	Long:  `Generates a ChatGPT prompt from a given Git repo. Specify the path to the repo as the first positional argument.`,
	Args: cobra.PositionalArgs(func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires a path to a repository")
		}
		return nil
	}),
	Run: func(cmd *cobra.Command, args []string) {
		repoPath = args[0]
		ignoreList := prompt.GenerateIgnoreList(repoPath, ignoreFilePath, !ignoreGitignore)
		repo, err := prompt.ProcessGitRepo(repoPath, ignoreList)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		if outputJSON {
			output, err := prompt.MarshalRepo(repo)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			if outputFile != "" {
				// if output file exists, throw error
				if _, err := os.Stat(outputFile); err == nil {
					fmt.Printf("Error: output file %s already exists\n", outputFile)
					os.Exit(1)
				}
				err = os.WriteFile(outputFile, []byte(output), 0644)
				if err != nil {
					fmt.Printf("Error: could not write to output file %s\n", outputFile)
					os.Exit(1)
				}
			} else {
				fmt.Println(string(output))
			}
			return
		}
		output, err := prompt.OutputGitRepo(repo, preambleFile)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		if outputFile != "" {
			// if output file exists, throw error
			if _, err := os.Stat(outputFile); err == nil {
				fmt.Printf("Error: output file %s already exists\n", outputFile)
				os.Exit(1)
			}
			err = os.WriteFile(outputFile, []byte(output), 0644)
			if err != nil {
				fmt.Printf("Error: could not write to output file %s\n", outputFile)
				os.Exit(1)
			}
		} else {
			fmt.Println(output)
		}
		if estimateTokens {
			fmt.Printf("Estimated number of tokens: %d\n", prompt.EstimateTokens(output))
		}
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)

	// see cmd/vars.go for the definition of these flags
	promptCmd.Flags().StringVarP(&preambleFile, "preamble", "p", "", "path to preamble text file")
	// output to file flag. Should be a string
	promptCmd.Flags().StringVarP(&outputFile, "output", "o", "", "path to output file")
	// estimate tokens. Should be a bool
	promptCmd.Flags().BoolVarP(&estimateTokens, "estimate", "e", false, "estimate the number of tokens in the output")
	// ignore file path. Should be a string
	promptCmd.Flags().StringVarP(&ignoreFilePath, "ignore", "i", "", "path to .gptignore file")
	// ignore gitignore. Should be a bool
	promptCmd.Flags().BoolVarP(&ignoreGitignore, "ignore-gitignore", "g", false, "ignore .gitignore file")
	// output JSON. Should be a bool
	promptCmd.Flags().BoolVarP(&outputJSON, "json", "j", false, "output JSON")
}
