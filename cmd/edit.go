package cmd

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/TypicalAM/goread/internal/backend/rss"
	"github.com/TypicalAM/goread/internal/config"
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/spf13/cobra"
)

var (
	exampleFiles embed.FS
	editCmd      = &cobra.Command{
		Use:       "edit [config|colorscheme|urls]",
		Short:     "Edit a goread configuration file",
		ValidArgs: []string{"config", "colorscheme", "urls"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(_ *cobra.Command, args []string) {
			if err := RunEdit(args[0]); err != nil {
				fmt.Fprintln(os.Stderr, errStyle.Render(fmt.Sprint("Encountered an error: ", err)))
				os.Exit(1)
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(editCmd)
}

func SetExampleFiles(fs embed.FS) {
	exampleFiles = fs
}

func RunEdit(filetype string) error {
	path := ""
	var err error
	switch filetype {
	case "config":
		if opts.configPath == "" {
			path, err = config.GetDefaultPath()
			if err != nil {
				return fmt.Errorf("failed to get default path for config: %w", err)
			}
		} else {
			path = opts.configPath
		}

	case "colorscheme":
		if opts.colorschemePath == "" {
			path, err = theme.GetDefaultPath()
			if err != nil {
				return fmt.Errorf("failed to get default path for theme: %w", err)
			}
		} else {
			path = opts.colorschemePath
		}

	case "urls":
		if opts.urlsPath == "" {
			path, err = rss.GetDefaultPath()
			if err != nil {
				return fmt.Errorf("failed to get default path for urls: %w", err)
			}
		} else {
			path = opts.urlsPath
		}
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directories: %w", err)
		}

		exampleFileData := make([]byte, 0)
		switch filetype {
		case "config":
			exampleFileData, err = exampleFiles.ReadFile("internal/test/example/goread.yml")
			if err != nil {
				return fmt.Errorf("failed to get example file data: %w", err)
			}

		case "colorscheme":
			exampleFileData, err = exampleFiles.ReadFile("internal/test/example/colorscheme.json")
			if err != nil {
				return fmt.Errorf("failed to get example file data: %w", err)
			}

		case "urls":
			exampleFileData, err = exampleFiles.ReadFile("internal/test/example/urls.yml")
			if err != nil {
				return fmt.Errorf("failed to get example file data: %w", err)
			}
		}

		if err := os.WriteFile(path, exampleFileData, 0600); err != nil {
			return fmt.Errorf("failed to write data to file: %w", err)
		}
	}

	pa := strings.Split(os.Getenv("EDITOR")+" "+path, " ")
	cmd := exec.Command(pa[0], pa[1:]...) //nolint:gosec
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
