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
func (p Overlay) WrapView(view string) string {
	// TODO: Add a padding guardrail and make sure the popup doesn't crash the program when the size is too small
	lines := strings.Split(view, "\n")

	// Overlay the background with the styled text.
	output := make([]string, len(lines))
	for i := 0; i < len(lines); i++ {
		output[i] = p.rowPrefix[i] + lines[i] + p.rowSuffix[i]
	}

	return p.textAbove + "\n" + strings.Join(output, "\n") + "\n" + p.textBelow
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
