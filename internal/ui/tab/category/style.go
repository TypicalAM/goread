package category

import (
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/TypicalAM/goread/internal/ui/popup"
	"github.com/charmbracelet/lipgloss"
)

// popupStyle is the style of the popup window.
type popupStyle struct {
	border    popup.TitleBorder
	listItem  lipgloss.Style
	itemTitle lipgloss.Style
	itemField lipgloss.Style
}

// newPopupStyle creates a new popup style.
func newPopupStyle(colors *theme.Colors, width, height int, headingText string) popupStyle {
	border := popup.NewTitleBorder(headingText, width, height, colors.Color1, lipgloss.NormalBorder())

	item := lipgloss.NewStyle().
		Margin(1, 4).
		PaddingLeft(1).
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(colors.Color3).
		Italic(true)

	itemTitle := lipgloss.NewStyle().
		Foreground(colors.Color3)

	itemField := lipgloss.NewStyle().
		Foreground(colors.Color2)

	return popupStyle{
		border:    border,
		listItem:  item,
		itemTitle: itemTitle,
		itemField: itemField,
	}
}
