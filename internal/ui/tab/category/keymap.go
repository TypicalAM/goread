package category

import "github.com/charmbracelet/bubbles/key"

// Keymap contains the key bindings for this tab
type Keymap struct {
	NewFeed    key.Binding
	EditFeed   key.Binding
	DeleteFeed key.Binding
}

// DefaultKeymap contains the default key bindings for this tab
var DefaultKeymap = Keymap{
	NewFeed: key.NewBinding(
		key.WithKeys("n", "ctrl+n"),
		key.WithHelp("n/ctrl+n", "New"),
	),
	EditFeed: key.NewBinding(
		key.WithKeys("e", "ctrl+e"),
		key.WithHelp("e/ctrl+e", "Edit"),
	),
	DeleteFeed: key.NewBinding(
		key.WithKeys("d", "ctrl+d"),
		key.WithHelp("d/ctrl+d", "Delete"),
	),
}

// SetEnabled allows to disable/enable shortcuts
func (m *Keymap) SetEnabled(enabled bool) {
	m.NewFeed.SetEnabled(enabled)
	m.EditFeed.SetEnabled(enabled)
	m.DeleteFeed.SetEnabled(enabled)
}
