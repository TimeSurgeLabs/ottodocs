/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/utils"
)

var VALID_MODELS = []string{"gpt-4", "gpt-4-0314", "gpt-4-32k", "gpt-4-32k-0314", "gpt-3.5-turbo", "gpt-3.5-turbo-0301"}

// setModelCmd represents the setModel command
var setModelCmd = &cobra.Command{
	Use:   "setModel",
	Short: "Set the model to use for documentation",
	Long: `Sets the model to use for documentation. Takes a valid ChatGPT API model name as a single positional argument.
Valid models are: gpt-4, gpt-4-0314, gpt-4-32k, gpt-4-32k-0314, gpt-3.5-turbo, gpt-3.5-turbo-0301
See here for more information: https://platform.openai.com/docs/models/model-endpoint-compatibility
`,
	Run: func(cmd *cobra.Command, args []string) {
		model := args[0]
		c, err := config.Load()
		if err != nil {
			log.Errorf("Error loading config: %s", err)
		}
		if !utils.Contains(VALID_MODELS, model) {
			log.Errorf("Invalid model: %s", model)
			log.Errorf("Valid models are: %s", VALID_MODELS)
			return
		}

		c.Model = model

		err = config.Save(c)
		if err != nil {
			log.Errorf("Error saving config: %s", err)
		}
		log.Infof("Set model to %s", model)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	RootCmd.AddCommand(setModelCmd)
}
