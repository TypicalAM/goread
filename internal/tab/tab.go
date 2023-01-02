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

// Tab is a general layout for a tab
type Tab interface {
	// General fields
	Title() string
	Loaded() bool
	Type() Type

	// Bubbletea methods
	Init() tea.Cmd
	Update(msg tea.Msg) (Tab, tea.Cmd)
	View() string
}

// NewTab is used to signal to the main model that a
// tab should be created
func NewTab(title string, tabType Type) tea.Cmd {
	return func() tea.Msg {
		return NewTabMessage{
			Title: title,
			Type:  tabType,
		}
	}
}

// The new tab message is sent when we want to enqueue a new tab
type NewTabMessage struct {
	// The new tab title
	Title string
	// The new tab type
	Type Type
}
