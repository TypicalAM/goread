package config

import (
	"fmt"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/backend/cache"
	"github.com/TypicalAM/goread/internal/backend/web"
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
func New(backendType string, urlPath string, colors colorscheme.Colorscheme) (Config, error) {
	// Create a new config
	config := Config{}

	// Set the url path
	config.urlPath = urlPath

	// Set the colorscheme
	config.Colors = colors

	// Determine the backend
	var backend backend.Backend
	switch backendType {
	case BackendWeb:
		backend = web.New(config.urlPath)
	case BackendCache:
		var err error
		backend, err = cache.New(config.urlPath)
		if err != nil {
			return Config{}, err
		}
	default:
		return Config{}, fmt.Errorf("invalid backend type: %s", backendType)
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
