package tab

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Type int

const (
	Welcome Type = iota
	Feed
	Category
)

// Tab is an interface outlining the methods that a tab should implement
// a bubbletea models' methods and also some more
type Tab interface {
	// general fields
	Title() string
	Type() Type
	Help() Help

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

// Help is a struct containing the keys and their descriptions
// for a given tab
type Help []KeyBind

// KeyBind is a struct containing the description of a tab
// and the keys that are used to interact with it
type KeyBind struct {
	Key         string
	Description string
}
