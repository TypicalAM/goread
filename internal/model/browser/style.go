package browser

import (
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/charmbracelet/lipgloss"
)

// style is the internal style of the browser
type style struct {
	colors        colorscheme.Colorscheme
	activeTab     lipgloss.Style
	activeTabIcon lipgloss.Style
	tabStyle      lipgloss.Style
	tabIcon       lipgloss.Style
	tabGap        lipgloss.Style
	statusBarGap  lipgloss.Style
	statusBarCell lipgloss.Style
	iconColors    map[tab.Type]lipgloss.Color
	icons         map[tab.Type]string
	texts         map[tab.Type]string
}

// newStyle creates a new style
func newStyle(colors colorscheme.Colorscheme) style {
	activeTab := lipgloss.NewStyle().
		Padding(0, 7, 0, 1).
		Italic(true).
		Bold(true)

	activeTabIcon := lipgloss.NewStyle().
		Padding(0, 0, 0, 3).
		Bold(true).
		Border(lipgloss.Border{Left: "┃"}, false, false, false, true).
		BorderForeground(colors.TextDark)

	tabStyle := lipgloss.NewStyle().
		Padding(0, 7, 0, 1).
		Background(colors.BgDark).
		Foreground(colors.TextDark)

	tabIcon := activeTabIcon.Copy().
		Background(colors.BgDark).
		BorderForeground(colors.BgDarker).
		BorderBackground(colors.BgDark)

	tabGap := lipgloss.NewStyle().
		Background(colors.BgDarker)

	statusBarGap := lipgloss.NewStyle().
		Background(colors.BgDark)

	statusBarCell := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(colors.BgDark)

	iconColors := map[tab.Type]lipgloss.Color{
		tab.Welcome:  colors.Color4,
		tab.Category: colors.Color5,
		tab.Feed:     colors.Color3,
	}

	icons := map[tab.Type]string{
		tab.Welcome:  "﫢",
		tab.Category: "﫜",
		tab.Feed:     "",
	}

	texts := map[tab.Type]string{
		tab.Welcome:  "MAIN",
		tab.Category: "CATEGORY",
		tab.Feed:     "FEED",
	}

	return style{
		colors:        colors,
		activeTab:     activeTab,
		activeTabIcon: activeTabIcon,
		tabStyle:      tabStyle,
		tabIcon:       tabIcon,
		tabGap:        tabGap,
		statusBarGap:  statusBarGap,
		statusBarCell: statusBarCell,
		iconColors:    iconColors,
		icons:         icons,
		texts:         texts,
	}
}

// Style the text depending on the type of the tab
func (s style) attachIconToTab(text string, tabType tab.Type, isActive bool) string {
	var iconStyle, textStyle lipgloss.Style
	if isActive {
		iconStyle = s.activeTabIcon
		textStyle = s.activeTab
	} else {
		iconStyle = s.tabIcon
		textStyle = s.tabStyle
	}

	// Cut the text if the tab length is too much to handle
	if len(text) > 12 {
		text = text[:12] + ""
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		iconStyle.Copy().Foreground(s.iconColors[tabType]).Render(s.icons[tabType]),
		textStyle.Render(text),
	)
}

// Style the status bar cell depending on the of the current tab
func (s style) styleStatusBarCell(tabType tab.Type) string {
	return s.statusBarCell.Copy().
		Background(s.iconColors[tabType]).
		Render(s.texts[tabType])
}
