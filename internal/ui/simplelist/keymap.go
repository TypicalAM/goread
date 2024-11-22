package simplelist

import "github.com/charmbracelet/bubbles/key"

// Keymap is the Keymap for the list
type Keymap struct {
	Open        key.Binding
	Up          key.Binding
	Down        key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	QuickSelect key.Binding
}

// DefaultKeymap is the default keymap for the list
var DefaultKeymap = Keymap{
	Open: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Open"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Move down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("shift+up", "K"),
		key.WithHelp("shift+↑/K", "Page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("shift+down", "J"),
		key.WithHelp("shift+↓/J", "Page down"),
	),
	QuickSelect: key.NewBinding(
		key.WithKeys("0", "1", "2", "3", "4", "5", "6", "7", "8", "9"),
		key.WithHelp("0-9", "Quick select"),
	),
}
