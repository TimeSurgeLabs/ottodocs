package cmd

import (
	"os"

	l "github.com/charmbracelet/log"
)

var force bool
var verbose bool
var question string

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
var noCommit bool
var push bool

var base string
var title string

var model string
var apiKey string
var ghToken string
var remote string
var userColor string
var ottoColor string

var issuePRNumber int
var useComments bool
var promptOnly bool
var countFinalTokens bool

var startLine int
var endLine int
var appendFile bool

var previousTag string
var currentTag string

var contextFiles []string

var displayHistory bool
var loadHistory string
var deleteHistory string
var readOnly bool
var clearHistory bool
var repoContext bool
var organization string

var log = l.NewWithOptions(os.Stderr, l.Options{
	Level:           l.InfoLevel,
	ReportTimestamp: false,
})
