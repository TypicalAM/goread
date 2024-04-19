package overview

import (
	"github.com/TypicalAM/goread/internal/backend/rss"
	"github.com/TypicalAM/goread/internal/theme"
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
	nameInput textinput.Model
	descInput textinput.Model
	style     popupStyle
	oldName   string
	focused   focusedField
	editing   bool
	width     int
	height    int
	reserved  bool
}

// NewPopup creates a new popup window in which the user can choose a new category.
func NewPopup(colors *theme.Colors, oldName, oldDesc string) Popup {
	width := 46
	height := 14

	editing := oldName != "" || oldDesc != ""
	reserved := oldName == rss.AllFeedsName || oldName == rss.DownloadedFeedsName

	nameInput := textinput.New()
	nameInput.CharLimit = 30
	nameInput.Width = width - 15
	nameInput.Prompt = "Name: "
	descInput := textinput.New()
	descInput.CharLimit = 30
	descInput.Width = width - 22
	descInput.Prompt = "Description: "

	focused := allField
	if oldName == rss.DownloadedFeedsName {
		focused = downloadedField
	}

	style := popupStyle{}
	if editing {
		style = newPopupStyle(colors, width, height, "Edit category")
	} else {
		style = newPopupStyle(colors, width, height, "New category")
	}

	if editing && !reserved {
		nameInput.SetValue(oldName)
		descInput.SetValue(oldDesc)
		focused = nameField
		nameInput.Focus()
	}

	return Popup{
		style:     style,
		nameInput: nameInput,
		descInput: descInput,
		oldName:   oldName,
		focused:   focused,
		editing:   editing,
		width:     width,
		height:    height,
		reserved:  reserved,
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
		if p.reserved {
			return p, nil
		}

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
				if p.editing {
					return p, nil
				}

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
				if p.editing {
					return p, nil
				}

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
				return p, confirm(rss.AllFeedsName, "", "", false)

			case downloadedField:
				return p, confirm(rss.DownloadedFeedsName, "", "", false)

			case nameField, descField:
				return p, confirm(p.nameInput.Value(), p.descInput.Value(), p.oldName, p.oldName != "")
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
	renderedChoices := make([]string, 3)

	titles := []string{rss.AllFeedsName, rss.DownloadedFeedsName, "New category"}
	descs := []string{"All available articles", "Downloaded articles", p.nameInput.View() + "\n" + p.descInput.View()}

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

	toList := p.style.list.Render(lipgloss.JoinVertical(lipgloss.Top, renderedChoices...))
	return p.style.border.Render(toList)
}

// GetSize returns the size of the popup.
func (p Popup) GetSize() (width, height int) {
	return p.width, p.height
}

// confirm creates a message that confirms the user's choice.
func confirm(name, desc, oldName string, edit bool) tea.Cmd {
	return func() tea.Msg { return ChosenCategoryMsg{name, desc, oldName, edit} }
}
