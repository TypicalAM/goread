package model

import (
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/TypicalAM/goread/internal/style"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusBarCell = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Foreground(style.GlobalColorscheme.BgDark)

	iconColors = map[tab.Type]lipgloss.Color{
		tab.Welcome:  style.GlobalColorscheme.Color4,
		tab.Category: style.GlobalColorscheme.Color5,
		tab.Feed:     style.GlobalColorscheme.Color3,
	}

	icons = map[tab.Type]string{
		tab.Welcome:  "﫢",
		tab.Category: "﫜",
		tab.Feed:     "",
	}

	texts = map[tab.Type]string{
		tab.Welcome:  "MAIN",
		tab.Category: "CATEGORY",
		tab.Feed:     "FEED",
	}
)

// Style the text depending on the type of the tab
func attachIconToTab(text string, tabType tab.Type, isActive bool) string {
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
func styleStatusBarCell(tabType tab.Type) string {
	return statusBarCell.Copy().
		Background(iconColors[tabType]).
		Render(texts[tabType])
}
