package lollypops

import (
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/TypicalAM/goread/internal/ui/popup"
	"github.com/charmbracelet/lipgloss"
)

// choiceStyle is the style of the choice popup
type choiceStyle struct {
	border       popup.TitleBorder
	button       lipgloss.Style
	activeButton lipgloss.Style
	question     lipgloss.Style
}

// newChoiceStyle creates a new style for the choice popup
func newChoiceStyle(colors *theme.Colors, width, height int) choiceStyle {
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

	return choiceStyle{
		border:       popup.NewTitleBorder("Confirm choice", width, height, colors.Color1, lipgloss.NormalBorder()),
		button:       buttonStyle,
		activeButton: activeButtonStyle,
		question:     question,
	}
}
