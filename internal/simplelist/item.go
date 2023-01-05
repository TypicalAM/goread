package simplelist

import (
	"github.com/charmbracelet/glamour"
)

// Item is a single item in a list
type Item struct {
	title       string
	description string
	moreContent string
}

// NewItem creates a new list item
func NewItem(title string, description string, moreContent string) Item {
	return Item{title, description, moreContent}
}

// FilterValue returns the value which is used for filtering
func (m Item) FilterValue() string {
	return m.title
}

// Title returns the title of the item
func (m Item) Title() string {
	return m.title
}

// Description returns the description of the item
func (m Item) Description() string {
	return m.description
}

// StyleContent styles the content of the item with glamour and returns the result
func (m Item) StyleContent(width int) (string, error) {
	// Create a renderer for the content
	g, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}

	// Style the content
	styledContent, err := g.Render(m.moreContent)
	if err != nil {
		return "", err
	}

	// Return the styled content
	return styledContent, nil
}
