package popup

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TitleBorder creates a fancy border style to wrap a popup
type TitleBorder struct {
	bottomBorder lipgloss.Style
	topBorder    string

	text       string
	borderType lipgloss.Border
	color      lipgloss.Color
}

// NewTitleBorder creates a new fancy border with a specific title
func NewTitleBorder(text string, width, height int, color lipgloss.Color, border lipgloss.Border) TitleBorder {
	tb := TitleBorder{
		text:       strings.Clone(text),
		borderType: border,
		color:      color,
	}

	tb.Resize(width, height)
	return tb
}

// Resize allows resizing the border and adjusting the top border
func (tb *TitleBorder) Resize(width, height int) {
	tb.bottomBorder = lipgloss.NewStyle().
		Width(width-2).
		Height(height-2).
		Border(lipgloss.NormalBorder(), false, true, true, true).
		BorderForeground(tb.color)

	textCopy := " " + strings.Clone(tb.text) + " "
	if width-len(textCopy) < 2 {
		textCopy = textCopy[:width-2]
	}

	textWidth := len(textCopy)
	if textWidth%2 == 1 {
		textCopy += " "
		textWidth++
	}

	fill := (width - 2 - textWidth) / 2

	title := strings.Builder{}
	title.WriteString(tb.borderType.TopLeft)
	title.WriteString(strings.Repeat(tb.borderType.Top, fill+width%2))
	title.WriteString(textCopy)
	title.WriteString(strings.Repeat(tb.borderType.Top, fill))
	title.WriteString(tb.borderType.TopRight)
	tb.topBorder = lipgloss.NewStyle().Foreground(tb.color).Render(title.String())
}

// Render renders the contnent with the fancy border
func (tb TitleBorder) Render(view string) string {
	return lipgloss.JoinVertical(lipgloss.Top, tb.topBorder, tb.bottomBorder.Render(view))
}
