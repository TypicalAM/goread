package goread

import (
	"fmt"
	"os"
	"time"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/colorscheme"
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
	dumpColors      bool
	testColors      bool
	resetCache      bool
	cacheSize       int
	cacheDuration   int
}

var (
	msgStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#6bae6c"))
	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e06c75"))

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
	rootCmd.Flags().StringVarP(&opts.cachePath, "cache_path", "", "", "The path to the cache file")
	rootCmd.Flags().StringVarP(&opts.colorschemePath, "colorscheme_path", "c", "", "The path to the colorscheme file")
	rootCmd.Flags().StringVarP(&opts.urlsPath, "urls_path", "u", "", "The path to the urls file")
	rootCmd.Flags().BoolVarP(&opts.testColors, "test_colors", "", false, "Test the colorscheme")
	rootCmd.Flags().BoolVarP(&opts.dumpColors, "dump_colors", "", false, "Dump the colors to the colorscheme file")
	rootCmd.Flags().StringVarP(&opts.getColors, "get_colors", "", "", "Get the colors from pywal and save them to the colorscheme file")
	rootCmd.Flags().BoolVarP(&opts.resetCache, "reset_cache", "", false, "Reset the cache")
	rootCmd.Flags().IntVarP(&opts.cacheSize, "cache_size", "", 0, "The size of the cache")
	rootCmd.Flags().IntVarP(&opts.cacheDuration, "cache_duration", "", 0, "The duration of the cache in hours")
}

// Execute executes the commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "There has been an error executing the commands: '%s'", err)
		os.Exit(1)
	}
}

// Run runs the program
func Run() error {
	colors, err := colorscheme.New(opts.colorschemePath)
	if err != nil {
		return err
	}

	_ = colors.Load()

	// Pretty printing colors
	if opts.testColors {
		fmt.Println(colors.PrettyPrint())
		return nil
	}

	// Dumping colors to file
	if opts.dumpColors {
		if err := colors.Save(); err != nil {
			fmt.Println(errStyle.Render("Failed to save the colorscheme"))
			return err
		}

		fmt.Println(msgStyle.Render(fmt.Sprint("The colorscheme was saved to ", colors.FilePath)))
		return nil
	}

	// Get colors from pywal
	if opts.getColors != "" {
		if err := colors.Convert(opts.getColors); err != nil {
			return err
		}

		if err = colors.Save(); err != nil {
			return err
		}

		fmt.Println(msgStyle.Render(fmt.Sprint("The colorscheme was saved to ", colors.FilePath)))
		fmt.Println(colors.PrettyPrint())
		return nil
	}

	// Set the cache size
	if opts.cacheSize > 0 {
		backend.DefaultCacheSize = opts.cacheSize
	}

	// Set the cache duration
	if opts.cacheDuration > 0 {
		backend.DefaultCacheDuration = time.Hour * time.Duration(opts.cacheDuration)
	}

	// Initialize the backend
	backend, err := backend.New(opts.urlsPath, opts.cachePath, opts.resetCache)
	if err != nil {
		return err
	}

	// Create the browser
	browser := browser.New(colors, backend)

	// Start the program
	if _, err = tea.NewProgram(browser).Run(); err != nil {
		return err
	}

	// Clean up the backend
	return backend.Close()
}
