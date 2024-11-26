package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/backend/cache"
	"github.com/TypicalAM/goread/internal/config"
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/TypicalAM/goread/internal/ui/browser"
)

// options denote the flags that can be given to the program
type options struct {
	cacheDir        string
	colorschemePath string
	urlsPath        string
	configPath      string
	getColors       string
	loadOPMLFrom    string
	exportOPMLTo    string
	cacheSize       int
	cacheDuration   int
	dumpColors      bool
	testColors      bool
	resetCache      bool
	urlsReadOnly    bool
}

var (
	msgStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#6bae6c"))
	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e06c75"))
	opts     = options{}
	rootCmd  = &cobra.Command{
		Use:   "goread",
		Short: "goread - a fancy TUI for reading RSS/Atom feeds",
		Run: func(_ *cobra.Command, _ []string) {
			if err := Run(); err != nil {
				fmt.Fprintln(os.Stderr, errStyle.Render(fmt.Sprint("Encountered an error: ", err)))
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.Flags().
		StringVarP(&opts.cacheDir, "cache_dir", "", "", "The path to the cache directory")
	rootCmd.PersistentFlags().
		StringVarP(&opts.colorschemePath, "colorscheme_path", "c", "", "The path to the colorscheme file")
	rootCmd.PersistentFlags().StringVarP(&opts.urlsPath, "urls_path", "u", "", "The path to the urls file")
	rootCmd.PersistentFlags().StringVarP(&opts.configPath, "config_path", "s", "", "The path to the configuration file")
	rootCmd.Flags().BoolVarP(&opts.testColors, "test_colors", "", false, "Test the colorscheme")
	rootCmd.Flags().
		BoolVarP(&opts.dumpColors, "dump_colors", "", false, "Dump the colors to the colorscheme file")
	rootCmd.Flags().
		StringVarP(&opts.getColors, "get_colors", "", "", "Get the colors from pywal and save them to the colorscheme file")
	rootCmd.Flags().BoolVarP(&opts.resetCache, "reset_cache", "", false, "Reset the cache")
	rootCmd.Flags().IntVarP(&opts.cacheSize, "cache_size", "", 0, "The size of the cache")
	rootCmd.Flags().
		IntVarP(&opts.cacheDuration, "cache_duration", "", 0, "The duration of the cache in hours")
	rootCmd.Flags().
		StringVarP(&opts.loadOPMLFrom, "load_opml", "i", "", "Import the feeds from an OPML file")
	rootCmd.Flags().
		StringVarP(&opts.exportOPMLTo, "export_opml", "e", "", "Export the feeds to an OPML file")
	rootCmd.Flags().
		BoolVarP(&opts.urlsReadOnly, "urls_readonly", "", false, "Feed urls config is read-only, skip saving the feed urls configuration")
}

func Execute() {
	_ = rootCmd.Execute()
}

// SetVersion sets the version of the program
func SetVersion(version string) {
	rootCmd.Version = version
}

// Run runs the program
func Run() error {
	if f, err := tea.LogToFile(filepath.Join(os.TempDir(), "goread.log"), ""); err == nil {
		defer f.Close()
	} else {
		log.Println("Failed to create log file")
		log.SetOutput(io.Discard)
	}

	log.Println("Starting goread")

	colors, err := theme.New(opts.colorschemePath)
	if err != nil {
		return err
	}

	if err = colors.Load(); err != nil {
		log.Println("Failed to load colorscheme: ", err)
	}

	// Pretty printing colors
	if opts.testColors {
		fmt.Println(colors.PrettyPrint())
		return nil
	}

	// Dumping colors to file
	if opts.dumpColors {
		if err := colors.Save(); err != nil {
			log.Println("Failed to save colorscheme: ", err)
			return err
		}

		log.Println("Colorscheme saved to ", colors.FilePath)
		fmt.Println(msgStyle.Render(fmt.Sprint("The colorscheme was saved to ", colors.FilePath)))
		return nil
	}

	// Get colors from pywal
	if opts.getColors != "" {
		if err := colors.Convert(opts.getColors); err != nil {
			log.Println("Failed to convert colorscheme: ", err)
			return err
		}

		if err = colors.Save(); err != nil {
			log.Println("Failed to save colorscheme: ", err)
			return err
		}

		log.Println("Colorscheme saved to ", colors.FilePath)
		fmt.Println(msgStyle.Render(fmt.Sprint("The colorscheme was saved to ", colors.FilePath)))
		fmt.Println(colors.PrettyPrint())
		return nil
	}

	// Set the cache size
	if opts.cacheSize > 0 {
		log.Println("Setting cache size to ", opts.cacheSize)
		cache.DefaultCacheSize = opts.cacheSize
	}

	// Set the cache duration
	if opts.cacheDuration > 0 {
		log.Println("Setting cache duration to ", opts.cacheDuration)
		cache.DefaultCacheDuration = time.Hour * time.Duration(opts.cacheDuration)
	}

	// Get the config
	cfg, err := config.New(opts.configPath)
	if err != nil {
		log.Println("Failed to initialize config: ", err)
		return err
	}

	if err := cfg.Load(); err != nil {
		log.Println("Failed to load config: ", err)
		return err
	}

	// Initialize the backend
	backend, err := backend.New(opts.urlsPath, opts.cacheDir, opts.resetCache)
	if err != nil {
		log.Println("Failed to initialize backend: ", err)
		return err
	}

	// Load the OPML file
	if opts.loadOPMLFrom != "" {
		log.Println("Loading OPML file: ", opts.loadOPMLFrom)

		if err := backend.Rss.LoadOPML(opts.loadOPMLFrom); err != nil {
			return err
		}

		fmt.Println(msgStyle.Render("Loaded OPML file successfully"))
		return backend.Close(opts.urlsReadOnly)
	}

	// Export the OPML file
	if opts.exportOPMLTo != "" {
		log.Println("Exporting OPML file to: ", opts.exportOPMLTo)

		if err := backend.Rss.ExportOPML(opts.exportOPMLTo); err != nil {
			return err
		}

		fmt.Println(msgStyle.Render("Exported OPML file successfully"))
		return backend.Close(opts.urlsReadOnly)
	}

	// Create the browser
	browser := browser.New(colors, backend)
	if _, err = tea.NewProgram(browser).Run(); err != nil {
		log.Println("Bubbletea program fail: ", err)
		return err
	}

	// Clean up the backend
	log.Println("Closing backend")
	return backend.Close(opts.urlsReadOnly)
}
