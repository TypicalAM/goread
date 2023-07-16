package feed

import (
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/lipgloss"
)

// style is the style of the feed tab.
type style struct {
	width         int
	height        int
	listWidth     int
	viewportWidth int

	loadingMsg      lipgloss.Style
	idleList        lipgloss.Style
	focusedList     lipgloss.Style
	idleViewport    lipgloss.Style
	focusedViewport lipgloss.Style

	errIcon string
}

// newStyle creates a new style for the feed tab.
func newStyle(colors colorscheme.Colorscheme, width, height int) style {
	listWidth := width/4 - 2
	viewportWidth := width - listWidth - 4

	loadingMsg := lipgloss.NewStyle().
		MarginLeft(3).
		MarginTop(1)

	errIconStyle := loadingMsg.Copy().
		Foreground(colors.Color4).
		SetString("")

	idleList := lipgloss.NewStyle().
		Width(listWidth).
		Height(height).
		Border(lipgloss.NormalBorder()).
		BorderForeground(colors.TextDark)

	focusedList := idleList.Copy().
		BorderForeground(colors.Color1)

	idleViewport := lipgloss.NewStyle().
		Width(viewportWidth).
		Height(height).
		Border(lipgloss.NormalBorder()).
		BorderForeground(colors.TextDark)

	focusedViewport := idleViewport.Copy().
		BorderForeground(colors.Color1)

	return style{
		width:           width,
		height:          height,
		listWidth:       listWidth,
		viewportWidth:   viewportWidth,
		loadingMsg:      loadingMsg,
		errIcon:         errIconStyle.String(),
		idleList:        idleList,
		focusedList:     focusedList,
		idleViewport:    idleViewport,
		focusedViewport: focusedViewport,
	}
}

// setSize sets the size of the style.
func (s style) setSize(width, height int) style {
	s.width = width
	s.height = height
	s.listWidth = width/4 - 2
	s.viewportWidth = width - s.listWidth - 4
	s.idleList = s.idleList.Width(s.listWidth).Height(height)
	s.focusedList = s.focusedList.Width(s.listWidth).Height(height)
	s.idleViewport = s.idleViewport.Width(s.viewportWidth).Height(height)
	s.focusedViewport = s.focusedViewport.Width(s.viewportWidth).Height(height)
	return s
}
