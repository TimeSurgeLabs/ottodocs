package main

import (
	"github.com/TimeSurgeLabs/ottodocs/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs")
	if err != nil {
		panic(err)
	}
}
