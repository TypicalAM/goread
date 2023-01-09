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
func parseCmdLine() (urlPath string, testColors, pywalConvert bool) {
	// Create the flagset
	urlPathPtr := flag.String("url_path", "", "The path to the url file")
	testColorsPtr := flag.Bool("colors", false, "Test the colorscheme")
	pywalConvertPtr := flag.Bool("pywal", false, "Convert a pywal colorscheme to a goread colorschemea and save it to the config directory")

	// Parse the flags
	flag.Parse()

	urlPath = *urlPathPtr
	testColors = *testColorsPtr
	pywalConvert = *pywalConvertPtr

	// Return the default path
	return urlPath, testColors, pywalConvert
}

func main() {
	// Parse the command line arguments
	urlPath, testColors, pywalConvert := parseCmdLine()

	// TODO: configurable
	colors := colorscheme.New("")

	// If the user wants to convert a pywal colorscheme to a goread colorscheme
	if pywalConvert {
		convertFromPywal(colors, "")
		os.Exit(0)
	}

	// If the user wants to test the colors, do that and exit
	if testColors {
		fmt.Println(colors.TestColors())
		os.Exit(0)
	}

	// Create the config
	cfg, err := config.New(urlPath, colors)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cfg.Close()

	// Create the main model
	model := browser.New(cfg)

	// Start the program
	p := tea.NewProgram(model)
	if _, err = p.Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
	}
}

func convertFromPywal(colors colorscheme.Colorscheme, pywalFilePath string) {
	// Convert the pywal colorscheme
	err := colors.Convert(pywalFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Save the colorscheme
	err = colors.Save()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Notify the user
	messageStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#6bae6c"))
	fmt.Println(messageStyle.Render("The new colorscheme was saved to the config directory\n"))
	fmt.Println(colors.TestColors())
}
