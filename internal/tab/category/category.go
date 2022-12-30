package category

import (
	"fmt"
	"strings"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/list"
	"github.com/TypicalAM/goread/internal/style"
	"github.com/TypicalAM/goread/internal/tab"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RSSCategoryTab is a tab to choose a category from a list of categories
type RSSCategoryTab struct {
	title       string
	index       int
	loaded      bool
	description string

	loadingSpinner spinner.Model
	list           list.List
	readerFunc     func(string) tea.Cmd
}

// New creates a new RssCategoryTab with sensible defautls
func New(title string, index int, readerFunc func(string) tea.Cmd) RSSCategoryTab {
	// Create a spinner for loading the data
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Foreground(style.BasicColorscheme.Color1)

	return RSSCategoryTab{
		loadingSpinner: spin,
		title:          title,
		index:          index,
		readerFunc:     readerFunc,
	}
}

// Implement the Tab interface
func (c RSSCategoryTab) Title() string {
	return c.title
}

func (c RSSCategoryTab) Index() int {
	return c.index
}

func (c RSSCategoryTab) Loaded() bool {
	return c.loaded
}

func (c RSSCategoryTab) Init() tea.Cmd {
	return tea.Batch(c.readerFunc(c.title), c.loadingSpinner.Tick)
}

func (c RSSCategoryTab) Update(msg tea.Msg) (tab.Tab, tea.Cmd) {
	// Update the list
	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)

	switch msg := msg.(type) {
	case backend.FetchSuccessMessage:
		// The data fetch was successfull
		if !c.loaded && style.WindowWidth != 0 && style.WindowHeight != 0 {
			c.list = list.NewList(c.title, style.WindowHeight-5)
			c.list.SetItems(msg.Items)
			c.loaded = true
			return c, nil
		}
	case tea.KeyMsg:
		// If the tab is not loaded, return
		if !c.loaded {
			return c, nil
		}

		// Check if the user opened a tab using the number pad
		if index, ok := c.list.HasItem(msg.String()); ok {
			return c, tab.NewTab(c.list.GetItem(index).FilterValue(), tab.Feed)
		}

		// If it isnt a number, check if it is an enter
		if msg.String() == "enter" {
			return c, tab.NewTab(c.list.SelectedItem().FilterValue(), tab.Feed)
		}
	default:
		if !c.loaded {
			c.loadingSpinner, cmd = c.loadingSpinner.Update(msg)
			return c, cmd
		}
	}

	// Return the cmd
	return c, cmd
}

// Showcase the list of the feeds that can be "clicked" to
// load & open the corresponding feed
func (c RSSCategoryTab) View() string {
	// If the list is not loaded, return a loading message
	if !c.loaded {
		// Create the loading message
		loadingMessage := lipgloss.NewStyle().
			MarginLeft(3).
			MarginTop(1).
			Render(fmt.Sprintf("%s Loading category %s", c.loadingSpinner.View(), c.title))

		// Display the loading message with padding
		return loadingMessage + strings.Repeat("\n", style.WindowHeight-3-lipgloss.Height(loadingMessage))
	}

	// Return the list view
	return c.list.View()
}

func (c RSSCategoryTab) Type() tab.TabType {
	return tab.Category
}

// Set the index of the tab
func (c RSSCategoryTab) SetIndex(index int) tab.Tab {
	c.index = index
	return c
}
