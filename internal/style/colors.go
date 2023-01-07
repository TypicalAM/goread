package style

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Create the basic colorscheme on startup
var GlobalColorscheme = Load("")

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

// Save saves the colorscheme to a JSON file
func (c Colorscheme) Save(path string) error {
	// Try to marshall the data
	yamlData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// Check if the path exists
	if path == "" {
		// Get the default path
		defaultPath, err := getDefaultPath()
		if err != nil {
			return err
		}

		// Set the path
		path = defaultPath
	}

	// Try to write the data to the file
	if err = os.WriteFile(path, yamlData, 0600); err != nil {
		// Try to create the directory
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}

		// Try to write to the file again
		err = os.WriteFile(path, yamlData, 0600)
		if err != nil {
			return err
		}
	}

	// Successfully wrote the file
	return nil
}

// Load loads a colorscheme from a JSON file
func Load(path string) Colorscheme {
	// Create the colorscheme
	colorscheme := NewDefaultColorscheme()

	// Check if the path is valid
	if path == "" {
		// Get the default path
		defaultPath, err := getDefaultPath()
		if err != nil {
			return colorscheme
		}

		// Set the path
		path = defaultPath
	}

	// Try to open the file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return colorscheme
	}

	// Try to decode the file
	err = json.Unmarshal(fileContent, &colorscheme)
	if err != nil {
		return colorscheme
	}

	// Successfully loaded the file
	return colorscheme
}

// Convert takes the information from a pywal file and converts it to a colorscheme
func (c *Colorscheme) Convert() error {
	// Get the path to the cache dir
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	// Try to open the file
	path := filepath.Join(cacheDir, "wal", "colors.json")
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Try to decode the file
	var walColorscheme map[string]interface{}
	err = json.Unmarshal(fileContent, &walColorscheme)
	if err != nil {
		return err
	}

	// Set the colors
	c.BgDark = lipgloss.Color(walColorscheme["special"].(map[string]interface{})["background"].(string))
	c.BgDarker = lipgloss.Color(walColorscheme["special"].(map[string]interface{})["background"].(string))
	c.Text = lipgloss.Color(walColorscheme["special"].(map[string]interface{})["foreground"].(string))
	c.TextDark = lipgloss.Color(walColorscheme["special"].(map[string]interface{})["foreground"].(string))
	c.Color1 = lipgloss.Color(walColorscheme["colors"].(map[string]interface{})["color1"].(string))
	c.Color2 = lipgloss.Color(walColorscheme["colors"].(map[string]interface{})["color2"].(string))
	c.Color3 = lipgloss.Color(walColorscheme["colors"].(map[string]interface{})["color3"].(string))
	c.Color4 = lipgloss.Color(walColorscheme["colors"].(map[string]interface{})["color4"].(string))
	c.Color5 = lipgloss.Color(walColorscheme["colors"].(map[string]interface{})["color5"].(string))
	c.Color6 = lipgloss.Color(walColorscheme["colors"].(map[string]interface{})["color6"].(string))
	c.Color7 = lipgloss.Color(walColorscheme["colors"].(map[string]interface{})["color7"].(string))

	// Successfully converted the colorscheme
	return nil
}

// TestColors displays the colorscheme in the terminal
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

// getDefaultPath returns the default path for the colorscheme file
func getDefaultPath() (string, error) {
	// Get the default config path
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// Create the config path
	return filepath.Join(configDir, "goread", "colorscheme.json"), nil
}
