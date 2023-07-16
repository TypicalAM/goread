package colorscheme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Colorscheme is a struct that contains all the colors for the application
type Colorscheme struct {
	FilePath string         `json:"-"`
	BgDark   lipgloss.Color `json:"bg_dark"`
	BgDarker lipgloss.Color `json:"bg_darker"`
	Text     lipgloss.Color `json:"text"`
	TextDark lipgloss.Color `json:"text_dark"`
	Color1   lipgloss.Color `json:"color1"`
	Color2   lipgloss.Color `json:"color2"`
	Color3   lipgloss.Color `json:"color3"`
	Color4   lipgloss.Color `json:"color4"`
	Color5   lipgloss.Color `json:"color5"`
	Color6   lipgloss.Color `json:"color6"`
	Color7   lipgloss.Color `json:"color7"`
}

// New will create a new colorscheme and try to load it
func New(path string) (*Colorscheme, error) {
	if path == "" {
		defaultPath, err := getDefaultPath()
		if err != nil {
			return nil, err
		}

		path = defaultPath
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	colors := Colorscheme{FilePath: path}
	if err = json.Unmarshal(fileContent, &colors); err != nil {
		return nil, err
	}

	return &colors, nil
}

// NewDefault returns a new default colorscheme
func NewDefault() *Colorscheme {
	return &Colorscheme{
		BgDark:   "#161622",
		BgDarker: "#11111a",
		Text:     "#FFFFFF",
		TextDark: "#47485b",
		Color1:   "#c29fec",
		Color2:   "#ddbec0",
		Color3:   "#89b4fa",
		Color4:   "#e06c75",
		Color5:   "#98c379",
		Color6:   "#fab387",
		Color7:   "#f1c1e4",
	}
}

// Save saves the colorscheme to a JSON file
func (c Colorscheme) Save() error {
	// Check if the path exists
	if c.FilePath == "" {
		// Get the default path
		defaultPath, err := getDefaultPath()
		if err != nil {
			return err
		}

		// Set the path
		c.FilePath = defaultPath
	}

	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// Try to write the data to the file
	if err = os.WriteFile(c.FilePath, jsonData, 0600); err != nil {
		// Try to create the directory
		err = os.MkdirAll(filepath.Dir(c.FilePath), 0755)
		if err != nil {
			return err
		}

		// Try to write to the file again
		err = os.WriteFile(c.FilePath, jsonData, 0600)
		if err != nil {
			return err
		}
	}

	// Successfully wrote the file
	return nil
}

// Convert takes the information from a pywal file and converts it to a colorscheme
func (c *Colorscheme) Convert(filePath string) error {
	if filePath == "" {
		// Get the path to the cache dir
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return err
		}

		filePath = filepath.Join(cacheDir, "wal", "colors.json")
	}

	// Try to open the file
	fileContent, err := os.ReadFile(filePath)
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

	// Iterate over every color
	for _, color := range []lipgloss.Color{c.BgDark, c.BgDarker, c.Text, c.TextDark, c.Color1, c.Color2, c.Color3, c.Color4, c.Color5, c.Color6, c.Color7} {
		// Create the foreground
		foreground := lipgloss.NewStyle().Foreground(color)

		// Create the background
		background := lipgloss.NewStyle().Background(color)

		// Add the color to the result
		result = append(result, fmt.Sprintf(
			"%s %s %s",
			foreground.Render("foreground"),
			background.Render("background"),
			color,
		))
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
