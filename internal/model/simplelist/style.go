package simplelist

import (
	"github.com/TypicalAM/goread/internal/style"
	"github.com/charmbracelet/lipgloss"
)

// listStyle is the style of the list.
type listStyle struct {
	titleStyle   lipgloss.Style
	noItemsStyle lipgloss.Style
	itemStyle    lipgloss.Style
}

// newDefaultListStyle returns a new list style with default values.
func newDefaultListStyle() listStyle {
	// Create the new style
	newStyle := listStyle{}

	// titleStyle is used to style the title of the list
	newStyle.titleStyle = lipgloss.NewStyle().
		Foreground(style.GlobalColorscheme.Color1).
		MarginLeft(3).
		PaddingBottom(1)

	// newItemsStyle is used to style the message when there are no items
	newStyle.noItemsStyle = lipgloss.NewStyle().
		MarginLeft(3).
		Foreground(style.GlobalColorscheme.Color2).
		Italic(true)

	// newItemsStyle is used to style the items in the list
	newStyle.itemStyle = lipgloss.NewStyle().
		MarginLeft(3).
		MarginRight(3).
		Foreground(style.GlobalColorscheme.Color2)

	// Return the new style
	return newStyle
}
