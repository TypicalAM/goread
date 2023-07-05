package feed

import (
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/lipgloss"
)

// popupStyle is the style of the popup window.
type popupStyle struct {
	general   lipgloss.Style
	heading   lipgloss.Style
	item      lipgloss.Style
	itemTitle lipgloss.Style
	itemField lipgloss.Style
}

// newPopupStyle creates a new popup style.
func newPopupStyle(colors colorscheme.Colorscheme, width, height int) popupStyle {
	general := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(colors.Color1)

	heading := lipgloss.NewStyle().
		Margin(1, 0, 1, 0).
		Width(width - 2).
		Align(lipgloss.Center).
		Italic(true)

	item := lipgloss.NewStyle().
		Margin(0, 4).
		PaddingLeft(1).
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(colors.Color3).
		Italic(true)

	itemTitle := lipgloss.NewStyle().
		Foreground(colors.Color3)

	itemField := lipgloss.NewStyle().
		Foreground(colors.Color2)

	return popupStyle{
		general:   general,
		heading:   heading,
		item:      item,
		itemTitle: itemTitle,
		itemField: itemField,
	}
}
