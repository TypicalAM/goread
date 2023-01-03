package main

import (
	"fmt"
	"os"

	"github.com/TypicalAM/goread/internal/config"
	"github.com/TypicalAM/goread/internal/model"
	"github.com/TypicalAM/goread/internal/style"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check if the user wants to test the colorscheme
	if len(os.Args) > 1 && os.Args[1] == "colors" {
		fmt.Println(style.GlobalColorscheme.TestColors())
		return
	}

	cfg, err := config.New(config.BackendCache, "", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create the main model
	model := model.New(cfg.Getbackend())

	// Start the program
	p := tea.NewProgram(model)
	if _, err = p.Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
	}
}
