/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chand1012/ottodocs/pkg/ai"
	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/constants"
	"github.com/chand1012/ottodocs/pkg/textfile"
	"github.com/chand1012/ottodocs/pkg/utils"
	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a file using AI",
	Long: `OttoDocs Edit allows you to use AI to help edit your code files. 
Provide a file name and a goal, and OttoDocs will return a generated version of the file.
You can even specify the starting and ending lines for the edit, or choose to append the results to the file:

Example: otto edit main.go --start 1 --end 10 --goal "Refactor the function"`,
	Aliases: []string{"e"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(l.DebugLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error("Requires a file name as an argument. Example: otto edit main.go")
			os.Exit(1)
		}

		c, err := config.Load()
		if err != nil {
			log.Errorf("Error loading config: %s", err)
			os.Exit(1)
		}

		fileName := args[0]

		// load the file
		contents, err := utils.LoadFile(fileName)
		if os.IsNotExist(err) {
			contents = ""
		} else if err != nil {
			log.Errorf("Error loading file: %s", err)
			os.Exit(1)
		}

		if endLine < 0 || startLine < 0 {
			log.Error("End line must be greater than or equal to start line and both must be greater than or equal to 0")
			os.Exit(1)
		}

		var editCode string
		if endLine != 0 {
			// get the lines to edit
			lines := strings.Split(contents, "\n")
			if endLine > len(lines) {
				log.Error("End line is greater than the number of lines in the file")
				os.Exit(1)
			}

			editCode = strings.Join(lines[startLine-1:endLine], "\n")
		} else {
			endLine = len(strings.Split(contents, "\n"))
		}

		log.Debugf("editing lines %d-%d", startLine, endLine)

		if chatPrompt == "" {
			chatPrompt, err = utils.Input("Goal: ")
			if err != nil {
				log.Errorf("Error prompting for goal: %s", err)
				os.Exit(1)
			}
		}

		var prompt string
		if editCode != "" {
			prompt = constants.EDIT_CODE_PROMPT + "\nEDIT: " + editCode + "\n\nGOAL: " + chatPrompt + "\n\nFILE: " + filePath + "\n\n" + contents
		} else {
			prompt = constants.EDIT_CODE_PROMPT + "\nGOAL: " + chatPrompt + "\n\nFILE: " + filePath + "\n\n" + contents
		}

		if len(contextFiles) > 0 {
			var contextContent string
			for _, contextFile := range contextFiles {
				contextContent, err = utils.LoadFile(contextFile)
				if err != nil {
					log.Errorf("Error loading context file: %s", err)
					continue
				}
				prompt += "\n\nCONTEXT: " + contextFile + "\n\n" + contextContent
			}
		}

		stream, err := ai.SimpleStreamRequest(prompt, c)
		if err != nil {
			log.Errorf("Error requesting from OpenAI: %s", err)
			os.Exit(1)
		}

		// print the response
		fmt.Println("New Code:")
		newCode, err := utils.PrintChatCompletionStream(stream)
		if err != nil {
			log.Errorf("Error printing chat completion stream: %s", err)
			os.Exit(1)
		}

		confirmMsg := "Would you like to overwrite the file with the new code? (y/N): "
		if appendFile {
			confirmMsg = "Would you like to append the new code to the file? (y/N): "
		}

		if !force {
			confirm, err := utils.Input(confirmMsg)
			if err != nil {
				log.Errorf("Error getting input: %s", err)
				os.Exit(1)
			}

			confirm = strings.ToLower(confirm)
			if confirm != "y" && confirm != "yes" {
				os.Exit(0)
			}
		}

		var finalCode string
		if appendFile {
			// add the new code to the end of the file
			finalCode = contents + "\n" + newCode + "\n"
		} else {
			// write the new code to the file
			finalCode, err = textfile.ReplaceLines(contents, startLine, endLine, newCode)
			if err != nil {
				log.Errorf("Error replacing lines: %s", err)
				os.Exit(1)
			}
		}

		err = utils.WriteFile(fileName, finalCode)
		if err != nil {
			log.Errorf("Error writing file: %s", err)
			os.Exit(1)
		}
		fmt.Println("File written successfully!")
	},
}

func init() {
	RootCmd.AddCommand(editCmd)

	editCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite of existing files")
	editCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	editCmd.Flags().BoolVarP(&appendFile, "append", "a", false, "Append to the end of a file instead of overwriting it")
	editCmd.Flags().IntVarP(&startLine, "start", "s", 1, "Start line")
	editCmd.Flags().IntVarP(&endLine, "end", "e", 0, "End line")
	editCmd.Flags().StringVarP(&chatPrompt, "goal", "g", "", "Goal of the edit")
	editCmd.Flags().StringSliceVarP(&contextFiles, "context", "c", []string{}, "Context files")
}
