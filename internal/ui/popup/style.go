package popup

import (
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// style is the style of the choice popup
type style struct {
	border       TitleBorder
	button       lipgloss.Style
	activeButton lipgloss.Style
	question     lipgloss.Style
}

// newStyle creates a new style for the choice popup
func newStyle(colors *theme.Colors, width, height int) style {
	buttonStyle := lipgloss.NewStyle().
		Foreground(colors.TextDark).
		Background(colors.BgDark).
		Padding(0, 2).
		Margin(0, 1)

	activeButtonStyle := buttonStyle.Copy().
		Foreground(colors.Text).
		Background(colors.Color3)

	question := lipgloss.NewStyle().
		Width(width).
		Margin(1, 0).
		Italic(true).
		Align(lipgloss.Center)

	return style{
		border:       NewTitleBorder("Confirm choice", width, height, colors.Color1, lipgloss.NormalBorder()),
		button:       buttonStyle,
		activeButton: activeButtonStyle,
		question:     question,
	}
}
