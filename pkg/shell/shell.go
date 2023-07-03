package shell

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/TimeSurgeLabs/ottodocs/pkg/utils"
)

const ZSH_HISTORY_PATH = ".zsh_history"
const BASH_HISTORY_PATH = ".bash_history"

// for now we only support bash and zsh for unix
// we will support more in the future if demand is there
// windows is on the roadmap

// GetShellHistory gets the most recently used command from the shell history
// file. It will attempt to open all of them, only getting the most recently
// modified one. Only get n lines. If the history file is zsh, it will just get
// the command and not the metadata.
func GetHistory(n int) ([]string, error) {
	shellHistories := []string{ZSH_HISTORY_PATH, BASH_HISTORY_PATH}
	var mostRecentHistory string
	var mostRecentTime int64
	var shellUsed string
	for _, history := range shellHistories {
		// get the home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		// get the history file path
		historyPath := filepath.Join(homeDir, history)
		// get the file info
		info, err := os.Stat(historyPath)
		if err != nil {
			continue
		}
		// check if the file is newer than the current most recent
		if info.ModTime().Unix() > mostRecentTime {
			mostRecentTime = info.ModTime().Unix()
			mostRecentHistory = historyPath
			shellUsed = history
		}
	}

	// if we didn't find a history file, return an error
	if mostRecentHistory == "" {
		return nil, os.ErrNotExist
	}

	// open the file
	contents, err := os.ReadFile(mostRecentHistory)
	if err != nil {
		return nil, err
	}

	// fmt.Println("shell used: ", shellUsed)
	// fmt.Println("contents: ", string(contents))

	var finalLines []string
	if shellUsed == ZSH_HISTORY_PATH {
		// split the file by newlines
		lines := strings.Split(string(contents), "\n")
		// reverse the lines
		utils.ReverseSlice(lines)
		for _, line := range lines {
			// get just the command. Split each line after the first semicolon
			// and get the first element
			split_command := strings.SplitN(line, ";", 2)
			if len(split_command) < 2 {
				continue
			}

			command := split_command[1]
			// if the command is empty, skip it
			if command == "" {
				continue
			}
			// add the command to the final lines
			finalLines = append(finalLines, command)
			if len(finalLines) >= n {
				break
			}
		}
	}

	if shellUsed == BASH_HISTORY_PATH {
		// just get the last 100 lines
		finalLines = strings.Split(string(contents), "\n")
		if len(finalLines) > n {
			finalLines = finalLines[len(finalLines)-100:]
		}
	}

	return finalLines, nil
}
