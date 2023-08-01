package feed

import "github.com/charmbracelet/bubbles/key"

// Keymap contains the key bindings for this tab
type Keymap struct {
	Open            key.Binding
	ToggleFocus     key.Binding
	RefreshArticles key.Binding
	SaveArticle     key.Binding
	DeleteFromSaved key.Binding
	CycleSelection  key.Binding
	MarkAsUnread    key.Binding
}

// DefaultKeymap contains the default key bindings for this tab
var DefaultKeymap = Keymap{
	Open: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Open"),
	),
	ToggleFocus: key.NewBinding(
		key.WithKeys("left", "right", "h", "l"),
		key.WithHelp("←/→", "Move left/right"),
	),
	RefreshArticles: key.NewBinding(
		key.WithKeys("r", "ctrl+r"),
		key.WithHelp("r/ctrl+r", "Refresh"),
	),
	SaveArticle: key.NewBinding(
		key.WithKeys("s", "ctrl+s"),
		key.WithHelp("s/ctrl+s", "Save"),
	),
	DeleteFromSaved: key.NewBinding(
		key.WithKeys("d", "ctrl+d"),
		key.WithHelp("d/ctrl+d", "Delete from saved"),
	),
	CycleSelection: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "Cycle selection"),
	),
	MarkAsUnread: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "Mark as unread"),
	),
}

// SetEnabled allows to disable/enable shortcuts
func (m *Keymap) SetEnabled(enabled bool) {
	m.Open.SetEnabled(enabled)
	m.ToggleFocus.SetEnabled(enabled)
	m.RefreshArticles.SetEnabled(enabled)
	m.SaveArticle.SetEnabled(enabled)
	m.DeleteFromSaved.SetEnabled(enabled)
	m.CycleSelection.SetEnabled(enabled)
	m.MarkAsUnread.SetEnabled(enabled)
}


