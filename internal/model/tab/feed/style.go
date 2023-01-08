package feed

import (
	"github.com/TypicalAM/goread/internal/style"
	"github.com/charmbracelet/lipgloss"
)

var (
	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(style.GlobalColorscheme.TextDark)

	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(style.GlobalColorscheme.Color1)
)
