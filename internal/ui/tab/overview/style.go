package overview

import (
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/TypicalAM/goread/internal/ui/popup"
	"github.com/charmbracelet/lipgloss"
)

// popupStyle is the style of the popup window.
type popupStyle struct {
	border              popup.TitleBorder
	list                lipgloss.Style
	choice              lipgloss.Style
	choiceTitle         lipgloss.Style
	choiceDesc          lipgloss.Style
	selectedChoice      lipgloss.Style
	selectedChoiceTitle lipgloss.Style
	selectedChoiceDesc  lipgloss.Style
}

// newPopupStyle creates a new popup style.
func newPopupStyle(colors *theme.Colors, width, height int, headingText string) popupStyle {
	border := popup.NewTitleBorder(headingText, width, height, colors.Color1, lipgloss.NormalBorder())

	list := lipgloss.NewStyle().
		Margin(1, 4).
		Width(width - 2).
		Height(10)

	choice := lipgloss.NewStyle().
		PaddingLeft(1).
		MarginBottom(1).
		Border(lipgloss.RoundedBorder(), false, false, false, true)

	choiceTitle := lipgloss.NewStyle().
		Foreground(colors.Text)

	choiceDesc := lipgloss.NewStyle().
		Foreground(colors.TextDark)

	selectedChoice := choice.Copy().
		Italic(true).
		BorderForeground(colors.Color3)

	selectedChoiceTitle := lipgloss.NewStyle().
		Foreground(colors.Color3)

	selectedChoiceDesc := lipgloss.NewStyle().
		Foreground(colors.Color2)

	return popupStyle{
		border:              border,
		list:                list,
		choice:              choice,
		choiceTitle:         choiceTitle,
		choiceDesc:          choiceDesc,
		selectedChoice:      selectedChoice,
		selectedChoiceTitle: selectedChoiceTitle,
		selectedChoiceDesc:  selectedChoiceDesc,
	}
}
