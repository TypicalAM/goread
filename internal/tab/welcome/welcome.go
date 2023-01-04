package welcome

import (
	"github.com/TypicalAM/goread/internal/backend"
	simpleList "github.com/TypicalAM/goread/internal/list"
	"github.com/TypicalAM/goread/internal/tab"
	tea "github.com/charmbracelet/bubbletea"
)

type Welcome struct {
	// General fields
	title  string
	loaded bool

	// The list of categorie
	list       simpleList.List
	readerFunc func() tea.Cmd

	availableWidth  int
	availableHeight int
}

// New creates a new RssFeedTab with sensible defaults
func New(availableWidth, availableHeight int, title string, readerFunc func() tea.Cmd) Welcome {
	return Welcome{
		availableWidth:  availableWidth,
		availableHeight: availableHeight,
		title:           title,
		readerFunc:      readerFunc,
	}
}

// Return the title of the tab
func (w Welcome) Title() string {
	return w.title
}

// Return the type of the tab
func (w Welcome) Type() tab.Type {
	return tab.Welcome
}

// Implement the bubbletea.Model interface
func (w Welcome) Init() tea.Cmd {
	return w.readerFunc()
}

// Update the variables
func (w Welcome) Update(msg tea.Msg) (tab.Tab, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Wait for items to be loaded
	if !w.loaded {
		if msg, ok := msg.(backend.FetchSuccessMessage); ok {
			// Initialize the list of categories, items will be set later
			w.list = simpleList.NewList("Categories", w.availableHeight-5)

			// Add the categories
			w.list.SetItems(msg.Items)
			w.loaded = true
		} else {
			return w, nil
		}
	}

	// Update the list
	w.list, cmd = w.list.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case backend.FetchSuccessMessage:
		// Update the list of categories
		w.list.SetItems(msg.Items)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Add a new tab with the selected category
			if !w.list.IsEmpty() {
				cmds = append(cmds, tab.NewTab(
					w.list.SelectedItem().FilterValue(),
					tab.Category,
				))
			}

		case "n":
			// Add a new category
			cmds = append(cmds, backend.NewItem(backend.Category))

		case "d":
			// Delete the selected category
			if !w.list.IsEmpty() {
				cmds = append(cmds, backend.DeleteItem(
					backend.Category,
					w.list.SelectedItem().FilterValue(),
				))
			}
		}
	}

	// Check the message type
	if msg, ok := msg.(tea.KeyMsg); ok {
		// Check if we need to open a new tab
		if index, ok := w.list.HasItem(msg.String()); ok {
			cmds = append(cmds, tab.NewTab(w.list.GetItem(index).FilterValue(), tab.Category))
		}
	}

	return w, tea.Batch(cmds...)
}

// View the tab
func (w Welcome) View() string {
	// Check if the program is loaded, if not, return an empty string
	if !w.loaded {
		return "Loading..."
	}

	// Return the view
	return w.list.View()
}
