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
	userField
)

// Popup is the category popup where a user can create a category.
type Popup struct {
	defaultPopup popup.Default
	style        popupStyle
	textInput    textinput.Model
	focused      focusedField
}

// NewPopup creates a new popup window in which the user can choose a new category.
func NewPopup(colors colorscheme.Colorscheme, bgRaw string, width, height int) Popup {
	return Popup{
		defaultPopup: popup.New(bgRaw, width, height),
		style:        newPopupStyle(colors, width, height),
		textInput:    textinput.New(),
		focused:      allField,
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
				p.focused = userField
				cmds = append(cmds, p.textInput.Focus())
			case userField:
				p.focused = allField
				p.textInput.Blur()
			}

		case "up", "k":
			switch p.focused {
			case allField:
				p.focused = userField
				cmds = append(cmds, p.textInput.Focus())
			case downloadedField:
				p.focused = allField
			case userField:
				p.focused = downloadedField
				p.textInput.Blur()
			}

		case "enter":
			switch p.focused {
			case allField:
				return p, confirmCategory(rss.AllFeedsName)

			case downloadedField:
				return p, confirmCategory(rss.DownloadedFeedsName)

			case userField:
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

// View renders the popup window.
func (p Popup) View() string {
	question := p.style.heading.Render("Choose a category")
	renderedChoices := make([]string, 3)

	titles := []string{rss.AllFeedsName, rss.DownloadedFeedsName, "New category"}
	descs := []string{"All the feeds from all the categories", "Feeds that have been downloaded", p.textInput.View()}

	var focused int
	switch p.focused {
	case allField:
		focused = 0
	case downloadedField:
		focused = 1
	case userField:
		focused = 2
	}

	for i := 0; i < 3; i++ {
		if i == focused {
			renderedChoices[i] = p.style.selectedChoice.Render(lipgloss.JoinVertical(
				lipgloss.Top,
				p.style.selectedChoiceTitle.Render(titles[i]),
				p.style.selectedChoiceDesc.Render(descs[i]),
			))
		} else {
			renderedChoices[i] = p.style.choice.Render(lipgloss.JoinVertical(
				lipgloss.Top,
				p.style.choiceTitle.Render(titles[i]),
				p.style.choiceDesc.Render(descs[i]),
			))
		}
	}

	toBox := p.style.list.Render(lipgloss.JoinVertical(lipgloss.Top, renderedChoices...))
	popup := lipgloss.JoinVertical(lipgloss.Top, question, toBox)
	return p.defaultPopup.Overlay(p.style.general.Render(popup))
}

// confirmCategory returns a tea.Cmd which relays the message to the browser.
func confirmCategory(name string) tea.Cmd {
	return func() tea.Msg {
		return ChosenCategoryMsg{Name: name}
	}
}
