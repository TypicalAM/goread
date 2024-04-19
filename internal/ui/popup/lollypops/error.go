package lollypops

import (
	"github.com/TypicalAM/goread/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ErrorResultMsg is the message sent when the user presses ok
type ErrorResultMsg struct{}

// AppError is a popup that presents an error to the user.
type AppError struct {
	style  errorStyle
	msg    string
	width  int
	height int
}

// NewError creates a new error popup.
func NewError(colors *theme.Colors, message string) AppError {
	width := len(message) + 16
	height := 7

	return AppError{
		style:  newErrorStyle(colors, width, height),
		msg:    message,
		width:  width,
		height: height,
	}
}

// Init initializes the popup.
func (ae AppError) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (ae AppError) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
		return ae, ae.confirm()
	}

	return ae, nil
}

// View renders the popup.
func (ae AppError) View() string {
	button := ae.style.activeButton.Render("OK")
	msg := ae.style.msg.Render(ae.msg)
	ui := lipgloss.JoinVertical(lipgloss.Center, msg, button)
	dialog := lipgloss.Place(ae.width-2, ae.height-2, lipgloss.Center, lipgloss.Center, ui)
	return ae.style.border.Render(dialog)
}

// GetSize returns the size of the popup.
func (ae AppError) GetSize() (width, height int) {
	return ae.width, ae.height
}

// confirm returns a tea.Cmd that tells the parent model about the confirmation.
func (ae AppError) confirm() tea.Cmd {
	return func() tea.Msg { return ErrorResultMsg{} }
}
