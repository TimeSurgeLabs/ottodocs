/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/git"
	"github.com/chand1012/ottodocs/pkg/utils"
	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Document a repository of files or a single file",
	Long: `Document an entire repository of files. Specify the path to the repo as the first positional argument. This command will recursively
search for files in the directory and document them. If a single file is specified, it will be documented.

Example:
otto docs . -i -w 
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(l.DebugLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var repoPath string
		if len(args) > 0 {
			repoPath = args[0]
		} else {
			repoPath = "."
		}

		if markdownMode && overwriteOriginal {
			log.Error("Error: cannot overwrite original file in markdown mode")
			os.Exit(1)
		}

		if markdownMode && outputFile == "" {
			log.Error("Error: must specify an output file in markdown mode")
			os.Exit(1)
		}

		if outputFile != "" {
			// if output file exists, throw error
			if _, err := os.Stat(outputFile); err == nil {
				log.Errorf("Error: output file %s already exists!", outputFile)
				os.Exit(1)
			}
		}

		conf, err := config.Load()
		if err != nil || conf.APIKey == "" {
			// if the API key is not set, prompt the user to config
			log.Error("Please config first.")
			log.Error("Run `ottodocs config -h` to learn how to config.")
			os.Exit(1)
		}

		log.Info("Getting path info...")
		info, err := os.Stat(repoPath)
		if err != nil {
			log.Errorf("Error getting file info: %s", err)
			os.Exit(1)
		}

		if info.IsDir() {
			// make sure its a git repo
			if !git.IsGitRepo(repoPath) {
				log.Error("Error: not a git repository")
				os.Exit(1)
			}

			log.Info("Getting repo...")
			repo, err := git.GetRepo(repoPath, ignoreFilePath, ignoreGitignore)
			if err != nil {
				log.Errorf("Error: %s", err)
				os.Exit(1)
			}

			log.Debug("Documenting repo...")
			for _, file := range repo.Files {
				var contents string

				path := filepath.Join(repoPath, file.Path)

				if outputFile != "" {
					fmt.Println("Documenting file", file.Path)
				}

				if chatPrompt == "" {
					chatPrompt = "Write documentation for the following code snippet. The file name is" + file.Path + ":"
				}

				log.Debugf("Loading file %s", path)
				fileContents, err := utils.LoadFile(path)
				if err != nil {
					log.Warnf("Error loading file %s: %s", path, err)
					continue
				}

				if inlineMode || !markdownMode {
					log.Debugf("Documenting inline file %s", path)
					contents, err = ai.SingleFile(path, fileContents, chatPrompt, conf)
				} else {
					log.Debugf("Documenting markdown for %s", path)
					contents, err = ai.Markdown(path, fileContents, chatPrompt, conf)
				}

				if err != nil {
					log.Warnf("Error documenting file %s: %s", path, err)
					continue
				}

				if outputFile != "" && markdownMode {
					log.Debug("Writing markdown to file...")
					// write the string to the output file
					// append if the file already exists
					file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
					if err != nil {
						log.Errorf("Error: %s", err)
						os.Exit(1)
					}

					_, err = file.WriteString(contents)
					if err != nil {
						log.Errorf("Error: %s", err)
						os.Exit(1)
					}

					file.Close()
				} else if overwriteOriginal {
					log.Debug("Overwriting original file...")
					// overwrite the original file
					// clear the contents of the file
					file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
					if err != nil {
						log.Errorf("Error: %s", err)
						os.Exit(1)
					}

					// write the new contents to the file
					_, err = file.WriteString(contents)
					if err != nil {
						log.Errorf("Error: %s", err)
						os.Exit(1)
					}

					file.Close()
				} else {
					// print the contents to stdout
					fmt.Println(contents)
				}
			}
		} else {
			log.Info("Documenting file...")
			if chatPrompt == "" {
				chatPrompt = "Write documentation for the following code snippet:"
			}

			filePath := repoPath

			var contents string

			fileContents, err := utils.LoadFile(filePath)
			if err != nil {
				log.Errorf("Error: %s", err)
				os.Exit(1)
			}

			if inlineMode || !markdownMode {
				log.Debug("Documenting inline...")
				contents, err = ai.SingleFile(filePath, fileContents, chatPrompt, conf)
			} else {
				log.Debug("Documenting markdown...")
				contents, err = ai.Markdown(filePath, fileContents, chatPrompt, conf)
			}

			if err != nil {
				log.Errorf("Error: %s", err)
				os.Exit(1)
			}

			if outputFile != "" {
				log.Debug("Writing to file...")
				// write the string to the output file
				err = os.WriteFile(outputFile, []byte(contents), 0644)
				if err != nil {
					log.Errorf("Error: %s", err)
					os.Exit(1)
				}
			} else if overwriteOriginal {
				log.Debug("Overwriting original file...")
				// overwrite the original file
				// clear the contents of the file
				file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					log.Errorf("Error: %s", err)
					os.Exit(1)
				}

				// write the new contents to the file
				_, err = file.WriteString(contents)
				if err != nil {
					log.Errorf("Error: %s", err)
					os.Exit(1)
				}
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
	docsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}
