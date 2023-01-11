package category

import (
	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/model/simplelist"
	"github.com/TypicalAM/goread/internal/model/tab"
	tea "github.com/charmbracelet/bubbletea"
)

// Model contains the state of this tab
type Model struct {
	colors colorscheme.Colorscheme
	width  int
	height int
	title  string
	loaded bool
	list   simplelist.Model

	// reader is a function which returns a tea.Cmd which will be executed
	// when the tab is initialized
	reader func(string) tea.Cmd
}

// New creates a new category tab with sensible defaults
func New(colors colorscheme.Colorscheme, width, height int, title string, reader func(string) tea.Cmd) Model {
	return Model{
		colors: colors,
		width:  width,
		height: height,
		title:  title,
		reader: reader,
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

// Help returns the help for the tab
func (m Model) Help() tab.Help {
	return tab.Help{
		tab.KeyBind{Key: "enter", Description: "Open"},
		tab.KeyBind{Key: "ctrl+n", Description: "New"},
		tab.KeyBind{Key: "ctrl+e", Description: "Edit"},
		tab.KeyBind{Key: "ctrl+d", Description: "Delete"},
	}
}

// Init initializes the tab
func (m Model) Init() tea.Cmd {
	return m.reader(m.title)
}

// Update the variables of the tab
func (m Model) Update(msg tea.Msg) (tab.Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case backend.FetchSuccessMessage:
		// The data fetch was successful
		if !m.loaded {
			m.list = simplelist.New(m.colors, m.title, m.height, false)
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

		case "ctrl+n":
			// Add a new category
			return m, backend.NewItem(backend.Feed, true, nil, nil)

		case "ctrl+e":
			// If the list is empty, return nothing
			if m.list.IsEmpty() {
				return m, nil
			}

			// Edit the selected feed
			feedPath := []string{m.title, m.list.SelectedItem().FilterValue()}
			item := m.list.SelectedItem().(simplelist.Item)
			return m, backend.NewItem(
				backend.Feed,
				false,
				feedPath,
				[]string{item.FilterValue(), item.Description()},
			)

		case "ctrl+d":
			// Delete the selected category
			if !m.list.IsEmpty() {
				return m, backend.DeleteItem(backend.Feed, m.list.SelectedItem().FilterValue())
			}

		default:
			// Check if we need to open a new feed
			if item, ok := m.list.GetItem(msg.String()); ok {
				return m, tab.NewTab(item.FilterValue(), tab.Feed)
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
