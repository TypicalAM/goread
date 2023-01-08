package browser

import (
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/TypicalAM/goread/internal/style"
	"github.com/charmbracelet/lipgloss"
)

var (
	activeTab = lipgloss.NewStyle().
			Padding(0, 7, 0, 1).
			Italic(true).
			Bold(true)

	activeTabIcon = lipgloss.NewStyle().
			Padding(0, 0, 0, 3).
			Bold(true).
			Border(lipgloss.Border{Left: "┃"}, false, false, false, true).
			BorderForeground(style.GlobalColorscheme.TextDark)

	tabStyle = lipgloss.NewStyle().
			Padding(0, 7, 0, 1).
			Background(style.GlobalColorscheme.BgDark).
			Foreground(style.GlobalColorscheme.TextDark)

	tabIcon = activeTabIcon.Copy().
		Background(style.GlobalColorscheme.BgDark).
		BorderForeground(style.GlobalColorscheme.BgDarker).
		BorderBackground(style.GlobalColorscheme.BgDark)

	tabGap = lipgloss.NewStyle().
		Background(style.GlobalColorscheme.BgDarker)

	statusBarGap = lipgloss.NewStyle().
			Background(style.GlobalColorscheme.BgDark)

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
		iconStyle = activeTabIcon
		textStyle = activeTab
	} else {
		iconStyle = tabIcon
		textStyle = tabStyle
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
