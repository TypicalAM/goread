package popup

import (
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// style is the style of the choice popup
type style struct {
	button       lipgloss.Style
	activeButton lipgloss.Style
	question     lipgloss.Style
	general      lipgloss.Style
}

// newStyle creates a new style for the choice popup
func newStyle(colors *theme.Colorscheme, width, height int) style {
	buttonStyle := lipgloss.NewStyle().
		Foreground(colors.TextDark).
		Background(colors.BgDark).
		Padding(0, 2).
		Margin(0, 1)

	activeButtonStyle := buttonStyle.Copy().
		Foreground(colors.Text).
		Background(colors.Color3).
		Underline(true)

	general := lipgloss.NewStyle().
		Foreground(colors.Text).
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(colors.Color1)

	question := lipgloss.NewStyle().
		Width(width).
		Margin(1, 0).
		Italic(true).
		Align(lipgloss.Center)

	return style{
		button:       buttonStyle,
		activeButton: activeButtonStyle,
		question:     question,
		general:      general,
	}
}
