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

// newListStyle creates a new listStyle
func newListStyle() listStyle {
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
		Foreground(style.GlobalColorscheme.Color2)

	// Return the new style
	return newStyle
}

// styleDescription will style the description of the item
func (s listStyle) styleDescription(description string) string {
	// Create the arrow style
	arrowStyle := lipgloss.NewStyle().
		MarginLeft(10).
		Foreground(style.GlobalColorscheme.Color3)

	// Create the text style
	textStyle := lipgloss.NewStyle().
		MarginLeft(1).
		Foreground(style.GlobalColorscheme.Color3)

	return arrowStyle.Render("⮡") + textStyle.Render(description)
}
