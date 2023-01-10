package config

import (
	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/colorscheme"
)

// Define the basic backend types
const (
	BackendWeb   = "web"
	BackendCache = "cache"
)

// Config is the configuration for the program
type Config struct {
	Colors  colorscheme.Colorscheme
	Backend backend.Backend
	urlPath string
}

// New returns a new Config
func New(colors colorscheme.Colorscheme, urlPath, cachePath string, resetCache bool) (Config, error) {
	// Create a new config
	config := Config{}

	// Set the url path
	config.urlPath = urlPath

	// Set the colorscheme
	config.Colors = colors

	// Get the backend
	backend, err := backend.New(urlPath, cachePath, resetCache)
	if err != nil {
		return config, err
	}

	// Set the backend
	config.Backend = backend

	// Return the config
	return config, nil
}

// Close closes the backend
func (c Config) Close() error {
	return c.Backend.Close()
}
