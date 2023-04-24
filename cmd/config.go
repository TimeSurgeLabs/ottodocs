/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"os"

	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configures ottodocs",
	Long:  `Configures ottodocs. Allows user to specify API Keys and the model with a single command.`,
	Run: func(cmd *cobra.Command, args []string) {
		// load the config
		c, err := config.Load()
		if err != nil {
			log.Errorf("Error loading config: %s", err)
			os.Exit(1)
		}

		// if none of the config options are provided, print a warning
		if apiKey == "" && model == "" && ghToken == "" {
			log.Warn("No configuration options provided")
			return
		}

		// if the api key is provided, set it
		if apiKey != "" {
			log.Info("Setting API key...")
			c.APIKey = apiKey
		}

		// if the model is provided, set it
		if model != "" {
			log.Info("Setting model...")
			c.Model = model
		}

		// if the gh token is provided, set it
		if ghToken != "" {
			log.Info("Setting GitHub token...")
			c.GHToken = ghToken
		}

		// save the config
		err = c.Save()
		if err != nil {
			log.Errorf("Error saving config: %s", err)
			os.Exit(1)
		}

		log.Info("Configuration saved successfully!")
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
}
