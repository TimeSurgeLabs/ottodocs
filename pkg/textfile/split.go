package textfile

import (
	"fmt"
	"strings"

	"github.com/chand1012/ottodocs/pkg/calc"
	"github.com/chand1012/ottodocs/pkg/constants"
	"github.com/chand1012/ottodocs/pkg/utils"
)

type SplitFile struct {
	Path      string
	Contents  string
	StartLine int
	EndLine   int
}

func (s *SplitFile) Hash() string {
	// path/to/file.go#5-12
	return fmt.Sprintf("%s#%d-%d", s.Path, s.StartLine, s.EndLine)
}

func Split(path string) ([]SplitFile, error) {
	// read the file with the os package
	// split the file into lines with a max of 4000 tokens
	// return a slice of SplitFile structs

	contents, err := utils.LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not load file: %s", err)
	}

	lines := strings.Split(contents, "\n")

	var splitFiles []SplitFile

	for _, line := range lines {
		if len(splitFiles) == 0 {
			splitFiles = append(splitFiles, SplitFile{
				Path:      path,
				Contents:  line,
				StartLine: 1,
				EndLine:   1,
			})
			continue
		}

		lastSplitFile := splitFiles[len(splitFiles)-1]

		if calc.EstimateTokens(lastSplitFile.Contents)+calc.EstimateTokens(line) > constants.OPENAI_MAX_TOKENS {
			splitFiles = append(splitFiles, SplitFile{
				Path:      path,
				Contents:  line,
				StartLine: lastSplitFile.EndLine + 1,
				EndLine:   lastSplitFile.EndLine + 1,
			})
			continue
		}

		lastSplitFile.Contents += "\n" + line
		lastSplitFile.EndLine++
		splitFiles[len(splitFiles)-1] = lastSplitFile
	}

	return splitFiles, nil
}
