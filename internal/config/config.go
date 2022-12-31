package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/backend/fake"
	"github.com/TypicalAM/goread/internal/backend/web"
)

// Define the basic backend types
const (
	BackendFake = "fake"
	BackendWeb  = "web"
)

// Config is the configuration for the program
type Config struct {
	backend    backend.Backend
	configPath string
	urlPath    string
}

// NewConfig returns a new Config
func New(backend string, configPath string, urlPath string) (Config, error) {
	// If the config path is not supplied, use the default
	if configPath == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return Config{}, err
		}

		configPath = filepath.Join(configDir, "goread/config.json")
	}

	// If the url path is not supplied, use the default
	if urlPath == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return Config{}, err
		}

		urlPath = filepath.Join(configDir, "goread/urls.json")
	}

	// Detect the backend
	switch backend {
	case BackendFake:
		// Return the fake backend
		return Config{
			backend:    fake.New(),
			configPath: configPath,
			urlPath:    urlPath,
		}, nil

	case BackendWeb:
		// Return the web backend
		return Config{
			backend:    web.New(),
			configPath: configPath,
			urlPath:    urlPath,
		}, nil

	default:
		// No backend was found
		return Config{}, fmt.Errorf("Unknown backend: %s", backend)
	}
}

// Getbackend returns the backend
func (c Config) Getbackend() backend.Backend {
	return c.backend
}
