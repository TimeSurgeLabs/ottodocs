/*
Copyright © 2023 Chandler <chandler@chand1012.dev>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Generate a pull request",
	Long:  `Generates a pull request from commit messages, title, and the git diff.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pr called")
	},
}

func init() {
	RootCmd.AddCommand(prCmd)

}
