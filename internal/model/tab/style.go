package tab

import (
	"github.com/TypicalAM/goread/internal/style"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusBarCell = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Foreground(style.GlobalColorscheme.BgDark)

	iconColors = map[Type]lipgloss.Color{
		Welcome:  style.GlobalColorscheme.Color4,
		Category: style.GlobalColorscheme.Color5,
		Feed:     style.GlobalColorscheme.Color3,
	}

	icons = map[Type]string{
		Welcome:  "﫢",
		Category: "﫜",
		Feed:     "",
	}

	texts = map[Type]string{
		Welcome:  "MAIN",
		Category: "CATEGORY",
		Feed:     "FEED",
	}
)

// Style the text depending on the type of the tab
func AttachIconToTab(text string, tabType Type, isActive bool) string {
	var iconStyle, textStyle lipgloss.Style
	if isActive {
		iconStyle = style.ActiveTabIcon
		textStyle = style.ActiveTab
	} else {
		iconStyle = style.TabIcon
		textStyle = style.TabStyle
	}

	// Cut the text if the tab length is too much to handle
	if len(text) > 12 {
		text = text[:12] + ""
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		iconStyle.Copy().Foreground(iconColors[tabType]).Render(icons[tabType]),
		textStyle.Render(text),
	)
}

// Style the status bar cell depending on the the of the current tab
func StyleStatusBarCell(tabType Type) string {
	return statusBarCell.Copy().
		Background(iconColors[tabType]).
		Render(texts[tabType])
}
