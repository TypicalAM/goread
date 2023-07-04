package category

import (
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/popup"
	"github.com/TypicalAM/goread/internal/rss"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type Popup struct {
	defaultPopup popup.Default
	colors       colorscheme.Colorscheme
	textInput    textinput.Model
	focused      focusedField
	renderStyle  lipgloss.Style
}

// NewPopup creates a new popup window in which the user can choose a new category.
func NewPopup(colors colorscheme.Colorscheme, bgRaw string, width, height int) Popup {
	return Popup{
		defaultPopup: popup.New(bgRaw, width, height),
		colors:       colors,
		textInput:    textinput.New(),
		focused:      allField,
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
func (p Popup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				p.textInput.Blur()
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
				p.textInput.Blur()
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

func (p Popup) View() string {
	headingStyle := lipgloss.NewStyle().
		Margin(2, 2).
		Width(p.defaultPopup.Width() - 2).
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
		Width(p.defaultPopup.Width() - 2).
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
	return p.defaultPopup.Overlay(p.renderStyle.Render(popup))
}

// confirmCategory returns a tea.Cmd which relays the message to the browser.
func confirmCategory(name string) tea.Cmd {
	return func() tea.Msg {
		return ChosenCategoryMsg{Name: name}
	}
}
