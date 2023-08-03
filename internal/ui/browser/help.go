package browser

import (
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/TypicalAM/goread/internal/ui/popup"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Help is a popup that displays the help page.
type Help struct {
	help     help.Model
	style    helpStyle
	keyBinds [][]key.Binding
	overlay  popup.Overlay
}

// newHelp returns a new Help popup.
func newHelp(colors *theme.Colors, bgRaw string, width, height int, binds [][]key.Binding) *Help {
	style := newHelpStyle(colors, width, height)
	help := help.New()
	help.Styles = style.help

	return &Help{
		help:     help,
		style:    style,
		keyBinds: binds,
		overlay:  popup.NewOverlay(bgRaw, width, height),
	}
}

// Init initializes the popup.
func (h Help) Init() tea.Cmd {
	return nil
}

// Update updates the popup, in this case it's just static text.
func (h Help) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return h, nil
}

// View renders the popup.
func (h Help) View() string {
	return h.overlay.WrapView(h.style.box.Render(lipgloss.JoinVertical(lipgloss.Center,
		h.style.title.Render("Help"),
		h.help.FullHelpView(h.keyBinds),
	)))
}
