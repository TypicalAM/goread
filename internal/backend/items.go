package backend

import (
	"github.com/TypicalAM/goread/internal/ui/tab"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Fetcher is a function that fetches the items
type Fetcher func(string) tea.Cmd

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
	Err         error
	Description string
}

// NewItemMsg is a message to tell the main model that a new item
// needs to be added to the list
type NewItemMsg struct {
	Sender    tab.Tab
	OldFields []string
	Editing   bool
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

// MakeChoiceMsg is a message to tell the main model that a choice
// needs to be made
type MakeChoiceMsg struct {
	Question string
	Default  bool
}

// MakeChoice is a function to tell the main model that a choice
// needs to be made
func MakeChoice(question string, defaultChoice bool) tea.Cmd {
	return func() tea.Msg {
		return MakeChoiceMsg{
			Question: question,
			Default:  defaultChoice,
		}
	}
}

// MarkAsRead is a message to tell the main model that a new item needs to be marked as read
type MarkAsReadMsg struct {
	Key   string
	Index int
}

// MarkAsRead is a function to tell the main model that a new item needs to be marked as read
func MarkAsRead(key string, index int) tea.Cmd {
	return func() tea.Msg {
		return MarkAsReadMsg{
			Key:   key,
			Index: index,
		}
	}
}

// MarkAsUnread is a message to tell the main model that a new item needs to be marked as unread
type MarkAsUnreadMsg struct {
	Key   string
	Index int
}

// MarkAsUnread is a function to tell the main model that a new item needs to be marked as unread
func MarkAsUnread(key string, index int) tea.Cmd {
	return func() tea.Msg {
		return MarkAsUnreadMsg{
			Key:   key,
			Index: index,
		}
	}
}
