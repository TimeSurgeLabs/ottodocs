/*
Copyright © 2024 TimeSurgeLabs <chandler@timesurgelabs.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var printVersion bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "otto",
	Short: "Document your code with ease",
	Long:  `Code documentation made easy using GPT.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// run the version command
		if printVersion {
			versionCmd.Run(cmd, args)
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// RootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().BoolVarP(&printVersion, "version", "V", false, "print version")
}
