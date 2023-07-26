package popup

import (
	"github.com/TypicalAM/goread/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChoiceResultMsg is the message sent when a choice is made.
type ChoiceResultMsg struct {
	Result bool
}

// Choice is a popup that presents a yes/no choice to the user.
type Choice struct {
	style        style
	question     string
	defaultPopup Default
	selected     bool
}

// NewChoice creates a new Choice popup.
func NewChoice(colors *theme.Colors, bgRaw string, width int, question string, defaultChoice bool) Choice {
	optWidth := len(question) + 16
	if optWidth > width {
		optWidth = width
	}

	height := 7

	return Choice{
		style:        newStyle(colors, optWidth, height),
		defaultPopup: New(bgRaw, optWidth, height),
		question:     question,
		selected:     defaultChoice,
	}
}

// Init initializes the popup.
func (c Choice) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (c Choice) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "enter":
			return c, c.makeChoice()

		case "left", "right", "tab":
			c.selected = !c.selected
			return c, nil

		case "y", "Y":
			c.selected = true
			return c, c.makeChoice()

		case "n", "N":
			c.selected = false
			return c, c.makeChoice()
		}
	}

	return c, nil
}

// View renders the popup.
func (c Choice) View() string {
	var okButton, cancelButton string
	if c.selected {
		okButton = c.style.activeButton.Render("Yes")
		cancelButton = c.style.button.Render("No")
	} else {
		okButton = c.style.button.Render("Yes")
		cancelButton = c.style.activeButton.Render("No")
	}

	question := c.style.question.Render(c.question)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)
	dialog := lipgloss.Place(c.defaultPopup.width-2, c.defaultPopup.height-2, lipgloss.Center, lipgloss.Center, ui)

	return c.defaultPopup.Overlay(c.style.general.Render(dialog))
}

// makeChoice returns a tea.Cmd that tells the parent model about the choice.
func (c Choice) makeChoice() tea.Cmd {
	return func() tea.Msg {
		return ChoiceResultMsg{c.selected}
	}
}
