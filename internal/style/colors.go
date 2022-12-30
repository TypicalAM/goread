package style

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Create the basic colorscheme on startup
var BasicColorscheme = NewDefaultColorscheme()

// The basic colorscheme to use in the app
type Colorscheme struct {
	// Background color
	BgDark   lipgloss.Color
	BgDarker lipgloss.Color

	// Text colors
	Text     lipgloss.Color
	TextDark lipgloss.Color

	// Accent colors
	Color1 lipgloss.Color
	Color2 lipgloss.Color
	Color3 lipgloss.Color
	Color4 lipgloss.Color
	Color5 lipgloss.Color
	Color6 lipgloss.Color
	Color7 lipgloss.Color
}

// A function which returns a new default colorscheme
func NewDefaultColorscheme() Colorscheme {
	return Colorscheme{
		BgDark:   lipgloss.Color("#161622"),
		BgDarker: lipgloss.Color("#11111a"),

		Text:     lipgloss.Color("#FFFFFF"),
		TextDark: lipgloss.Color("#47485b"),

		Color1: lipgloss.Color("#c29fec"),
		Color2: lipgloss.Color("#ddbec0"),
		Color3: lipgloss.Color("#89b4fa"),
		Color4: lipgloss.Color("#e06c75"),
		Color5: lipgloss.Color("#98c379"),
		Color6: lipgloss.Color("#fab387"),
		Color7: lipgloss.Color("#f1c1e4"),
	}
}

// Create a function which displays formatted text in every color available
func (c Colorscheme) TestColors() string {
	// Create the result variable
	result := []string{"A table of all the colors:"}

	// Loop through the colors and add the text to the result
	v := reflect.ValueOf(c)
	typeOfV := v.Type()
	for i := 0; i < v.NumField(); i++ {
		result = append(
			result,
			fmt.Sprintf(
				"%s: %s%s\t%s",
				typeOfV.Field(i).Name,
				strings.Repeat(" ", 15-len(typeOfV.Field(i).Name)),
				lipgloss.NewStyle().
					Foreground(v.Field(i).Interface().(lipgloss.Color)).
					Render("Foreground"),
				lipgloss.NewStyle().
					Background(v.Field(i).Interface().(lipgloss.Color)).
					Render("Background"),
			),
		)
	}

	return strings.Join(result, "\n")
}
