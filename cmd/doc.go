/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/chand1012/ottodocs/ai"
	"github.com/chand1012/ottodocs/config"
	"github.com/chand1012/ottodocs/utils"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Document a file",
	Long:  `Document a file using the OpenAI ChatGPT API.`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			log.Error("Please provide a path to a file to document.")
			log.Error("Run `ottodocs doc -h` for more information.")
			os.Exit(1)
		}

		if chatPrompt == "" {
			chatPrompt = "Write documentation for the following code snippet:"
		}

		conf, err := config.Load()
		if err != nil || conf.APIKey == "" {
			// if the API key is not set, prompt the user to login
			log.Error("Please login first.")
			log.Error("Run `ottodocs login` to login.")
			os.Exit(1)
		}

		var contents string

		fileContents, err := utils.LoadFile(filePath)
		if err != nil {
			log.Errorf("Error: %s", err)
			os.Exit(1)
		}

		if inlineMode || !markdownMode {
			contents, err = ai.SingleFile(filePath, fileContents, chatPrompt, conf.APIKey, conf.Model)
		} else {
			contents, err = ai.Markdown(filePath, fileContents, chatPrompt, conf.APIKey, conf.Model)
		}

		if err != nil {
			log.Errorf("Error: %s", err)
			os.Exit(1)
		}

		if outputFile != "" {
			// write the string to the output file
			err = os.WriteFile(outputFile, []byte(contents), 0644)
			if err != nil {
				log.Errorf("Error: %s", err)
				os.Exit(1)
			}
		} else if overwriteOriginal {
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
	},
}

func init() {
	RootCmd.AddCommand(docCmd)

	// see cmd/vars.go for the definition of these flags
	docCmd.Flags().StringVarP(&filePath, "file", "f", "", "The file to document")
	docCmd.Flags().StringVarP(&chatPrompt, "prompt", "p", "", "The prompt to use for the document")
	docCmd.Flags().StringVarP(&outputFile, "output", "o", "", "The output file to write the documentation to. ")
	docCmd.Flags().BoolVarP(&inlineMode, "inline", "i", false, "Inline mode. Adds the documentation to the code.")
	docCmd.Flags().BoolVarP(&markdownMode, "markdown", "m", false, "Markdown mode. Outputs the documentation in markdown.")
	docCmd.Flags().BoolVarP(&overwriteOriginal, "overwrite", "w", false, "Overwrite the original file.")
}
