package popup

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/ansi"
)

// Popup is the popup window interface. In can be implemented in other packages and use the `Default` popup to overlay the content on top of the background.
type Popup interface {
	tea.Model
}

// Default is a default popup window.
type Default struct {
	textAbove string
	textBelow string
	rowPrefix []string
	rowSuffix []string
	width     int
	height    int
}

// New creates a new default popup window.
func New(bgRaw string, width, height int) Default {
	bg := strings.Split(bgRaw, "\n")
	bgWidth := ansi.PrintableRuneWidth(bg[0])
	bgHeight := len(bg)

	startRow := (bgHeight - height) / 2
	startCol := (bgWidth - width) / 2

	rowPrefix := make([]string, height)
	rowSuffix := make([]string, height)

	for i, text := range bg[startRow : startRow+height] {
		rowPrefix[i] = text[:findPrintableIndex(text, startCol)]
		rowSuffix[i] = text[findPrintableIndex(text, startCol+width):]
	}

	prefix := strings.Join(bg[:startRow], "\n")
	suffix := strings.Join(bg[startRow+height:], "\n")

	return Default{
		rowPrefix: rowPrefix,
		rowSuffix: rowSuffix,
		width:     width,
		height:    height,
		textAbove: prefix,
		textBelow: suffix,
	}
}

// Overlay overlays the given text on top of the background.
func (p Default) Overlay(text string) string {
	// TODO: Add a padding guardrail
	lines := strings.Split(text, "\n")

	// Overlay the background with the styled text.
	output := make([]string, len(lines))
	for i := 0; i < len(lines); i++ {
		output[i] = p.rowPrefix[i] + lines[i] + p.rowSuffix[i]
	}

	return p.textAbove + "\n" + strings.Join(output, "\n") + "\n" + p.textBelow
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
