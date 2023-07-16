package category

import (
	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/model/simplelist"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Keymap contains the key bindings for this tab
type Keymap struct {
	SelectFeed key.Binding
	NewFeed    key.Binding
	EditFeed   key.Binding
	DeleteFeed key.Binding
}

// DefaultKeymap contains the default key bindings for this tab
var DefaultKeymap = Keymap{
	SelectFeed: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Open"),
	),
	NewFeed: key.NewBinding(
		key.WithKeys("n", "ctrl+n"),
		key.WithHelp("n/C-n", "New"),
	),
	EditFeed: key.NewBinding(
		key.WithKeys("e", "ctrl+e"),
		key.WithHelp("e/C-e", "Edit"),
	),
	DeleteFeed: key.NewBinding(
		key.WithKeys("d", "ctrl+d"),
		key.WithHelp("d/C-d", "Delete"),
	),
}

// ShortHelp returns the short help for this tab
func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.SelectFeed, k.NewFeed, k.EditFeed, k.DeleteFeed,
	}
}

// FullHelp returns the full help for this tab
func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SelectFeed, k.NewFeed, k.EditFeed, k.DeleteFeed},
	}
}

// Model contains the state of this tab
type Model struct {
	colors colorscheme.Colorscheme
	width  int
	height int
	title  string
	loaded bool
	list   simplelist.Model
	keymap Keymap

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
		keymap: DefaultKeymap,
	}
}

// Title returns the title of the tab
func (m Model) Title() string {
	return m.title
}

// Style returns the style of the tab
func (m Model) Style() tab.Style {
	return tab.Style{
		Color: m.colors.Color5,
		Icon:  "﫜",
		Name:  "CATEGORY",
	}
}

// SetSize sets the dimensions of the tab
func (m Model) SetSize(width, height int) tab.Tab {
	m.width = width
	m.height = height
	m.list.SetHeight(m.height)
	return m
}

// GetKeyBinds returns the key bindings of the tab
func (m Model) GetKeyBinds() []key.Binding {
	return m.keymap.ShortHelp()
}

// Init initializes the tab
func (m Model) Init() tea.Cmd {
	return m.reader(m.title)
}

// Update updates the variables of the tab
func (m Model) Update(msg tea.Msg) (tab.Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case backend.FetchSuccessMsg:
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
		switch {
		case key.Matches(msg, m.keymap.SelectFeed):
			if !m.list.IsEmpty() {
				return m, tab.NewTab(m, m.list.SelectedItem().FilterValue())
			}

			// If the list is empty, return nothing
			return m, nil

		case key.Matches(msg, m.keymap.NewFeed):
			// Add a new feed
			return m, backend.NewItem(m, false, make([]string, 2))

		case key.Matches(msg, m.keymap.EditFeed):
			// If the list is empty, return nothing
			if m.list.IsEmpty() {
				return m, nil
			}

			// Edit the selected feed
			item := m.list.SelectedItem().(simplelist.Item)
			fields := []string{item.Title(), item.Description()}
			return m, backend.NewItem(m, true, fields)

		case key.Matches(msg, m.keymap.DeleteFeed):
			if !m.list.IsEmpty() {
				// Delete the selected feed
				delItemName := m.list.SelectedItem().FilterValue()
				itemCount := len(m.list.Items())

				// Move the selection to the next item
				if itemCount == 1 {
					m.list.SetIndex(0)
				} else {
					m.list.SetIndex(m.list.Index() % (itemCount - 1))
				}

				return m, backend.DeleteItem(m, delItemName)
			}

		default:
			// Check if we need to open a new feed
			if item, ok := m.list.GetItem(msg.String()); ok {
				return m, tab.NewTab(m, item.FilterValue())
			}
		}
	}

	// Update the list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View returns the view of the tab
func (m Model) View() string {
	// Check if the program is loaded, if not, return a loading message
	if !m.loaded {
		return "Loading..."
	}

	// Return the list view
	return m.list.View()
}
