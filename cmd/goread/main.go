package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/config"
	"github.com/TypicalAM/goread/internal/model/browser"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// parseCmdLine parses the command line arguments
func parseCmdLine() (urlPath, backend string, testColors, pywalConvert bool) {
	// Create the flagset
	backendPtr := flag.String("backend", "cache", "The backend to use for the config file")
	urlPathPtr := flag.String("url_path", "", "The path to the url file")
	testColorsPtr := flag.Bool("colors", false, "Test the colorscheme")
	pywalConvertPtr := flag.Bool("pywal", false, "Convert a pywal colorscheme to a goread colorschemea and save it to the config directory")

	// Parse the flags
	flag.Parse()

	backend = *backendPtr
	urlPath = *urlPathPtr
	testColors = *testColorsPtr
	pywalConvert = *pywalConvertPtr

	// Return the default path
	return urlPath, backend, testColors, pywalConvert
}

func main() {
	// Parse the command line arguments
	urlPath, backend, testColors, pywalConvert := parseCmdLine()

	// If the user wants to convert a pywal colorscheme to a goread colorscheme
	if pywalConvert {
		convertFromPywal()
		os.Exit(0)
	}

	// If the user wants to test the colors, do that and exit
	if testColors {
		fmt.Println(colorscheme.Global.TestColors())
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
	model := browser.New(cfg.Getbackend())

	// Start the program
	p := tea.NewProgram(model)
	if _, err = p.Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
	}
}

func convertFromPywal() {
	// Convert the pywal colorscheme
	// TODO: Make this configurable
	err := colorscheme.Global.Convert("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Save the colorscheme
	err = colorscheme.Global.Save("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Notify the user
	messageStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#6bae6c"))
	fmt.Println(messageStyle.Render("The new colorscheme was saved to the config directory\n"))
	fmt.Println(colorscheme.Global.TestColors())
}
