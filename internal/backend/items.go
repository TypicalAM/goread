package backend

import (
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// FetchSuccessMsg is a message that is sent when the fetching of the
// categories or feeds was successful
type FetchSuccessMsg struct {
	Items []list.Item
}

// FetchArticleSuccessMsg is a message that is sent when the fetching of the
// articles was successful
type FetchArticleSuccessMsg struct {
	Items           []list.Item
	ArticleContents []string
}

// FetchErrorMsgssage that is sent when the fetching of the
// categories or feeds failed
type FetchErrorMsg struct {
	Description string
	Err         error
}

// NewItemMsg is a message to tell the main model that a new item
// needs to be added to the list
type NewItemMsg struct {
	Sender    tab.Tab
	Editing   bool
	OldFields []string
}

// NewItem is a function to tell the main model that a new item
// needs to be added to the list
func NewItem(sender tab.Tab, editing bool, fields []string) tea.Cmd {
	return func() tea.Msg {
		return NewItemMsg{
			Sender:    sender,
			Editing:   editing,
			OldFields: fields,
		}
	}
}

// DeleteItemMsg is a message to tell the main model that a new item
// needs to be removed from the list
type DeleteItemMsg struct {
	Sender tab.Tab
	Key    string
}

// DeleteItem is a function to tell the main model that a new item
// needs to be removed from the list
func DeleteItem(sender tab.Tab, key string) tea.Cmd {
	return func() tea.Msg {
		return DeleteItemMsg{
			Sender: sender,
			Key:    key,
		}
	}
}

// DownloadItemMsg is a message to tell the main model that a new item
// needs to be downloaded
type DownloadItemMsg struct {
	Key   string
	Index int
}

// DownloadItem is a function to tell the main model that a new item
// needs to be downloaded
func DownloadItem(key string, index int) tea.Cmd {
	return func() tea.Msg {
		return DownloadItemMsg{
			Key:   key,
			Index: index,
		}
	}
}
