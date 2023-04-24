/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/chand1012/ottodocs/pkg/config"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Add an API key to your configuration",
	Long: `Add an API key to your configuration. 
This API key will be used to authenticate with the OpenAI ChatGPT API.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if the api key is not provided, prompt the user for it
		if apiKey == "" {
			fmt.Print("Please provide an API key: ")
			fmt.Scanln(&apiKey)
		}
		// save the API key to a configuration file at ~/.ottodocs/config.json
		conf, err := config.Load()
		if err != nil {
			log.Errorf("Error: %s", err)
			os.Exit(1)
		}
		conf.APIKey = apiKey
		err = conf.Save()
		if err != nil {
			log.Errorf("Error: %s", err)
			os.Exit(1)
		}
		fmt.Println("API key saved successfully!")
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&apiKey, "apikey", "k", "", "API key to add to configuration")
}
