package popup

import (
	"strings"

	"github.com/muesli/ansi"
)

// Overlay allows you to overlay text on top of a background and achieve a popup.
type Overlay struct {
	textAbove string
	textBelow string
	rowPrefix []string
	rowSuffix []string
	width     int
	height    int
}

// NewOverlay creates a new overlay and computes the necessary indices.
func NewOverlay(bgRaw string, width, height int) Overlay {
	bg := strings.Split(bgRaw, "\n")
	bgWidth := ansi.PrintableRuneWidth(bg[0])
	bgHeight := len(bg)

	if height > bgHeight {
		height = bgHeight
	}
	if width > bgWidth {
		width = bgWidth
	}

	startRow := (bgHeight - height) / 2
	startCol := (bgWidth - width) / 2

	rowPrefix := make([]string, height)
	rowSuffix := make([]string, height)

	for i, text := range bg[startRow : startRow+height] {
		popupStart := findPrintIndex(text, startCol)
		popupEnd := findPrintIndex(text, startCol+width)

		if popupStart != -1 {
			rowPrefix[i] = text[:popupStart]
		} else {
			rowPrintable := ansi.PrintableRuneWidth(text)
			rowPrefix[i] = text + strings.Repeat(" ", startCol-rowPrintable)
		}

		if popupEnd != -1 {
			rowSuffix[i] = text[popupEnd:]
		} else {
			rowSuffix[i] = ""
		}
	}

	prefix := strings.Join(bg[:startRow], "\n")
	suffix := strings.Join(bg[startRow+height:], "\n")

	return Overlay{
		rowPrefix: rowPrefix,
		rowSuffix: rowSuffix,
		width:     width,
		height:    height,
		textAbove: prefix,
		textBelow: suffix,
	}
}

// WrapView overlays the given text on top of the background.
// TODO: Maybe handle the box here. It's a bit weird to have to do it in the view.
func (p Overlay) WrapView(view string) string {
	var b strings.Builder
	b.WriteString(p.textAbove)
	b.WriteRune('\n')

	lines := strings.Split(view, "\n")
	for i := 0; i < len(lines) && i < p.height; i++ {
		b.WriteString(p.rowPrefix[i])
		b.WriteString(lines[i])
		b.WriteString(p.rowSuffix[i])
		b.WriteRune('\n')
	}

	b.WriteString(p.textBelow)
	return b.String()
}

// Width returns the width of the popup window.
func (p Overlay) Width() int {
	return p.width
}

// Height returns the height of the popup window.
func (p Overlay) Height() int {
	return p.height
}

// findPrintIndex finds the print index, that is what string index corresponds to the given printable rune index.
func findPrintIndex(str string, index int) int {
	for i := len(str) - 1; i >= 0; i-- {
		if ansi.PrintableRuneWidth(str[:i]) == index {
			return i
		}
	}

	return -1
}
