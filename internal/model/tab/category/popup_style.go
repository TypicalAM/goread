package category

import (
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/lipgloss"
)

// popupStyle is the style of the popup window.
type popupStyle struct {
	general        lipgloss.Style
	heading        lipgloss.Style
	choiceSection  lipgloss.Style
	choice         lipgloss.Style
	selectedChoice lipgloss.Style
}

// newPopupStyle creates a new popup style.
func newPopupStyle(colors colorscheme.Colorscheme, width, height int) popupStyle {
	heading := lipgloss.NewStyle().
		Margin(2, 2).
		Width(width - 2).
		Align(lipgloss.Center).
		Italic(true)

	choiceSection := lipgloss.NewStyle().
		Padding(2).
		Width(width - 2).
		Height(10)

	choice := lipgloss.NewStyle().
		PaddingLeft(2).
		MarginBottom(1).
		Border(lipgloss.RoundedBorder(), false, false, false, true)

	selectedChoice := choice.Copy().
		BorderForeground(colors.Color4)

	general := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(colors.Color1)

	return popupStyle{
		heading:        heading,
		choiceSection:  choiceSection,
		choice:         choice,
		selectedChoice: selectedChoice,
		general:        general,
	}
}
