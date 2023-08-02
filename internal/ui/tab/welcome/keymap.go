package welcome

import "github.com/charmbracelet/bubbles/key"

// Keymap contains the key bindings for this tab
type Keymap struct {
	NewCategory    key.Binding
	EditCategory   key.Binding
	DeleteCategory key.Binding
}

// DefaultKeymap contains the default key bindings for this tab
var DefaultKeymap = Keymap{
	NewCategory: key.NewBinding(
		key.WithKeys("n", "ctrl+n"),
		key.WithHelp("n/ctrl+n", "New"),
	),
	EditCategory: key.NewBinding(
		key.WithKeys("e", "ctrl+e"),
		key.WithHelp("e/ctrl+e", "Edit"),
	),
	DeleteCategory: key.NewBinding(
		key.WithKeys("d", "ctrl+d"),
		key.WithHelp("d/ctrl+d", "Delete"),
	),
}

// SetEnabled allows to disable/enable shortcuts
func (m *Keymap) SetEnabled(enabled bool) {
	m.NewCategory.SetEnabled(enabled)
	m.EditCategory.SetEnabled(enabled)
	m.DeleteCategory.SetEnabled(enabled)
}
