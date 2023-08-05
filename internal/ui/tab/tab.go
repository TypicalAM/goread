package tab

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Style is a struct that holds the style of a tab
type Style struct {
	Color lipgloss.Color
	Icon  string
	Name  string
}

// Tab is an interface outlining the methods that a tab should implement including bubbletea's model methods
type Tab interface {
	tea.Model
	help.KeyMap

	Title() string
	Style() Style
	SetSize(width, height int) Tab
}
