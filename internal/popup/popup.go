package popup

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/ansi"
)

// Popup is a popup window allowing the user to select an item from a list of items.
type Popup struct {
	ogSection   []string
	section     []string
	width       int
	height      int
	prefix      string
	suffix      string
	startCol    int
	renderStyle lipgloss.Style
}

// New creates a new popup window.
func New(bgRaw string, width, height int) Popup {
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

	return Popup{
		ogSection: ogSection,
		section:   section,
		width:     width,
		height:    height,
		prefix:    prefix,
		suffix:    suffix,
		startCol:  startCol,
		renderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Width(width - 2).
			Height(height - 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#888B7E")),
	}
}

// Init the popup window.
func (p Popup) Init() tea.Cmd {
	return nil
}

// Update the popup window.
func (p Popup) Update(msg tea.Msg) (Popup, tea.Cmd) {
	_ = msg
	return p, nil
}

// Render the popup window.
func (p Popup) View() string {
	// Question
	headingStyle := lipgloss.NewStyle().
		Width(p.width-2).
		Margin(2, 2).
		Align(lipgloss.Center).
		Italic(true)
	question := headingStyle.Render("Which category do you want to add?")

	itemStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		MarginBottom(1)

	activeTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e56996")).
		Width(p.width - 2).
		Underline(true)

	activeDescStyle := lipgloss.NewStyle().
		Width(p.width - 2).
		Foreground(lipgloss.Color("#FFFFFF"))

	activeItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#e56996")).
		PaddingLeft(2).
		Border(lipgloss.NormalBorder(), false, false, false, true)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888B7E"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888B7E"))

	textDesc := activeItemStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		activeTextStyle.Render("All Categories"),
		activeDescStyle.Render("All categories are shown."),
	))

	item1 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		textDesc,
	)

	textDesc = lipgloss.JoinVertical(lipgloss.Left, textStyle.Render("Downloaded"), descStyle.Render("Downloaded categories are shown."))
	item2 := itemStyle.Render(textDesc)

	textDesc = lipgloss.JoinVertical(lipgloss.Top, textStyle.Render("Not Downloaded"), descStyle.Render("Not downloaded categories are shown."))
	item3 := itemStyle.Render(textDesc)

	popup := lipgloss.JoinVertical(lipgloss.Top, question, item1, item2, item3)

	popupSplit := strings.Split(p.renderStyle.Render(popup), "\n")

	// Overlay the background with the styled text.
	for i, text := range p.ogSection {
		p.section[i] = text[:findPrintableIndex(text, p.startCol)] +
			popupSplit[i] +
			text[findPrintableIndex(text, p.startCol+p.width):]
	}

	return fmt.Sprintf("%s\n%s\n%s", p.prefix, strings.Join(p.section, "\n"), p.suffix)
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
