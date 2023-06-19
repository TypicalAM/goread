package popup

import (
<<<<<<< HEAD
	"fmt"
	"strings"

=======
	"strings"

	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/bubbles/textinput"
>>>>>>> 5ea651a (feat: added a basic popup window)
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/ansi"
)

<<<<<<< HEAD
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
=======
// focusedPopupOption is the currently focused popup option.
type focusedPopupOption int

const (
	allField focusedPopupOption = iota
	downloadedField
	newCategoryField
)

// Popup is a popup window allowing the user to select an item from a list of items.
type Popup struct {
	ogSection    []string
	section      []string
	width        int
	height       int
	prefix       string
	suffix       string
	colors       colorscheme.Colorscheme
	startCol     int
	textInput    textinput.Model
	renderStyle  lipgloss.Style
	focusedField focusedPopupOption
}

// New creates a new popup window.
func New(colors colorscheme.Colorscheme, bgRaw string, width, height int) Popup {
>>>>>>> 5ea651a (feat: added a basic popup window)
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

<<<<<<< HEAD
=======
	text := textinput.New()

>>>>>>> 5ea651a (feat: added a basic popup window)
	return Popup{
		ogSection: ogSection,
		section:   section,
		width:     width,
		height:    height,
		prefix:    prefix,
		suffix:    suffix,
		startCol:  startCol,
<<<<<<< HEAD
=======
		colors:    colors,
		textInput: text,
>>>>>>> 5ea651a (feat: added a basic popup window)
		renderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Width(width - 2).
			Height(height - 2).
			Border(lipgloss.NormalBorder()).
<<<<<<< HEAD
			BorderForeground(lipgloss.Color("#888B7E")),
=======
			BorderForeground(colors.Color1),
>>>>>>> 5ea651a (feat: added a basic popup window)
	}
}

// Init the popup window.
func (p Popup) Init() tea.Cmd {
<<<<<<< HEAD
	return nil
=======
	return textinput.Blink
>>>>>>> 5ea651a (feat: added a basic popup window)
}

// Update the popup window.
func (p Popup) Update(msg tea.Msg) (Popup, tea.Cmd) {
<<<<<<< HEAD
	_ = msg
	return p, nil
=======
	var cmds []tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "down", "j":
			switch p.focusedField {
			case allField:
				p.focusedField = downloadedField
			case downloadedField:
				p.focusedField = newCategoryField
			case newCategoryField:
				p.focusedField = allField
				cmds = append(cmds, p.textInput.Focus())
			}

		case "up", "k":
			switch p.focusedField {
			case allField:
				p.focusedField = newCategoryField
			case downloadedField:
				p.focusedField = allField
			case newCategoryField:
				p.focusedField = downloadedField
			}
		}
	}

	if p.textInput.Focused() {
		var cmd tea.Cmd
		p.textInput, cmd = p.textInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
>>>>>>> 5ea651a (feat: added a basic popup window)
}

// Render the popup window.
func (p Popup) View() string {
	// Question
	headingStyle := lipgloss.NewStyle().
<<<<<<< HEAD
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

=======
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
	if p.focusedField == allField {
		renderedChoices[0] = selectedChoiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[0], descriptions[0]))
	} else {
		renderedChoices[0] = choiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[0], descriptions[0]))
	}

	if p.focusedField == downloadedField {
		renderedChoices[1] = selectedChoiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[1], descriptions[1]))
	} else {
		renderedChoices[1] = choiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[1], descriptions[1]))
	}

	if p.focusedField == newCategoryField {
		renderedChoices[2] = selectedChoiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[2], descriptions[2]))
	} else {
		if p.textInput.Focused() {
			renderedChoices[2] = choiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[2], "focused", p.textInput.View()))
		}
		renderedChoices[2] = choiceStyle.Render(lipgloss.JoinVertical(lipgloss.Top, choices[2], p.textInput.View()))
	}

	toBox := choiceSectionStyle.Render(lipgloss.JoinVertical(lipgloss.Top, renderedChoices...))
	popup := lipgloss.JoinVertical(lipgloss.Top, question, toBox)
>>>>>>> 5ea651a (feat: added a basic popup window)
	popupSplit := strings.Split(p.renderStyle.Render(popup), "\n")

	// Overlay the background with the styled text.
	for i, text := range p.ogSection {
		p.section[i] = text[:findPrintableIndex(text, p.startCol)] +
			popupSplit[i] +
			text[findPrintableIndex(text, p.startCol+p.width):]
	}

<<<<<<< HEAD
	return fmt.Sprintf("%s\n%s\n%s", p.prefix, strings.Join(p.section, "\n"), p.suffix)
=======
	return lipgloss.JoinVertical(
		lipgloss.Top,
		p.prefix,
		strings.Join(p.section, "\n"),
		p.suffix,
	)
>>>>>>> 5ea651a (feat: added a basic popup window)
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
