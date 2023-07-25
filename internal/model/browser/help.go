package browser

import (
	"github.com/TypicalAM/goread/internal/model/popup"
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Help is a popup that displays the help page.
type Help struct {
	help         help.Model
	boxStyle     lipgloss.Style
	keyBinds     [][]key.Binding
	defaultPopup popup.Default
}

// newHelp returns a new Help popup.
func newHelp(colors *theme.Colors, bgRaw string, width, height int, binds [][]key.Binding) *Help {
	help := help.New()
	help.Styles.FullDesc = lipgloss.NewStyle().Foreground(colors.Text)
	help.Styles.FullKey = lipgloss.NewStyle().Foreground(colors.Text)
	help.Styles.FullSeparator = lipgloss.NewStyle().Foreground(colors.TextDark)

	return &Help{
		help: help,
		boxStyle: lipgloss.NewStyle().
			Width(width - 2).
			Height(height - 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(colors.Color1),
		keyBinds:     binds,
		defaultPopup: popup.New(bgRaw, width, height),
	}
}

// Init initalizes the popup.
func (h Help) Init() tea.Cmd {
	return nil
}

// Update updates the popup, in this case it's just static text
func (h Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return h, nil
}

// View renders the popup.
func (h Help) View() string {
	hView := h.help.FullHelpView(h.keyBinds)
	boxed := h.boxStyle.Render(hView)
	return h.defaultPopup.Overlay(boxed)
}
