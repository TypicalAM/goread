package feed

import (
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/popup"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChosenFeedMsg is the message displayed when a category is successfully chosen.
type ChosenFeedMsg struct {
	Name           string
	URL            string
	OldName        string
	ParentCategory string
	IsEdit         bool
}

// focusedField is the field that is currently focused.
type focusedField int

const (
	nameField focusedField = iota
	urlField
)

// Popup is the feed popup where a user can create/edit a feed.
type Popup struct {
	defaultPopup   popup.Default
	style          popupStyle
	nameInput      textinput.Model
	urlInput       textinput.Model
	oldName        string
	oldURL         string
	parentCategory string
	focused        focusedField
}

// NewPopup returns a new feed popup.
func NewPopup(colors colorscheme.Colorscheme, bgRaw string, width, height int,
	oldName, oldURL, parentCategory string) Popup {

	style := newPopupStyle(colors, width, height)
	defaultPopup := popup.New(bgRaw, width, height)
	nameInput := textinput.New()
	URLInput := textinput.New()

	if oldName != "" || oldURL != "" {
		nameInput.SetValue(oldName)
		URLInput.SetValue(oldURL)
	}

	nameInput.Focus()

	return Popup{
		defaultPopup:   defaultPopup,
		style:          style,
		nameInput:      nameInput,
		urlInput:       URLInput,
		oldName:        oldName,
		oldURL:         oldURL,
		parentCategory: parentCategory,
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
		case "down", "up":
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
			return p, confirmFeed(
				p.nameInput.Value(),
				p.urlInput.Value(),
				p.oldName,
				p.parentCategory,
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
	question := p.style.heading.Render("Add/Edit a feed")
	title := p.style.itemTitle.Render("Write a name and an URL for the feed")
	name := p.style.itemField.Render(p.nameInput.View())
	url := p.style.itemField.Render(p.urlInput.View())
	item := p.style.item.Render(lipgloss.JoinVertical(lipgloss.Top, title, name, url))
	popup := lipgloss.JoinVertical(lipgloss.Top, question, item)
	return p.defaultPopup.Overlay(p.style.general.Render(popup))
}

// confirmFeed returns a tea.Cmd which relays the message to the browser.
func confirmFeed(name, url, oldName, parentCategory string, isEdit bool) tea.Cmd {
	return func() tea.Msg {
		return ChosenFeedMsg{
			Name:           name,
			URL:            url,
			OldName:        oldName,
			ParentCategory: parentCategory,
			IsEdit:         isEdit,
		}
	}
}
