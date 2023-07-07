package browser

import (
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/charmbracelet/lipgloss"
)

// style is the internal style of the browser
type style struct {
	colors               colorscheme.Colorscheme
	activeTab            lipgloss.Style
	activeTabIcon        lipgloss.Style
	tab                  lipgloss.Style
	tabIcon              lipgloss.Style
	tabGap               lipgloss.Style
	statusBarGap         lipgloss.Style
	statusBarCell        lipgloss.Style
	offlineStatusBarCell lipgloss.Style
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

	return style{
		colors:               colors,
		activeTab:            activeTab,
		activeTabIcon:        activeTabIcon,
		tab:                  tabStyle,
		tabIcon:              tabIcon,
		tabGap:               tabGap,
		statusBarGap:         statusBarGap,
		statusBarCell:        statusBarCell,
		offlineStatusBarCell: statusBarCell.Copy().Background(colors.TextDark),
	}
}

// attachIcon attaches an icon based on the tab type
func (s style) attachIcon(tabToStyle tab.Tab, title string, active bool) string {
	var iconStyle, textStyle lipgloss.Style
	if active {
		iconStyle, textStyle = s.activeTabIcon, s.activeTab
	} else {
		iconStyle, textStyle = s.tabIcon, s.tab
	}

	// Cut the text if the tab length is too much to handle
	// TODO: Why 12 ???
	if len(title) > 12 {
		title = title[:12] + ""
	}

	tabStyle := tabToStyle.Style()
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		iconStyle.Foreground(tabStyle.Color).Render(tabStyle.Icon),
		textStyle.Render(title),
	)
}

// styleStatusBarCell styles the status bar cell based on the tab type
func (s style) styleStatusBarCell(tabToStyle tab.Tab, offline bool) string {
	tabStyle := tabToStyle.Style()
	if offline {
		return s.offlineStatusBarCell.
			Render(tabStyle.Name)
	}

	return s.statusBarCell.
		Background(tabStyle.Color).
		Render(tabStyle.Name)
}
