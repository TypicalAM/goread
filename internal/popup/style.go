package popup

import (
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/lipgloss"
)

// choiceStyle is the style of the choice popup
type choiceStyle struct {
	button       lipgloss.Style
	activeButton lipgloss.Style
	question     lipgloss.Style
	general      lipgloss.Style
}

// newChoiceStyle creates a new choiceStyle
func newChoiceStyle(colors colorscheme.Colorscheme, width, height int) choiceStyle {
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

	return choiceStyle{
		button:       buttonStyle,
		activeButton: activeButtonStyle,
		question:     question,
		general:      general,
	}
}
