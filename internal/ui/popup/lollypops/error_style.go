package lollypops

import (
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/TypicalAM/goread/internal/ui/popup"
	"github.com/charmbracelet/lipgloss"
)

// errorStyle is the style of the error popup
type errorStyle struct {
	border       popup.TitleBorder
	activeButton lipgloss.Style
	msg          lipgloss.Style
}

// newErrorStyle creates a new style for the error popup
func newErrorStyle(colors *theme.Colors, width, height int) errorStyle {
	errorColor := lipgloss.Color("#f08ca8")
	buttonStyle := lipgloss.NewStyle().
		Foreground(colors.TextDark).
		Background(colors.BgDark).
		Padding(0, 2).
		Margin(0, 1)

	activeButtonStyle := buttonStyle.Copy().
		Foreground(colors.Text).
		Background(colors.Color3)

	msg := lipgloss.NewStyle().
		Width(width).
		Margin(1, 0).
		Align(lipgloss.Center)

	return errorStyle{
		border:       popup.NewTitleBorder("Error", width, height, errorColor, lipgloss.NormalBorder()),
		activeButton: activeButtonStyle,
		msg:          msg,
	}
}
