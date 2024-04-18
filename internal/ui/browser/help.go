package browser

import (
	"strings"

	"github.com/TypicalAM/goread/internal/theme"
	"github.com/TypicalAM/goread/internal/ui/popup"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/ansi"
)

// Help is a popup that displays the help page.
type Help struct {
	border   popup.TitleBorder
	help     help.Model
	box      lipgloss.Style
	keyBinds [][]key.Binding
	overlay  popup.Overlay
}

// newHelp returns a new Help popup.
func newHelp(colors *theme.Colors, bgRaw string, binds [][]key.Binding) *Help {
	helpModel := help.New()
	helpModel.Styles = help.Styles{}
	helpModel.Styles.FullDesc = lipgloss.NewStyle().
		Foreground(colors.Text)
	helpModel.Styles.FullKey = lipgloss.NewStyle().
		Foreground(colors.Color2)
	helpModel.Styles.FullSeparator = lipgloss.NewStyle().
		Foreground(colors.TextDark)

	rendered := helpModel.FullHelpView(binds)
	width := ansi.PrintableRuneWidth(rendered[:strings.IndexRune(rendered, '\n')-1]) + 6
	height := strings.Count(rendered, "\n") + 7

	border := popup.NewTitleBorder("Help", width, height, colors.Color1, lipgloss.NormalBorder())
	return &Help{
		help:     helpModel,
		border:   border,
		box:      lipgloss.NewStyle().Margin(2, 2),
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
	list := h.box.Render(h.help.FullHelpView(h.keyBinds))
	return h.overlay.WrapView(h.border.Render(list))
}
