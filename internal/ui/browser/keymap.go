package browser

import "github.com/charmbracelet/bubbles/key"

// Keymap contains the key bindings for the browser
type Keymap struct {
	CloseTab          key.Binding
	NextTab           key.Binding
	PrevTab           key.Binding
	ShowHelp          key.Binding
	ToggleOfflineMode key.Binding
}

// DefaultKeymap contains the default key bindings for the browser
var DefaultKeymap = Keymap{
	CloseTab: key.NewBinding(
		key.WithKeys("c", "ctrl+w"),
		key.WithHelp("c", "Close tab"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "Next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("Shift+Tab", "Previous tab"),
	),
	ShowHelp: key.NewBinding(
		key.WithKeys("h", "ctrl+h"),
		key.WithHelp("h", "Help"),
	),
	ToggleOfflineMode: key.NewBinding(
		key.WithKeys("o", "ctrl+o"),
		key.WithHelp("o", "Offline mode"),
	),
}

// SetEnabled allows to disable/enable shortcuts
func (k *Keymap) SetEnabled(enabled bool) {
	k.CloseTab.SetEnabled(enabled)
	k.NextTab.SetEnabled(enabled)
	k.PrevTab.SetEnabled(enabled)
	k.ShowHelp.SetEnabled(enabled)
	k.ToggleOfflineMode.SetEnabled(enabled)
}
