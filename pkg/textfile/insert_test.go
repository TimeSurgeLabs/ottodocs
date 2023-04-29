package textfile

import (
	"testing"
)

func TestReplaceLines(t *testing.T) {
	testCases := []struct {
		name           string
		code           string
		startLine      int
		endLine        int
		newText        string
		expectedResult string
		expectError    bool
	}{
		{
			name:           "Replace single line",
			code:           "Line 1\nLine 2\nLine 3",
			startLine:      2,
			endLine:        2,
			newText:        "New Line 2",
			expectedResult: "Line 1\nNew Line 2\nLine 3",
			expectError:    false,
		},
		{
			name:           "Replace multiple lines",
			code:           "Line 1\nLine 2\nLine 3",
			startLine:      1,
			endLine:        2,
			newText:        "New Line 1 and 2",
			expectedResult: "New Line 1 and 2\nLine 3",
			expectError:    false,
		},
		{
			name:           "Invalid startLine and endLine",
			code:           "Line 1\nLine 2\nLine 3",
			startLine:      3,
			endLine:        1,
			newText:        "Invalid Test",
			expectedResult: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ReplaceLines(tc.code, tc.startLine, tc.endLine, tc.newText)

			if tc.expectError && err == nil {
				t.Errorf("Expected error, but got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tc.expectedResult {
				t.Errorf("Expected result: %q, but got: %q", tc.expectedResult, result)
			}
		})
	}
}

func TestInsertLine(t *testing.T) {
	testCases := []struct {
		name           string
		code           string
		lineNumber     int
		newText        string
		expectedResult string
		expectError    bool
	}{
		{
			name:           "Insert at beginning",
			code:           "Line 1\nLine 2\nLine 3",
			lineNumber:     1,
			newText:        "New Line 0",
			expectedResult: "New Line 0\nLine 1\nLine 2\nLine 3",
			expectError:    false,
		},
		{
			name:           "Insert in the middle",
			code:           "Line 1\nLine 2\nLine 3",
			lineNumber:     2,
			newText:        "New Line 2",
			expectedResult: "Line 1\nNew Line 2\nLine 2\nLine 3",
			expectError:    false,
		},
		{
			name:           "Invalid line number",
			code:           "Line 1\nLine 2\nLine 3",
			lineNumber:     0,
			newText:        "Invalid Test",
			expectedResult: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := InsertLine(tc.code, tc.lineNumber, tc.newText)

			if tc.expectError && err == nil {
				t.Errorf("Expected error, but got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tc.expectedResult {
				t.Errorf("Expected result: %q, but got: %q", tc.expectedResult, result)
			}
		})
	}
}
