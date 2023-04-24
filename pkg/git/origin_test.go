package git

import (
	"testing"
)

func TestExtractOriginInfo(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		gitURL        string
		expectedOwner string
		expectedRepo  string
		expectedErr   error
	}{
		{
			name:          "SSH URL",
			gitURL:        "git@github.com:owner/repo.git",
			expectedOwner: "owner",
			expectedRepo:  "repo",
			expectedErr:   nil,
		},
		{
			name:          "HTTPS URL",
			gitURL:        "https://github.com/owner/repo.git",
			expectedOwner: "owner",
			expectedRepo:  "repo",
			expectedErr:   nil,
		},
		{
			name:          "Short SSH URL",
			gitURL:        "git@github.com:owner/repo",
			expectedOwner: "owner",
			expectedRepo:  "repo",
			expectedErr:   nil,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			owner, repo, err := ExtractOriginInfo(test.gitURL)

			if owner != test.expectedOwner {
				t.Errorf("Expected owner '%s', but got '%s'", test.expectedOwner, owner)
			}

			if repo != test.expectedRepo {
				t.Errorf("Expected repo '%s', but got '%s'", test.expectedRepo, repo)
			}

			if err != test.expectedErr {
				t.Errorf("Expected error '%v', but got '%v'", test.expectedErr, err)
			}
		})
	}
}
