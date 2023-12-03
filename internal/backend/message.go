package backend

import (
	"github.com/TypicalAM/goread/internal/ui/tab"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// ArticleItem is an item that contains article data.
type ArticleItem struct {
	list.Item
	ArtTitle        string
	Desc            string
	RawDesc         string
	MarkdownContent string
	FeedURL         string
}

// FilterValue fulfills the list.Item interface
func (a ArticleItem) FilterValue() string {
	return a.ArtTitle
}

// Title fulfills the list.DefaultItem interface
func (a ArticleItem) Title() string {
	return a.ArtTitle
}

// Description fulfills the list.DefaultItem interface
func (a ArticleItem) Description() string {
	return a.Desc
}

// Fetcher fetches the data, it is used by tabs to query data.
type Fetcher func(feedname string) tea.Cmd

// ArticleFetcher fetches the article data, it is used by tabs to query data.
type ArticleFetcher func(feedname string, refresh bool) tea.Cmd

// FetchSuccessMsg is sent on fetch success.
type FetchSuccessMsg struct{ Items []list.Item }

// FetchArticleSuccessMsg is sent on article fetch success.
type FetchArticleSuccessMsg struct {
	Items []list.Item
}

// FetchErrorMsg is sent on fetch error.
type FetchErrorMsg struct {
	Err         error
	Description string
}

// NewItemMsg contains info the browser needs to know to add a new item.
type NewItemMsg struct{ Sender tab.Tab }

// NewItem is called from a tab to tell the browser that a new item needs to be added.
func NewItem(sender tab.Tab) tea.Cmd {
	return func() tea.Msg { return NewItemMsg{sender} }
}

// EditItemMsg contains info the browser needs to know to edit an item.
type EditItemMsg struct {
	Sender    tab.Tab
	OldFields []string
}

// EditItem is called from a tab to tell the browser that an item needs to be edited.
func EditItem(sender tab.Tab, fields []string) tea.Cmd {
	return func() tea.Msg { return EditItemMsg{sender, fields} }
}

// DeleteItemMsg contains info the browser needs to know to delete an item.
type DeleteItemMsg struct {
	Sender   tab.Tab
	ItemName string
}

// DeleteItem is called from a tab to tell the browser that an item needs to be deleted.
func DeleteItem(sender tab.Tab, itemName string) tea.Cmd {
	return func() tea.Msg { return DeleteItemMsg{sender, itemName} }
}

// DownloadItemMsg contains info the browser needs to know to download an item.
type DownloadItemMsg struct {
	FeedName string
	Index    int
}

// DownloadItem is called from a tab to tell the browser that an item needs to be downloaded.
func DownloadItem(feedName string, index int) tea.Cmd {
	return func() tea.Msg { return DownloadItemMsg{feedName, index} }
}

// MakeChoiceMsg contains info needed to create a binary choice prompt.
type MakeChoiceMsg struct {
	Question string
	Default  bool
}

// MakeChoice is called from a tab to tell the browser that a binary choice prompt needs to be created.
func MakeChoice(question string, defaultChoice bool) tea.Cmd {
	return func() tea.Msg { return MakeChoiceMsg{question, defaultChoice} }
}

// MarkAsReadMsg contains info needed to mark an item as read.
type MarkAsReadMsg string

// MarkAsRead is called from a tab to tell the browser that an item needs to be marked as read.
func MarkAsRead(url string) tea.Cmd {
	return func() tea.Msg { return MarkAsReadMsg(url) }
}

// MarkAsUnreadMsg contains info needed to mark an item as unread.
type MarkAsUnreadMsg string

// MarkAsUnread is called from a tab to tell the browser that an item needs to be marked as unread.
func MarkAsUnread(url string) tea.Cmd {
	return func() tea.Msg { return MarkAsUnreadMsg(url) }
}

// SetEnableKeybindMsg contains the desired state of the keybinds.
type SetEnableKeybindMsg bool

// SetEnableKeybind is called from a tab to tell the broswer that keybinds should be enabled/disabled.
func SetEnableKeybind(enable bool) tea.Cmd {
	return func() tea.Msg { return SetEnableKeybindMsg(enable) }
}

// StartQuittingMsg prompts the browser to start quitting (and perform a last browser redraw).
type StartQuittingMsg struct{}

// StartQuitting is called from a tab to tell the browser to start quitting.
func StartQuitting() tea.Cmd { return func() tea.Msg { return StartQuittingMsg{} } }
