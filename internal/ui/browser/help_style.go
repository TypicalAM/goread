package browser

import (
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

// helpStyle is the style for the help popup.
type helpStyle struct {
	help  help.Styles
	title lipgloss.Style
	box   lipgloss.Style
}

// newHelpStyle creates a new style for the help popup title and the help model.
func newHelpStyle(colors *theme.Colors, width, height int) helpStyle {
	styles := help.Styles{}
	styles.FullDesc = lipgloss.NewStyle().
		Foreground(colors.Text)
	styles.FullKey = lipgloss.NewStyle().
		Foreground(colors.Color2)
	styles.FullSeparator = lipgloss.NewStyle().
		Foreground(colors.TextDark)

	return helpStyle{
		help: styles,
		title: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Margin(1, 0).
			Width(width - 2).
			Foreground(colors.Text).
			Italic(true),
		box: lipgloss.NewStyle().
			Width(width - 2).
			Height(height - 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(colors.Color1),
	}
}
