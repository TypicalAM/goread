package popup

import (
	"strings"

	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/rss"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/ansi"
)

// ChosenCategoryMsg is the message displayed when a category is successfully chosen.
type ChosenCategoryMsg struct {
	Name string
}

// focusedField is the field that is currently focused.
type focusedField int

const (
	allField focusedField = iota
	downloadedField
	newCategoryField
)

// Popup is a popup window allowing the user to select an item from a list of items.
type Popup struct {
	ogSection   []string
	section     []string
	width       int
	height      int
	prefix      string
	suffix      string
	colors      colorscheme.Colorscheme
	textInput   textinput.Model
	focused     focusedField
	startCol    int
	renderStyle lipgloss.Style
}

// New creates a new popup window.
func New(colors colorscheme.Colorscheme, bgRaw string, width, height int) Popup {
	bg := strings.Split(bgRaw, "\n")
	bgWidth := ansi.PrintableRuneWidth(bg[0])
	bgHeight := len(bg)

	startRow := (bgHeight - height) / 2
	startCol := (bgWidth - width) / 2

	ogSection := make([]string, height)
	section := make([]string, height)

	prefix := strings.Join(bg[:startRow], "\n")
	suffix := strings.Join(bg[startRow+height:], "\n")
	copy(ogSection, bg[startRow:startRow+height])

	text := textinput.New()

	return Popup{
		ogSection: ogSection,
		section:   section,
		width:     width,
		height:    height,
		prefix:    prefix,
		suffix:    suffix,
		startCol:  startCol,
		colors:    colors,
		textInput: text,
		renderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Width(width - 2).
			Height(height - 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(colors.Color1),
	}
}

// Init the popup window.
func (p Popup) Init() tea.Cmd {
	return textinput.Blink
}

// Update the popup window.
func (p Popup) Update(msg tea.Msg) (Popup, tea.Cmd) {
	var cmds []tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "down", "j":
			switch p.focused {
			case allField:
				p.focused = downloadedField
			case downloadedField:
				p.focused = newCategoryField
				cmds = append(cmds, p.textInput.Focus())
			case newCategoryField:
				p.focused = allField
			}

		case "up", "k":
			switch p.focused {
			case allField:
				p.focused = newCategoryField
				cmds = append(cmds, p.textInput.Focus())
			case downloadedField:
				p.focused = allField
			case newCategoryField:
				p.focused = downloadedField
			}

		case "enter":
			switch p.focused {
			case allField:
				return p, confirmCategory(rss.AllFeedsName)

			case downloadedField:
				return p, confirmCategory(rss.DownloadedFeedsName)

			case newCategoryField:
				// TODO: Validate the name
				return p, confirmCategory(p.textInput.Value())
			}
		}
	}

	if p.textInput.Focused() {
		var cmd tea.Cmd
		p.textInput, cmd = p.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

// Render the popup window.
func (p Popup) View() string {
	// Question
	headingStyle := lipgloss.NewStyle().
		Margin(2, 2).
		Width(p.width - 2).
		Align(lipgloss.Center).
		Italic(true)

	question := headingStyle.Render("Which category do you want to add?")

	choices := []string{"All", "Downloaded", "New Category"}
	descriptions := []string{
		"All the feeds from all the categories",
		"Downloaded articles",
		"",
	}
	choiceSectionStyle := lipgloss.NewStyle().
		Padding(2).
		Width(p.width - 2).
		Height(10)

	choiceStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		MarginBottom(1).
		Border(lipgloss.RoundedBorder(), false, false, false, true)

	selectedChoiceStyle := choiceStyle.Copy().
		BorderForeground(p.colors.Color4)

	renderedChoices := make([]string, len(choices))
	if p.focused == allField {
		renderedChoices[0] = selectedChoiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[0], descriptions[0]))
	} else {
		renderedChoices[0] = choiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[0], descriptions[0]))
	}

	if p.focused == downloadedField {
		renderedChoices[1] = selectedChoiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[1], descriptions[1]))
	} else {
		renderedChoices[1] = choiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[1], descriptions[1]))
	}

	if p.focused == newCategoryField {
		renderedChoices[2] = selectedChoiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[2], p.textInput.View()))
	} else {
		renderedChoices[2] = choiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[2], p.textInput.View()))
	}

	toBox := choiceSectionStyle.Render(lipgloss.JoinVertical(lipgloss.Top, renderedChoices...))
	popup := lipgloss.JoinVertical(lipgloss.Top, question, toBox)
	popupSplit := strings.Split(p.renderStyle.Render(popup), "\n")

	// Overlay the background with the styled text.
	for i, text := range p.ogSection {
		p.section[i] = text[:findPrintableIndex(text, p.startCol)] +
			popupSplit[i] +
			text[findPrintableIndex(text, p.startCol+p.width):]
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		p.prefix,
		strings.Join(p.section, "\n"),
		p.suffix,
	)
}

// findPrintableIndex finds the index of the last printable rune at the given index.
func findPrintableIndex(str string, index int) int {
	for i := len(str) - 1; i >= 0; i-- {
		if ansi.PrintableRuneWidth(str[:i]) == index {
			return i
		}
	}
	return -1
}

// confirmCategory returns a tea.Cmd which relays the message to the browser.
func confirmCategory(name string) tea.Cmd {
	return func() tea.Msg {
		return ChosenCategoryMsg{Name: name}
	}
}
