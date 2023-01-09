package feed

import (
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/lipgloss"
)

var (
	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorscheme.Global.TextDark)

	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorscheme.Global.Color1)
)
