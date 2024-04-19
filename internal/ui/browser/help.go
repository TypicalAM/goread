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

// closeHelpMsg is the message sent when the user presses a button to close the help
type closeHelpMsg struct{}

// Help is a popup that displays the help page.
type Help struct {
	border   popup.TitleBorder
	help     help.Model
	box      lipgloss.Style
	keyBinds [][]key.Binding
	width    int
	height   int
}

// newHelp returns a new Help popup.
func newHelp(colors *theme.Colors, binds [][]key.Binding) *Help {
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
	height := strings.Count(rendered, "\n") + 5

	border := popup.NewTitleBorder("Help", width, height, colors.Color1, lipgloss.NormalBorder())
	return &Help{
		help:     helpModel,
		border:   border,
		box:      lipgloss.NewStyle().Margin(1, 2, 1, 4),
		keyBinds: binds,
		width:    width,
		height:   height,
	}
}

// GetSize returns the size of the popup.
func (h Help) GetSize() (width int, height int) {
	return h.width, h.height
}

// Init initializes the popup.
func (h Help) Init() tea.Cmd {
	return nil
}

// Update updates the popup, in this case it's just static text.
func (h Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		return h, h.confirm()
	}

	return h, nil
}

// View renders the popup.
func (h Help) View() string {
	list := h.box.Render(h.help.FullHelpView(h.keyBinds))
	return h.border.Render(list)
}

// confirm returns a tea.Cmd that tells the parent model about the confirmation.
func (h Help) confirm() tea.Cmd {
	return func() tea.Msg { return closeHelpMsg{} }
}
