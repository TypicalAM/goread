package popup

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/ansi"
)

// Popup is the popup window interface. In can be implemented in other packages and use the `Default` popup to overlay the content
// on top of the background.
type Popup interface {
	tea.Model
}

// Default is a default popup window.
type Default struct {
	ogSection []string
	section   []string
	width     int
	height    int
	prefix    string
	suffix    string
	startCol  int
}

// New creates a new default popup window.
func New(bgRaw string, width, height int) Default {
	bg := strings.Split(bgRaw, "\n")
	bgWidth := ansi.PrintableRuneWidth(bg[0])
	bgHeight := len(bg)

	startRow := (bgHeight - height) / 2
	startCol := (bgWidth - width) / 2

	ogSection := make([]string, height)
	section := make([]string, height)

	prefix := strings.Join(bg[:startRow], "\n")
	suffix := strings.Join(bg[startRow+height:], "\n")
	copy(ogSection, bg[startRow:startRow+height])

	return Default{
		ogSection: ogSection,
		section:   section,
		width:     width,
		height:    height,
		prefix:    prefix,
		suffix:    suffix,
		startCol:  startCol,
	}
}

// Overlay overlays the given text on top of the background.
func (p Default) Overlay(text string) string {
	// TODO: Add a padding guardrail
	lines := strings.Split(text, "\n")

	// Overlay the background with the styled text.
	// TODO: Use a string builder
	for i, text := range p.ogSection {
		p.section[i] = text[:findPrintableIndex(text, p.startCol)] +
			lines[i] +
			text[findPrintableIndex(text, p.startCol+p.width):]
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		p.prefix,
		strings.Join(p.section, "\n"),
		p.suffix,
	)
}

// Width returns the width of the popup window.
func (p Default) Width() int {
	return p.width
}

// Height returns the height of the popup window.
func (p Default) Height() int {
	return p.height
}

// findPrintableIndex finds the index of the last printable rune at the given index.
func findPrintableIndex(str string, index int) int {
	for i := len(str) - 1; i >= 0; i-- {
		if ansi.PrintableRuneWidth(str[:i]) == index {
			return i
		}
	}

	return -1
}
