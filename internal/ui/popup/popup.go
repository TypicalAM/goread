package popup

import tea "github.com/charmbracelet/bubbletea"

// Window represents a popup window.
type Window interface {
	tea.Model

	GetSize() (width, height int)
}
