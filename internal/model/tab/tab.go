package tab

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Type int

const (
	Welcome Type = iota + 1
	Feed
	Category
)

// Tab is an interface outlining the methods that a tab should implement including bubbletea's model methods
type Tab interface {
	// general fields
	Title() string
	Type() Type
	SetSize(width, height int) Tab
	GetKeyBinds() []key.Binding

	// bubbletea methods
	Init() tea.Cmd
	Update(msg tea.Msg) (Tab, tea.Cmd)
	View() string
}

// NewTab returns a tea.Cmd which sends a message to the main
// model to create a new tab
func NewTab(title string, tabType Type) tea.Cmd {
	return func() tea.Msg {
		return NewTabMessage{
			Title: title,
			Type:  tabType,
		}
	}
}

// NewTabMessage is a tea.Msg that signals that a new tab should be created
type NewTabMessage struct {
	Title string
	Type  Type
}
