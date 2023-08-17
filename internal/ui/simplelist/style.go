package simplelist

import (
	"fmt"
	"strings"

	"github.com/TypicalAM/goread/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// listStyle is the style of the list.
type listStyle struct {
	colors       *theme.Colors
	titleStyle   lipgloss.Style
	noItemsStyle lipgloss.Style
	itemStyle    lipgloss.Style

	bracketStyle lipgloss.Style
	numberStyle  lipgloss.Style
}

// newListStyle creates a new listStyle
func newListStyle(colors *theme.Colors) listStyle {
	titleStyle := lipgloss.NewStyle().
		Foreground(colors.Color1).
		MarginLeft(3).
		PaddingBottom(1)

	noItemsStyle := lipgloss.NewStyle().
		MarginLeft(3).
		Foreground(colors.Color2).
		Italic(true)

	itemStyle := lipgloss.NewStyle().
		MarginLeft(3).
		Foreground(colors.Color2)

	bracketStyle := lipgloss.NewStyle().
		Foreground(colors.Color7)

	numberStyle := lipgloss.NewStyle().
		Foreground(colors.Color6)

	return listStyle{
		colors:       colors,
		titleStyle:   titleStyle,
		noItemsStyle: noItemsStyle,
		itemStyle:    itemStyle,
		bracketStyle: bracketStyle,
		numberStyle:  numberStyle,
	}
}

// styleDescription will style the description of the item
func (s listStyle) styleDescription(description string) string {
	// Create the arrow style
	arrowStyle := lipgloss.NewStyle().
		MarginLeft(10).
		Foreground(s.colors.Color3)

	// Create the text style
	textStyle := lipgloss.NewStyle().
		MarginLeft(1).
		Foreground(s.colors.Color3)

	return arrowStyle.Render("⮡") + textStyle.Render(description)
}

// styleIndex will style the index of the item
func (s listStyle) styleIndex(index int, isSelected bool) string {
	var b strings.Builder
	b.WriteString("   ")
	b.WriteString(s.bracketStyle.Render("["))

	// If the index is the active index render it differently
	style := s.numberStyle.Copy()
	if isSelected {
		style = style.Background(s.colors.Text)
	}

	// Check if the index is a digit
	switch {
	case index < 10:
		b.WriteString(style.Render(fmt.Sprintf("%d", index)))

	case index < 36:
		b.WriteString(style.Render(fmt.Sprintf("%c", index+87)))

	default:
		b.WriteString(style.Render("-"))
	}

	// Render the whole index
	b.WriteString(s.bracketStyle.Render("]"))
	return b.String()
}
