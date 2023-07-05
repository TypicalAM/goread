package keybind

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/TypicalAM/goread/internal/model/browser"
	"github.com/TypicalAM/goread/internal/model/tab/category"
	"github.com/TypicalAM/goread/internal/model/tab/feed"
	"github.com/TypicalAM/goread/internal/model/tab/welcome"
)

// CustomKeymap is the custom keymap for the app.
type CustomKeymap struct {
	path           string
	BrowserMap     browser.Keymap  `json:"browser_map"`
	WelcomeTabMap  welcome.Keymap  `json:"welcome_tab_map"`
	CategoryTabMap category.Keymap `json:"category_tab_map"`
	FeedTabMap     feed.Keymap     `json:"feed_tab_map"`
}

// Populate will try to populate the keymap from a file
func Populate(path string) error {
	// Create a new keymap
	keymap := CustomKeymap{path: path}

	// Check if we can load from a file
	return keymap.load()
}

// Save saves the keymap to a JSON file
func (km CustomKeymap) Save() error {
	// Try to marshall the data
	yamlData, err := json.Marshal(km)
	if err != nil {
		return err
	}

	// Check if the path exists
	if km.path == "" {
		// Get the default path
		defaultPath, err := getDefaultPath()
		if err != nil {
			return err
		}

		// Set the path
		km.path = defaultPath
	}

	// Try to write the data to the file
	if err = os.WriteFile(km.path, yamlData, 0600); err != nil {
		// Try to create the directory
		err = os.MkdirAll(filepath.Dir(km.path), 0755)
		if err != nil {
			return err
		}

		// Try to write to the file again
		err = os.WriteFile(km.path, yamlData, 0600)
		if err != nil {
			return err
		}
	}

	// Successfully wrote the file
	return nil
}

// Load loads a keymap from a JSON file
func (km *CustomKeymap) load() error {
	// Check if the path is valid
	if km.path == "" {
		// Get the default path
		defaultPath, err := getDefaultPath()
		if err != nil {
			return err
		}

		// Set the path
		km.path = defaultPath
	}

	// Try to open the file
	fileContent, err := os.ReadFile(km.path)
	if err != nil {
		return err
	}

	// Try to decode the file
	err = json.Unmarshal(fileContent, &km)
	if err != nil {
		return err
	}

	// Successfully loaded the file
	return nil
}

// getDefaultPath returns the default path for the keybind file
func getDefaultPath() (string, error) {
	// Get the default config path
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// Create the config path
	return filepath.Join(configDir, "goread", "keybinds.json"), nil
}
