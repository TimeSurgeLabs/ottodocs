/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/TimeSurgeLabs/ottodocs/pkg/calc"
	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/constants"
	"github.com/TimeSurgeLabs/ottodocs/pkg/git"
	"github.com/TimeSurgeLabs/ottodocs/pkg/textfile"
	"github.com/TimeSurgeLabs/ottodocs/pkg/utils"
	"github.com/chand1012/memory"
	l "github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
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

		var messages []openai.ChatCompletionMessage
		var newCode string

		var prompt string
		if editCode != "" {
			prompt = "EDIT: " + editCode + "\n\nGOAL: " + chatPrompt + "\n\nFILE: " + filePath + "\n\n" + contents + "\n\nBe sure to only output the edited code, do not print the entire file."
		} else {
			prompt = "GOAL: " + chatPrompt + "\n\nFILE: " + filePath + "\n\n" + contents
		}

		client := openai.NewClient(c.APIKey)

		if repoContext {
			repo, err := git.GetRepo(".", "", false)
			if err != nil {
				log.Errorf("Error getting repo: %s", err)
				os.Exit(1)
			}

			m, _, err := memory.New(":memory:")
			if err != nil {
				log.Errorf("Error creating memory: %s", err)
				os.Exit(1)
			}

			for _, file := range repo.Files {
				if file.Path == fileName {
					continue
				}
				err = m.Add(file.Path, file.Contents)
				if err != nil {
					log.Errorf("Error indexing file: %s", err)
					os.Exit(1)
				}
			}

			queryConstructorPrompt := []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Convert this goal into a terms to use for a search query. They do not have to be organized, nor a complete sentence. Use no form of punctuation or quotations. Only return the query and nothing else: " + chatPrompt,
				},
			}

			utils.PrintColoredText("Otto: ", c.OttoColor)
			fmt.Println("Ok! Here is the query, taking your input into account.")
			utils.PrintColoredText("Otto: ", c.OttoColor)
			stream, err := client.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
				Model:    "gpt-3.5-turbo",
				Messages: queryConstructorPrompt,
			})

			if err != nil {
				log.Errorf("Error requesting from OpenAI: %s", err)
				os.Exit(1)
			}

			query, err := utils.PrintChatCompletionStream(stream)
			if err != nil {
				log.Errorf("Error printing chat completion stream: %s", err)
				os.Exit(1)
			}

			utils.PrintColoredText("Otto: ", c.OttoColor)
			fmt.Println("Searching repo for files that match the query...")

			resp, err := m.Search(query)
			if err != nil {
				log.Errorf("Error searching memory: %s", err)
				os.Exit(1)
			}

			contextFiles = []string{}
			for _, file := range resp {
				contextFiles = append(contextFiles, file.ID)
			}
			log.Debugf("context files: %s", contextFiles)
		}
		tokens, err := calc.PreciseTokens(prompt)
		if err != nil {
			log.Errorf("Error calculating tokens: %s", err)
			os.Exit(1)
		}

		maxTokens := calc.GetMaxTokens(c.Model)
		if len(contextFiles) > 0 {
			var contextContent string
			for _, contextFile := range contextFiles {
				contextContent, err = utils.LoadFile(contextFile)
				if err != nil {
					log.Errorf("Error loading context file: %s", err)
					continue
				}
				contentTokens, err := calc.PreciseTokens(contextContent)
				if err != nil {
					log.Errorf("Error loading context file: %s", err)
					continue
				}
				if contentTokens+tokens > maxTokens {
					break
				}
				prompt += "\n\nCONTEXT: " + contextFile + "\n\n" + contextContent
				tokens += contentTokens
			}
		}

		log.Debugf("prompt: %s", prompt)

		messages = []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: constants.EDIT_CODE_PROMPT,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		}

		for {

			stream, err := client.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
				Model:    c.Model,
				Messages: messages,
			})

			if err != nil {
				log.Errorf("Error requesting from OpenAI: %s", err)
				os.Exit(1)
			}

			// print the response
			utils.PrintColoredTextLn("New Code:", c.OttoColor)
			newCode, err = utils.PrintChatCompletionStream(stream)
			if err != nil {
				log.Errorf("Error printing chat completion stream: %s", err)
				os.Exit(1)
			}

			confirmMsg := "Would you like to write the file with the new code? (y/N). Type your input to keep editing: "
			if appendFile {
				confirmMsg = "Would you like to append the new code to the file? (y/N). Type your input to keep editing: "
			}

			if !force {
				confirm, err := utils.Input(confirmMsg)
				if err != nil {
					log.Errorf("Error getting input: %s", err)
					os.Exit(1)
				}

				confirm = strings.ToLower(confirm)
				if confirm == "n" || confirm == "no" {
					os.Exit(0)
				} else if confirm == "y" || confirm == "yes" {
					break
				} else {
					codeTokens, err := calc.PreciseTokens(newCode)
					if err != nil {
						log.Errorf("Error calculating tokens: %s", err)
						os.Exit(1)
					}

					maxTokens := calc.GetMaxTokens(c.Model) - codeTokens

					var newMessages []openai.ChatCompletionMessage
					newMessages = []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: "Use the following input to edit the code: " + confirm + "\n\nMake sure to only output the code, do not print anything else.",
						},
						{
							Role:    openai.ChatMessageRoleAssistant,
							Content: newCode,
						},
					}
					utils.ReverseSlice(messages)
					for _, message := range messages {
						if calc.PreciseTokensFromMessages(newMessages, c.Model) < maxTokens {
							newMessages = append(newMessages, message)
						}
					}
					utils.ReverseSlice(newMessages)
					messages = newMessages
					utils.PrintColoredText("Otto: ", c.OttoColor)
					fmt.Println("Ok! Here is the new code, taking your input into account.")
				}
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
	editCmd.Flags().BoolVarP(&repoContext, "repo", "r", false, "Use the current repo as context")
	editCmd.Flags().IntVarP(&startLine, "start", "s", 1, "Start line")
	editCmd.Flags().IntVarP(&endLine, "end", "e", 0, "End line")
	editCmd.Flags().StringVarP(&chatPrompt, "goal", "g", "", "Goal of the edit")
	editCmd.Flags().StringSliceVarP(&contextFiles, "context", "c", []string{}, "Context files")
}
