package backend

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// ItemType is the type of the added item
type ItemType int

const (
	Category ItemType = iota + 1
	Feed
)

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
	Type      ItemType
	New       bool
	Fields    []string
	ItemPath  []string
	OldFields []string
}

// NewItem is a function to tell the main model that a new item
// needs to be added to the list
func NewItem(itemType ItemType, newItem bool, itemPath []string, oldFields []string) tea.Cmd {
	return func() tea.Msg {
		var textFields []string
		if itemType == Category {
			textFields = []string{"Name", "Description"}
		} else {
			textFields = []string{"Name", "URL"}
		}

		return NewItemMessage{
			Type:      itemType,
			New:       newItem,
			Fields:    textFields,
			ItemPath:  itemPath,
			OldFields: oldFields,
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

// DownloadItemMessage is a message to tell the main model that a new item
// needs to be downloaded
type DownloadItemMessage struct {
	Key   string
	Index int
}

// DownloadItem is a function to tell the main model that a new item
// needs to be downloaded
func DownloadItem(key string, index int) tea.Cmd {
	return func() tea.Msg {
		return DownloadItemMessage{
			Key:   key,
			Index: index,
		}
	}
}
