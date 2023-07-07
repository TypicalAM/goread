package backend

import (
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// FetchSuccessMessage is a message that is sent when the fetching of the
// categories or feeds was successful
type FetchSuccessMessage struct {
	Items []list.Item
}

// FetchArticleSuccessMessage is a message that is sent when the fetching of the
// articles was successful
type FetchArticleSuccessMessage struct {
	Items        []list.Item
	Descriptions []string
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
	Sender    tab.Tab
	Editing   bool
	OldFields []string
}

// NewItem is a function to tell the main model that a new item
// needs to be added to the list
func NewItem(sender tab.Tab, editing bool, fields []string) tea.Cmd {
	return func() tea.Msg {
		return NewItemMessage{
			Sender:    sender,
			Editing:   editing,
			OldFields: fields,
		}
	}
}

// DeleteItemMessage is a message to tell the main model that a new item
// needs to be removed from the list
type DeleteItemMessage struct {
	Sender tab.Tab
	Key    string
}

// DeleteItem is a function to tell the main model that a new item
// needs to be removed from the list
func DeleteItem(sender tab.Tab, key string) tea.Cmd {
	return func() tea.Msg {
		return DeleteItemMessage{
			Sender: sender,
			Key:    key,
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
