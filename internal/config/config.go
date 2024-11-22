package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/TypicalAM/goread/internal/ui/browser"
	"github.com/TypicalAM/goread/internal/ui/simplelist"
	"github.com/TypicalAM/goread/internal/ui/tab/category"
	"github.com/TypicalAM/goread/internal/ui/tab/feed"
	"github.com/TypicalAM/goread/internal/ui/tab/overview"
	"github.com/charmbracelet/bubbles/key"
	"gopkg.in/yaml.v3"
)

var Default = Config{}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type Config struct {
	Keymap map[string]KeymapConfig `yaml:"keymap"`

	filePath string
}

type KeymapConfig map[string]KeyList

type KeyList []string

// New will create a new config structure
func New(path string) (*Config, error) {
	log.Println("Creating new settings")
	if path == "" {
		defaultPath, err := GetDefaultPath()
		if err != nil {
			return nil, fmt.Errorf("cfg.New: %w", err)
		}

		// Set the path
		path = defaultPath
	}

	cfg := Default
	cfg.filePath = path
	return &cfg, nil
}

// Load will try to load the config structure from a file
func (cfg *Config) Load() error {
	log.Println("Loading config from", cfg.filePath)
	data, err := os.ReadFile(cfg.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("cfg.Load: %w", err)
	}

	if err = yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("cfg.Load: %w", err)
	}

	// Apply the config
	allowedKeymaps := []string{"browser", "overview", "category", "feed", "list"}
	for keyCategory, keymap := range cfg.Keymap {
		if !slices.Contains(allowedKeymaps, keyCategory) {
			return fmt.Errorf("cfg.Load: unrecognized keymap: %s", keyCategory)
		}

		var keymapType reflect.Type
		var keymapValue reflect.Value
		switch keyCategory {
		case "browser":
			keymapType = reflect.TypeOf(browser.DefaultKeymap)
			keymapValue = reflect.ValueOf(&browser.DefaultKeymap)
		case "overview":
			keymapType = reflect.TypeOf(overview.DefaultKeymap)
			keymapValue = reflect.ValueOf(&overview.DefaultKeymap)
		case "category":
			keymapType = reflect.TypeOf(category.DefaultKeymap)
			keymapValue = reflect.ValueOf(&category.DefaultKeymap)
		case "feed":
			keymapType = reflect.TypeOf(feed.DefaultKeymap)
			keymapValue = reflect.ValueOf(&feed.DefaultKeymap)
		case "list":
			keymapType = reflect.TypeOf(simplelist.DefaultKeymap)
			keymapValue = reflect.ValueOf(&simplelist.DefaultKeymap)
		}

		fields := reflect.VisibleFields(keymapType)
		snakeToOriginal := make(map[string]string, len(fields))
		origHelpText := make(map[string]string, len(fields))
		for _, field := range fields {
			snakeToOriginal[toSnakeCase(field.Name)] = field.Name
			origHelpText[field.Name] = keymapValue.
				Elem().
				FieldByName(field.Name).
				MethodByName("Help").
				Call([]reflect.Value{})[0].
				FieldByName("Desc").
				String()
		}

		for name, keys := range keymap {
			if _, ok := snakeToOriginal[name]; !ok {
				b := strings.Builder{}
				for name := range snakeToOriginal {
					b.WriteRune(' ')
					b.WriteString(name)
				}
				return fmt.Errorf("cfg.Load: %s doesn't exist on %s, available options are:%s", name, keyCategory, b.String())
			}

			if len(keys) == 0 {
				return fmt.Errorf("cfg.Load: option %s on category %s doesn't have any keys bound", name, keyCategory)
			}

			specialKeys := map[string]string{"up": "↑", "down": "↓", "left": "←", "right": "→"}
			s := strings.Builder{}
			for _, key := range keys {
				if specialKey, ok := specialKeys[key]; ok {
					s.WriteString(specialKey)
				} else {
					s.WriteString(key)
				}
				s.WriteRune('/')
			}
			helpKeyStr := s.String()[:s.Len()-1] // cut off the last '/'

			origName := snakeToOriginal[name]
			newBind := key.NewBinding(key.WithKeys(keys...), key.WithHelp(helpKeyStr, origHelpText[origName]))
			keymapValue.Elem().FieldByName(origName).Set(reflect.ValueOf(newBind))
		}
	}
	return nil
}

// GetDefaultPath will return the default path for the config file
func GetDefaultPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cfg.getDefaultPath: %w", err)
	}

	return filepath.Join(configDir, "goread", "goread.yml"), nil
}

// toSnakeCase transforms the pascalCase string to snake_case
func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
