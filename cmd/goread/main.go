package main

import (
	"fmt"
	"os"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/model"
	"github.com/TypicalAM/goread/internal/style"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check if the user wants to test the colorscheme
	if len(os.Args) > 1 && os.Args[1] == "colors" {
		fmt.Println(style.BasicColorscheme.TestColors())
		return
	}

	// Create the main model
	model := model.New(backend.FakeBackend{})

	// Start the program
	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
	}
}
