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
	Name    string
	Desc    string
	OldName string
	IsEdit  bool
}

// focusedField is the field that is currently focused.
type focusedField int

const (
	allField focusedField = iota
	downloadedField
	nameField
	descField
)

// Popup is the category popup where a user can create a category.
type Popup struct {
	defaultPopup popup.Default
	style        popupStyle
	nameInput    textinput.Model
	descInput    textinput.Model
	focused      focusedField
	oldName      string
}

// NewPopup creates a new popup window in which the user can choose a new category.
func NewPopup(colors colorscheme.Colorscheme, bgRaw string, width, height int, oldName, oldDesc string) Popup {
	defultPopup := popup.New(bgRaw, width, height)
	style := newPopupStyle(colors, width, height)
	nameInput := textinput.New()
	nameInput.CharLimit = 30
	nameInput.Width = width - 15
	nameInput.Prompt = "Name: "
	descInput := textinput.New()
	descInput.CharLimit = 30
	descInput.Width = width - 22
	descInput.Prompt = "Description: "
	focusedField := allField

	if oldName != "" || oldDesc != "" {
		nameInput.SetValue(oldName)
		descInput.SetValue(oldDesc)
		focusedField = nameField
		nameInput.Focus()
	}

	return Popup{
		defaultPopup: defultPopup,
		style:        style,
		nameInput:    nameInput,
		descInput:    descInput,
		oldName:      oldName,
		focused:      focusedField,
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
		case "down", "tab":
			switch p.focused {
			case allField:
				p.focused = downloadedField
			case downloadedField:
				p.focused = nameField
				cmds = append(cmds, p.nameInput.Focus())
			case nameField:
				p.focused = descField
				p.nameInput.Blur()
				cmds = append(cmds, p.descInput.Focus())
			case descField:
				p.focused = allField
				p.descInput.Blur()
			}

		case "up":
			switch p.focused {
			case allField:
				p.focused = descField
				cmds = append(cmds, p.descInput.Focus())
			case downloadedField:
				p.focused = allField
			case nameField:
				p.focused = downloadedField
				p.nameInput.Blur()
			case descField:
				p.focused = nameField
				p.descInput.Blur()
				cmds = append(cmds, p.nameInput.Focus())
			}

		case "enter":
			switch p.focused {
			case allField:
				return p, confirmCategory(rss.AllFeedsName, "", "", false)

			case downloadedField:
				return p, confirmCategory(rss.DownloadedFeedsName, "", "", false)

			case nameField, descField:
				return p, confirmCategory(p.nameInput.Value(), p.descInput.Value(), p.oldName, p.oldName != "")
			}
		}
	}

	if p.nameInput.Focused() {
		var cmd tea.Cmd
		p.nameInput, cmd = p.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	if p.descInput.Focused() {
		var cmd tea.Cmd
		p.descInput, cmd = p.descInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

// View renders the popup window.
func (p Popup) View() string {
	question := p.style.heading.Render("Choose a category")
	renderedChoices := make([]string, 3)

	titles := []string{rss.AllFeedsName, rss.DownloadedFeedsName, "New category"}
	descs := []string{"All the feeds", "Saved Feeds", p.nameInput.View() + "\n" + p.descInput.View()}

	var focused int
	switch p.focused {
	case allField:
		focused = 0
	case downloadedField:
		focused = 1
	case nameField, descField:
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
func confirmCategory(name, desc, oldName string, isEdit bool) tea.Cmd {
	return func() tea.Msg {
		return ChosenCategoryMsg{
			Name:    name,
			Desc:    desc,
			OldName: oldName,
			IsEdit:  isEdit,
		}
	}
}
