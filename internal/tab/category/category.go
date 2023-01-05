package category

import (
	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/list"
	"github.com/TypicalAM/goread/internal/tab"
	tea "github.com/charmbracelet/bubbletea"
)

// Model contains the state of this tab
type Model struct {
	title           string
	loaded          bool
	list            list.List
	availableWidth  int
	availableHeight int

	// readerFunc is a function which returns a tea.Cmd which will be executed
	// when the tab is initialized
	readerFunc func(string) tea.Cmd
}

// New creates a new Category with sensible defaults
func New(availableWidth, availableHeight int, title string, readerFunc func(string) tea.Cmd) Model {
	return Model{
		availableWidth:  availableWidth,
		availableHeight: availableHeight,
		title:           title,
		readerFunc:      readerFunc,
	}
}

// Title returns the title of the tab
func (m Model) Title() string {
	return m.title
}

// Type returns the type of the tab
func (m Model) Type() tab.Type {
	return tab.Category
}

// Init initializes the tab
func (m Model) Init() tea.Cmd {
	return m.readerFunc(m.title)
}

// Update the variables of the tab
func (m Model) Update(msg tea.Msg) (tab.Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case backend.FetchSuccessMessage:
		// The data fetch was successful
		if !m.loaded {
			m.list = list.NewList(m.title, m.availableHeight-5)
			m.loaded = true
		}

		// Set the items of the list
		m.list.SetItems(msg.Items)
		return m, nil

	case tea.KeyMsg:
		// If the tab is not loaded, return
		if !m.loaded {
			return m, nil
		}

		// Handle the key message
		switch msg.String() {
		case "enter":
			if !m.list.IsEmpty() {
				return m, tab.NewTab(m.list.SelectedItem().FilterValue(), tab.Feed)
			}

			// If the list is empty, return nothing
			return m, nil

		case "n":
			// Add a new category
			return m, backend.NewItem(backend.Feed)

		case "d":
			// Delete the selected category
			if !m.list.IsEmpty() {
				return m, backend.DeleteItem(backend.Feed, m.list.SelectedItem().FilterValue())
			}

		default:
			// Check if we need to open a new feed
			if index, ok := m.list.HasItem(msg.String()); ok {
				return m, tab.NewTab(m.list.GetItem(index).FilterValue(), tab.Feed)
			}
		}
	}

	// Update the list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View the tab
func (m Model) View() string {
	// Check if the program is loaded, if not, return a loading message
	if !m.loaded {
		return "Loading..."
	}

	// Return the list view
	return m.list.View()
}
