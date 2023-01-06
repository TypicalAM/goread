package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TypicalAM/goread/internal/config"
	model "github.com/TypicalAM/goread/internal/model/main"
	"github.com/TypicalAM/goread/internal/style"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/sys/unix"
)

// parseCmdLine parses the command line arguments
func parseCmdLine() (urlPath string, backend string, testColors bool, err error) {
	// Create the flagset
	backendPtr := flag.String("backend", "cache", "The backend to use for the config file")
	urlPathPtr := flag.String("url_path", "", "The path to the url file")
	testColorsPtr := flag.Bool("colors", false, "Test the colorscheme")

	// Parse the flags
	flag.Parse()

	backend = *backendPtr
	urlPath = *urlPathPtr
	testColors = *testColorsPtr

	// Check if the backend is valid
	if backend != config.BackendCache && backend != config.BackendWeb {
		return "", "", false, fmt.Errorf("invalid backend: %s", backend)
	}

	// Check if the config path is valid and writeable
	configDir := filepath.Dir(urlPath)
	if unix.Access(configDir, unix.W_OK) != nil {
		return "", "", false, fmt.Errorf("config file directory is not writable: %s", configDir)
	}

	// Return the default path
	return urlPath, backend, testColors, nil
}

func main() {
	// Parse the command line arguments
	urlPath, backend, testColors, err := parseCmdLine()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// If the user wants to test the colors, do that and exit
	if testColors {
		fmt.Println(style.GlobalColorscheme.TestColors())
		os.Exit(0)
	}

	// Create the config
	cfg, err := config.New(backend, urlPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cfg.Close()

	// Create the main model
	model := model.New(cfg.Getbackend())

	// Start the program
	p := tea.NewProgram(model)
	if _, err = p.Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
	}
}
