package style

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	ColumnStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(GlobalColorscheme.TextDark)

	FocusedStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(GlobalColorscheme.Color1)

	tabBorder = lipgloss.Border{Left: "┃"}

	ActiveTab = lipgloss.NewStyle().
			Padding(0, 7, 0, 1).
			Italic(true).
			Bold(true)

	ActiveTabIcon = lipgloss.NewStyle().
			Padding(0, 0, 0, 3).
			Bold(true).
			Border(tabBorder, false, false, false, true).
			BorderForeground(GlobalColorscheme.TextDark)

	TabStyle = lipgloss.NewStyle().
			Padding(0, 7, 0, 1).
			Background(GlobalColorscheme.BgDark).
			Foreground(GlobalColorscheme.TextDark)

	TabIcon = ActiveTabIcon.Copy().
		Background(GlobalColorscheme.BgDark).
		BorderForeground(GlobalColorscheme.BgDarker).
		BorderBackground(GlobalColorscheme.BgDark)

	TabGap = lipgloss.NewStyle().
		Background(GlobalColorscheme.BgDarker)

	StatusBarGap = lipgloss.NewStyle().
			Background(GlobalColorscheme.BgDark)
)

// Utility function to output the bigger nubmer
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Style an index value
func Index(index int, isSelected bool) string {
	// Define the styles used in the index styling
	bracketStyle := lipgloss.NewStyle().
		Foreground(GlobalColorscheme.Color7)
	numberStyle := lipgloss.NewStyle().
		Foreground(GlobalColorscheme.Color6)

	// If the index is the active index render it differently
	if isSelected {
		numberStyle = numberStyle.Background(GlobalColorscheme.Text)
	}

	// Check if the index is a digit
	var indexString string
	if index < 10 {
		// Show a digit
		indexString = fmt.Sprintf("%d", index)
	} else {
		// Show a letter
		indexString = fmt.Sprintf("%c", index+87)
	}
	// Render the whole index
	return bracketStyle.Render("[") +
		numberStyle.Render(indexString) +
		bracketStyle.Render("]")
}
