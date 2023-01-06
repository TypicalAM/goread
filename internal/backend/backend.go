package backend

import (
	"github.com/TypicalAM/goread/internal/rss"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// ItemType is the type of the added item
type ItemType int

const (
	Category ItemType = iota
	Feed
)

// The Backend is the model for the backend of the program. It is responsible
// for managing the data in the categories and feeds
type Backend interface {
	// Name returns the name of the backend to show in the logs
	Name() string
	// FetchCategories returns a tea.Cmd which gets the category list
	// fron the backend
	FetchCategories() tea.Cmd
	// FetchFeeds returns a tea.Cmd which gets the feed list from
	// the backend via a string key
	FetchFeeds(catName string) tea.Cmd
	// FetchArticles returns a tea.Cmd which gets the articles from
	// the backend via a string key
	FetchArticles(feedName string) tea.Cmd
	// Rss returns the pointer to the RSS struct, it can be used
	// to add and remove items
	Rss() *rss.Rss
	// Close closes the backend
	Close() error
}

// FetchSuccessMessage is a message that is sent when the fetching of the
// categories or feeds was successful
type FetchSuccessMessage struct {
	Items []list.Item
}

// FetchErrorMessage is a message that is sent when the fetching of the
// categories or feeds failed
type FetchErrorMessage struct {
	Description string
	Err         error
}

// NewItemMessage is a message to tell the main model that a new item
// needs to be added to the list
type NewItemMessage struct {
	Type     ItemType
	New      bool
	Fields   []string
	ItemPath []string
}

// NewItem is a function to tell the main model that a new item
// needs to be added to the list
func NewItem(itemType ItemType, newItem bool, itemPath []string) tea.Cmd {
	return func() tea.Msg {
		var textFields []string
		if itemType == Category {
			textFields = []string{"Name", "Description"}
		} else {
			textFields = []string{"Name", "URL"}
		}

		return NewItemMessage{
			Type:     itemType,
			New:      newItem,
			Fields:   textFields,
			ItemPath: itemPath,
		}
	}
}

// DeleteItemMessage is a message to tell the main model that a new item
// needs to be removed from the list
type DeleteItemMessage struct {
	Type ItemType
	Key  string
}

// DeleteItem is a function to tell the main model that a new item
// needs to be removed from the list
func DeleteItem(itemType ItemType, key string) tea.Cmd {
	return func() tea.Msg {
		return DeleteItemMessage{
			Type: itemType,
			Key:  key,
		}
	}
}
