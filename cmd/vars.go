package cmd

import (
	"os"

	l "github.com/charmbracelet/log"
)

var verbose bool

var repoPath string
var preambleFile string
var outputFile string
var estimateTokens bool
var ignoreFilePath string
var ignoreGitignore bool
var outputJSON bool

var filePath string
var chatPrompt string
var inlineMode bool
var markdownMode bool
var overwriteOriginal bool

var conventional bool // use conventional commits
var plain bool
var auto bool
var push bool

var base string
var title string

var model string
var apiKey string
var ghToken string
var remote string

var log = l.NewWithOptions(os.Stderr, l.Options{
	Level:           l.InfoLevel,
	ReportTimestamp: false,
})
