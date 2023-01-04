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
	title  string
	loaded bool

	loadingSpinner spinner.Model
	list           list.List
	readerFunc     func(string) tea.Cmd

	availableWidth  int
	availableHeight int
}

// New creates a new RssCategoryTab with sensible defaults
func New(availableWidth, availableHeight int, title string, readerFunc func(string) tea.Cmd) RSSCategoryTab {
	// Create a spinner for loading the data
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Foreground(style.GlobalColorscheme.Color1)

	return RSSCategoryTab{
		availableWidth:  availableWidth,
		availableHeight: availableHeight,
		loadingSpinner:  spin,
		title:           title,
		readerFunc:      readerFunc,
	}
}

// Implement the Tab interface
func (c RSSCategoryTab) Title() string {
	return c.title
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
		// The data fetch was successful
		if !c.loaded {
			c.list = list.NewList(c.title, c.availableHeight-5)
			c.loaded = true
		}

		// Set the items of the list
		c.list.SetItems(msg.Items)
		return c, nil

	case tea.KeyMsg:
		// If the tab is not loaded, return
		if !c.loaded {
			return c, nil
		}

		// Handle the list messages
		switch msg.String() {
		case "enter":

			// If it isnt a number, check if it is an enter
			if !c.list.IsEmpty() {
				return c, tab.NewTab(c.list.SelectedItem().FilterValue(), tab.Feed)
			}

		case "n":
			// Add a new category
			return c, backend.NewItem(backend.Feed)

		case "d":
			// Delete the selected category
			if !c.list.IsEmpty() {
				return c, backend.DeleteItem(
					backend.Feed,
					c.list.SelectedItem().FilterValue(),
				)
			}
		}

		// Check if the user opened a tab using the number pad
		if index, ok := c.list.HasItem(msg.String()); ok {
			return c, tab.NewTab(c.list.GetItem(index).FilterValue(), tab.Feed)
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
		return loadingMessage + strings.Repeat("\n", c.availableHeight-3-lipgloss.Height(loadingMessage))
	}

	// Return the list view
	return c.list.View()
}

func (c RSSCategoryTab) Type() tab.Type {
	return tab.Category
}
