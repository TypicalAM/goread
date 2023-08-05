package tab

import tea "github.com/charmbracelet/bubbletea"

// NewTab returns a tea.Cmd which sends a message to the main model to create a new tab
func NewTab(sender Tab, title string) tea.Cmd {
	return func() tea.Msg {
		return NewTabMsg{
			Sender: sender,
			Title:  title,
		}
	}
}

// NewTabMsg is a tea.Msg that signals that a new tab should be created.
type NewTabMsg struct {
	Sender Tab
	Title  string
}
