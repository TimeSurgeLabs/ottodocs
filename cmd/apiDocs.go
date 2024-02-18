/*
Copyright © 2024 TimeSurgeLabs <chandler@timesurgelabs.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	l "github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/TimeSurgeLabs/ottodocs/pkg/ai"
	"github.com/TimeSurgeLabs/ottodocs/pkg/config"
	"github.com/TimeSurgeLabs/ottodocs/pkg/git"
	"github.com/TimeSurgeLabs/ottodocs/pkg/utils"
)

// apiDocsCmd represents the apiDocs command
var apiDocsCmd = &cobra.Command{
	Use:   "apiDocs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: run,
	PreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetLevel(l.DebugLevel)
		}
	},
}

func init() {
	RootCmd.AddCommand(apiDocsCmd)

	apiDocsCmd.Flags().BoolVarP(&overwriteOriginal, "overwrite", "w", false, "Overwrite the original file.")
	apiDocsCmd.Flags().BoolVarP(&appendFile, "append", "a", false, "Append to the original file if the file exists.")
	apiDocsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging.")
	apiDocsCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Path to the output file.")
	// apiDocsCmd.Flags().StringSliceVarP(&routerFiles, "routerFiles", "r", []string{}, "Files that contain router information.")
	// apiDocsCmd.Flags().StringSliceVarP(&contextFiles, "contextFiles", "c", []string{}, "Files that contain context information.")
}

func run(cmd *cobra.Command, args []string) {
	var repoPath string
	if len(args) > 0 {
		repoPath = args[0]
	} else {
		repoPath = "."
	}

	conf, err := config.Load()
	if err != nil || conf.APIKey == "" {
		// if the API key is not set, prompt the user to config
		log.Error("Please config first.")
		log.Error("Run `ottodocs config -h` to learn how to config.")
		os.Exit(1)
	}

	info, err := os.Stat(repoPath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if outputFile == "" {
		outputFile = "api.md"
	}

	if !info.IsDir() {
		log.Error("Error: path is not a directory")
		os.Exit(1)
	}

	if !git.IsGitRepo(repoPath) {
		log.Error("Error: not a git repository")
		os.Exit(1)
	}

	fmt.Println("Getting repo...")
	repo, err := git.GetRepo(repoPath, ignoreFilePath, ignoreGitignore)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// for now we'll try this with ChatGPT 4 turbo preview as
	// it has a massive context limit
	conf.Model = "gpt-4-turbo-preview"
	if conf.BaseURL != "" {
		log.Warn("Using custom models is not supported for this command. The OpenAPI API will be used.")
		conf.BaseURL = ""
	}

	var files []string
	for _, file := range repo.Files {
		path := filepath.Join(repoPath, file.Path)
		contents, err := utils.LoadFile(path)
		if err != nil {
			log.Warnf("Error loading file %s: %s", path, err)
			continue
		}
		contents = `# ` + file.Path + "\n\n" + contents + "\n\n---\n\n"
		files = append(files, contents)
	}

	fmt.Println("Documenting repo...")
	content, err := ai.APIDocs(files, conf)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	exists := false
	// check if the output file exists
	if _, err := os.Stat(outputFile); err == nil {
		if !overwriteOriginal {
			log.Errorf("Error: output file %s already exists!", outputFile)
			os.Exit(1)
		}
		exists = true
	}
	var file *os.File
	if !exists {
		// write the string to the output file
		file, err = os.Create(outputFile)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	} else {
		// append if the file already exists
		file, err = os.OpenFile(outputFile, os.O_WRONLY, 0644)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	_, err = file.WriteString(content)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	fmt.Printf("API documentation written to %s\n", outputFile)
}
