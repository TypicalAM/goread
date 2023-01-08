package config

import (
	"fmt"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/backend/cache"
	"github.com/TypicalAM/goread/internal/backend/web"
)

// Define the basic backend types
const (
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

	// Set the url path
	config.urlPath = urlPath

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
