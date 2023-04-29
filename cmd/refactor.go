/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
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

// refactorCmd represents the refactor command
var refactorCmd = &cobra.Command{
	Use:   "refactor",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: []string{"r"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(l.DebugLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error("Requires a file name as an argument. Example: otto refactor main.go")
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
		if err != nil {
			log.Errorf("Error loading file: %s", err)
			os.Exit(1)
		}

		if endLine < 0 || startLine < 0 {
			log.Error("End line must be greater than or equal to start line and both must be greater than or equal to 0")
			os.Exit(1)
		}

		if endLine != 0 {
			// get the lines to refactor
			lines := strings.Split(contents, "\n")
			if endLine > len(lines) {
				log.Error("End line is greater than the number of lines in the file")
				os.Exit(1)
			}

			contents = strings.Join(lines[startLine-1:endLine], "\n")
		} else {
			endLine = len(strings.Split(contents, "\n"))
		}

		log.Debugf("Refactoring lines %d-%d", startLine, endLine)

		if chatPrompt == "" {
			chatPrompt, err = utils.Input("Goal: ")
			if err != nil {
				log.Errorf("Error prompting for goal: %s", err)
				os.Exit(1)
			}
		}

		prompt := constants.REFACTOR_CODE_PROMPT + "Goal: " + chatPrompt + "\n\n" + strings.TrimRight(contents, " \n")

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

		if !force {
			confirm, err := utils.Input("Would you like to overwrite the file with the new code? (y/N): ")
			if err != nil {
				log.Errorf("Error getting input: %s", err)
				os.Exit(1)
			}

			confirm = strings.ToLower(confirm)
			if confirm != "y" && confirm != "yes" {
				os.Exit(0)
			}
		}

		// write the new code to the file
		finalCode, err := textfile.ReplaceLines(contents, startLine, endLine, newCode)
		if err != nil {
			log.Errorf("Error replacing lines: %s", err)
			os.Exit(1)
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
	RootCmd.AddCommand(refactorCmd)

	refactorCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite of existing files")
	refactorCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	refactorCmd.Flags().IntVarP(&startLine, "start", "s", 1, "Start line")
	refactorCmd.Flags().IntVarP(&endLine, "end", "e", 0, "End line")
	refactorCmd.Flags().StringVarP(&chatPrompt, "goal", "g", "", "Goal of the refactor")
}
