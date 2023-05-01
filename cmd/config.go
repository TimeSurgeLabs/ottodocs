/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/chand1012/ottodocs/pkg/utils"
	"github.com/spf13/cobra"
)

var VALID_MODELS = []string{"gpt-4", "gpt-4-0314", "gpt-4-32k", "gpt-4-32k-0314", "gpt-3.5-turbo", "gpt-3.5-turbo-0301"}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configures ottodocs",
	Long: `Configures ottodocs. Allows user to specify OpenAI API Key, GitHub Token, and the model with a single command.

Default model is gpt-3.5-turbo.
Valid models are: gpt-4, gpt-4-0314, gpt-4-32k, gpt-4-32k-0314, gpt-3.5-turbo, gpt-3.5-turbo-0301
See here for more information on the Models available: https://platform.openai.com/docs/models/model-endpoint-compatibility	

GitHub Tokens need access to the repo scope.

OpenAI API Key Generation: https://platform.openai.com/account/api-keys
GitHub Token Generation: https://github.com/settings/tokens
`,
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		// load the config
		c, err := config.Load()
		if err != nil {
			log.Errorf("Error loading config: %s", err)
			os.Exit(1)
		}

		// if none of the config options are provided, print a warning
		if apiKey == "" && model == "" && ghToken == "" && userColor == "" && ottoColor == "" {
			log.Warn("No configuration options provided")
			os.Exit(0)
		}

		// if the api key is provided, set it
		if apiKey != "" {
			fmt.Println("Setting API key...")
			c.APIKey = apiKey
		}

		// if the model is provided, set it
		if model != "" {
			fmt.Println("Setting model...")
			if !utils.Contains(VALID_MODELS, model) {
				log.Errorf("Invalid model: %s", model)
				log.Errorf("Valid models are: %s", VALID_MODELS)
				os.Exit(1)
			}
			c.Model = model
		}

		// if the gh token is provided, set it
		if ghToken != "" {
			fmt.Println("Setting GitHub token...")
			c.GHToken = ghToken
		}

		// if the userColor is provided, set it
		if userColor != "" {
			fmt.Println("Setting user color...")
			c.UserColor = userColor
		}

		// if the ottoColor is provided, set it
		if ottoColor != "" {
			fmt.Println("Setting Otto color...")
			c.OttoColor = ottoColor
		}

		// save the config
		err = c.Save()
		if err != nil {
			log.Errorf("Error saving config: %s", err)
			os.Exit(1)
		}

		fmt.Println("Configuration saved successfully!")
	},
}

func init() {
	RootCmd.AddCommand(configCmd)

	// get api key
	configCmd.Flags().StringVarP(&apiKey, "apikey", "k", "", "API key to add to configuration")
	// get model
	configCmd.Flags().StringVarP(&model, "model", "m", "", "Model to use for documentation")
	// set gh token
	configCmd.Flags().StringVarP(&ghToken, "ghtoken", "t", "", "GitHub token to use for documentation")
	// set user color
	configCmd.Flags().StringVarP(&userColor, "userColor", "u", "", "User color for configuration")
	// set otto color
	configCmd.Flags().StringVarP(&ottoColor, "ottoColor", "o", "", "Otto color for configuration")
}