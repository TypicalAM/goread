package feed

import (
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// style is the style of the feed tab.
type style struct {
	listItems       list.DefaultItemStyles
	link            lipgloss.Style
	loadingMsg      lipgloss.Style
	idleList        lipgloss.Style
	focusedList     lipgloss.Style
	idleViewport    lipgloss.Style
	focusedViewport lipgloss.Style
	errIcon         string
	width           int
	height          int
	listWidth       int
	viewportWidth   int
}

// newStyle creates a new style for the feed tab.
func newStyle(colors *colorscheme.Colorscheme, width, height int) style {
	listWidth := width/4 - 2
	viewportWidth := width - listWidth - 4

	link := lipgloss.NewStyle().
		Background(colors.Color1).
		Underline(true)

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

	// Create the styles for the list items
	delegateStyles := list.NewDefaultItemStyles()
	delegateStyles.SelectedTitle = delegateStyles.SelectedTitle.Copy().
		BorderForeground(colors.Color3).
		Foreground(colors.Color3).
		Italic(true)

	delegateStyles.SelectedDesc = delegateStyles.SelectedDesc.Copy().
		BorderForeground(colors.Color3).
		Foreground(colors.Color2).
		Height(2).
		Italic(true)

	delegateStyles.NormalDesc = delegateStyles.NormalDesc.Copy().
		Foreground(colors.TextDark).
		Height(2)

	return style{
		width:           width,
		height:          height,
		listWidth:       listWidth,
		viewportWidth:   viewportWidth,
		link:            link,
		loadingMsg:      loadingMsg,
		errIcon:         errIconStyle.String(),
		idleList:        idleList,
		focusedList:     focusedList,
		idleViewport:    idleViewport,
		focusedViewport: focusedViewport,
		listItems:       delegateStyles,
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

// popupStyle is the style of the popup window.
type popupStyle struct {
	general   lipgloss.Style
	heading   lipgloss.Style
	item      lipgloss.Style
	itemTitle lipgloss.Style
	itemField lipgloss.Style
}

// newPopupStyle creates a new popup style.
func newPopupStyle(colors *colorscheme.Colorscheme, width, height int) popupStyle {
	general := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Width(width - 2).
		Height(height - 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(colors.Color1)

	heading := lipgloss.NewStyle().
		Margin(1, 0, 1, 0).
		Width(width - 2).
		Align(lipgloss.Center).
		Italic(true)

	item := lipgloss.NewStyle().
		Margin(0, 4).
		PaddingLeft(1).
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(colors.Color3).
		Italic(true)

	itemTitle := lipgloss.NewStyle().
		Foreground(colors.Color3)

	itemField := lipgloss.NewStyle().
		Foreground(colors.Color2)

	return popupStyle{
		general:   general,
		heading:   heading,
		item:      item,
		itemTitle: itemTitle,
		itemField: itemField,
	}
}
