package category

import (
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChosenFeedMsg is the message displayed when a category is successfully chosen.
type ChosenFeedMsg struct {
	Name    string
	URL     string
	OldName string
	Parent  string
	IsEdit  bool
}

// focusedField is the field that is currently focused.
type focusedField int

const (
	nameField focusedField = iota
	urlField
)

// Popup is the feed popup where a user can create/edit a feed.
type Popup struct {
	nameInput textinput.Model
	urlInput  textinput.Model
	style     popupStyle
	oldName   string
	oldURL    string
	parent    string
	focused   focusedField
	editing   bool
	width     int
	height    int
}

// NewPopup returns a new feed popup.
func NewPopup(colors *theme.Colors, oldName, oldURL, parent string) Popup {
	width := 40
	height := 7

	editing := oldName != "" || oldURL != ""

	nameInput := textinput.New()
	nameInput.CharLimit = 30
	nameInput.Prompt = "Name: "
	nameInput.Width = width - 20
	urlInput := textinput.New()
	urlInput.CharLimit = 150
	urlInput.Width = width - 20
	urlInput.Prompt = "URL: "

	style := popupStyle{}
	if editing {
		style = newPopupStyle(colors, width, height, "Edit feed")
	} else {
		style = newPopupStyle(colors, width, height, "New feed")
	}

	if editing {
		nameInput.SetValue(oldName)
		urlInput.SetValue(oldURL)
	}

	nameInput.Focus()

	return Popup{
		style:     style,
		nameInput: nameInput,
		urlInput:  urlInput,
		oldName:   oldName,
		oldURL:    oldURL,
		parent:    parent,
		editing:   editing,
		width:     width,
		height:    height,
	}
}

// Init initializes the popup.
func (p Popup) Init() tea.Cmd {
	return textinput.Blink
}

// Update updates the popup.
func (p Popup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "down", "up", "tab":
			switch p.focused {
			case nameField:
				p.focused = urlField
				p.nameInput.Blur()
				cmds = append(cmds, p.urlInput.Focus())

			case urlField:
				p.focused = nameField
				p.urlInput.Blur()
				cmds = append(cmds, p.nameInput.Focus())
			}

		case "enter":
			return p, confirm(
				p.nameInput.Value(),
				p.urlInput.Value(),
				p.oldName,
				p.parent,
				p.oldName != "",
			)
		}
	}

	if p.nameInput.Focused() {
		var cmd tea.Cmd
		p.nameInput, cmd = p.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	if p.urlInput.Focused() {
		var cmd tea.Cmd
		p.urlInput, cmd = p.urlInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

// View renders the popup.
func (p Popup) View() string {
	itemText := ""
	if p.editing {
		itemText = "Your feed"
	} else {
		itemText = "New feed"
	}

	itemTitle := p.style.itemTitle.Render(itemText)
	name := p.style.itemField.Render(p.nameInput.View())
	url := p.style.itemField.Render(p.urlInput.View())
	listItem := p.style.listItem.Render(lipgloss.JoinVertical(lipgloss.Left, itemTitle, name, url))
	return p.style.border.Render(listItem)
}

// GetSize returns the size of the popup.
func (p Popup) GetSize() (width, height int) {
	return p.width, p.height
}

// confirm creates a message that confirms the user's choice.
func confirm(name, url, oldName, parent string, edit bool) tea.Cmd {
	return func() tea.Msg { return ChosenFeedMsg{name, url, oldName, parent, edit} }
}
