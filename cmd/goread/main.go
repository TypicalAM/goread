package main

import (
	"fmt"
	"os"
	"time"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/config"
	"github.com/TypicalAM/goread/internal/model/browser"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// options denote the flags that can be given to the program
type options struct {
	cachePath       string
	colorschemePath string
	urlsPath        string
	getColors       string
	testColors      bool
	resetCache      bool
	cacheSize       int
	cacheDuration   int
}

var (
	opts    = options{}
	rootCmd = &cobra.Command{
		Use:   "goread",
		Short: "goread is a command line tool for reading RSS and ATOM feeds",
		Long:  `goread is a fancy TUI for reading and categorizing different RSS and ATOM feeds`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := Run(); err != nil {
				fmt.Fprintf(os.Stderr, "There has been an error executing the commands: '%s'", err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&opts.cachePath, "cache_path", "c", "", "The path to the cache file")
	rootCmd.Flags().StringVarP(&opts.colorschemePath, "colorscheme_path", "s", "", "The path to the colorscheme file")
	rootCmd.Flags().StringVarP(&opts.urlsPath, "urls_path", "u", "", "The path to the urls file")
	rootCmd.Flags().BoolVarP(&opts.testColors, "test_colors", "t", false, "Test the colorscheme")
	rootCmd.Flags().StringVarP(&opts.getColors, "get_colors", "g", "", "Get the colors from pywal and save them to the colorscheme file")
	rootCmd.Flags().BoolVarP(&opts.resetCache, "reset_cache", "r", false, "Reset the cache")
	rootCmd.Flags().IntVarP(&opts.cacheSize, "cache_size", "z", 0, "The size of the cache")
	rootCmd.Flags().IntVarP(&opts.cacheDuration, "cache_duration", "d", 0, "The duration of the cache in hours")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "There has been an error executing the commands: '%s'", err)
		os.Exit(1)
	}
}

func Run() error {
	// Initialize the colorscheme
	colors := colorscheme.New(opts.colorschemePath)

	// If the user wants to test the colors, do that and exit
	if opts.testColors {
		fmt.Println(colors.TestColors())
		return nil
	}

	// If the user wants to get the colors from pywal, do that and exit
	if opts.getColors != "" {
		// Get the colors from pywal
		err := colors.Convert(opts.getColors)
		if err != nil {
			return err
		}

		// Save the colorscheme
		err = colors.Save()
		if err != nil {
			return err
		}

		// Notify the user
		messageStyle := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#6bae6c"))
		fmt.Println(messageStyle.Render("The new colorscheme was saved to the config directory\n"))
		fmt.Println(colors.TestColors())

		return colors.Convert(opts.getColors)
	}

	// Set the cache size
	if opts.cacheSize > 0 {
		backend.DefaultCacheSize = opts.cacheSize
	}

	// Set the cache duration
	if opts.cacheDuration > 0 {
		backend.DefaultCacheDuration = time.Hour * time.Duration(opts.cacheDuration)
	}

	// Initialize the cfg
	cfg, err := config.New(colors, opts.urlsPath, opts.cachePath, opts.resetCache)
	if err != nil {
		return err
	}

	// Create the browser
	browser := browser.New(cfg)

	// Start the program
	p := tea.NewProgram(browser)
	if _, err = p.Run(); err != nil {
		return err
	}

	// Close the config
	return cfg.Close()
}

func main() {
	Execute()
}
