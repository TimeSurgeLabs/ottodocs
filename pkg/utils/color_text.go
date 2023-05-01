package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// PrintColoredText prints the given text in the specified CSS color code using fmt.Print.
func PrintColoredText(text, cssColorCode string) {
	// Create a new Lipgloss style with a foreground color
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(cssColorCode))

	// Style the input text
	styledText := style.Render(text)

	// Print the styled text to the terminal without a newline
	fmt.Print(styledText)
}

// PrintColoredTextLn prints the given text in the specified CSS color code using fmt.Println.
func PrintColoredTextLn(text, cssColorCode string) {
	// Create a new Lipgloss style with a foreground color
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(cssColorCode))

	// Style the input text
	styledText := style.Render(text)

	// Print the styled text to the terminal with a newline
	fmt.Println(styledText)
}
