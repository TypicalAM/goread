package list

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/reflow/wrap"
)

type Item struct {
	title       string
	description string
	moreContent string
}

func (m Item) FilterValue() string {
	return m.title
}

func (m Item) Title() string {
	return m.title
}

func (m Item) Description() string {
	return m.description
}

// NewListItem creates a new list item
func NewListItem(title string, description string, moreContent string) Item {
	return Item{title, description, moreContent}
}

// StyleContent styles the content of the item with glamour
// and returns the result
func (m Item) StyleContent(width int) string {
	// Create a renderer for the content
	g, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return fmt.Sprintf("We have encountered an error styling the content: %s", err)
	}

	// Style the content
	styledContent, err := g.Render(m.moreContent)
	if err != nil {
		return fmt.Sprintf("We have encountered an error styling the content: %s", err)
	}

	// Return the styled content
	return styledContent
}

// Wrap description to a given width
func (m Item) WrapDescription(width int) Item {
	m.description = wrap.String(m.description, width-3)
	return m
}
