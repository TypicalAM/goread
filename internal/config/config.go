package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/backend/cache"
	"github.com/TypicalAM/goread/internal/backend/fake"
	"github.com/TypicalAM/goread/internal/backend/web"
)

// Define the basic backend types
const (
	BackendFake  = "fake"
	BackendWeb   = "web"
	BackendCache = "cache"
)

// Config is the configuration for the program
type Config struct {
	backend backend.Backend
	urlPath string
}

// New returns a new Config
func New(backendType string, urlPath string) (Config, error) {
	// Create a new config
	config := Config{}

	// If the config path is not supplied, use the default
	if urlPath == "" {
		// Get the default config path
		configDir, err := os.UserConfigDir()
		if err != nil {
			return Config{}, err
		}

		// Create the config path
		config.urlPath = filepath.Join(configDir, "goread", "urls.yml")
	}

	// Determine the backend
	var backend backend.Backend
	switch backendType {
	case BackendFake:
		backend = fake.New()
	case BackendWeb:
		backend = web.New()
	case BackendCache:
		backend = cache.New()
	default:
		return Config{}, fmt.Errorf("invalid backend type: %s", backendType)
	}

	// Set the backend
	config.backend = backend

	// Return the config
	return config, nil
}

// Getbackend returns the backend
func (c Config) Getbackend() backend.Backend {
	return c.backend
}

// Close closes the backend
func (c Config) Close() error {
	return c.backend.Close()
}
