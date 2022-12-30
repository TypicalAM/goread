package list

type ListItem struct {
	title       string
	description string
	moreContent string
}

func (m ListItem) FilterValue() string {
	return m.title
}

func (m ListItem) Title() string {
	return m.title
}

func (m ListItem) Description() string {
	return m.description
}

// NewListItem creates a new list item
func NewListItem(title string, description string, moreContent string) ListItem {
	return ListItem{title, description, moreContent}
}

// GetContent returns the content
func (m ListItem) GetContent() string {
	return m.moreContent
}
