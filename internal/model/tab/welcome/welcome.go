package welcome

import (
	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/model/popup"
	"github.com/TypicalAM/goread/internal/model/simplelist"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Keymap contains the key bindings for this tab
type Keymap struct {
	SelectCategory key.Binding
	NewCategory    key.Binding
	EditCategory   key.Binding
	DeleteCategory key.Binding
}

// DefaultKeymap contains the default key bindings for this tab
var DefaultKeymap = Keymap{
	SelectCategory: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Open"),
	),
	NewCategory: key.NewBinding(
		key.WithKeys("n", "ctrl+n"),
		key.WithHelp("n/C-n", "New"),
	),
	EditCategory: key.NewBinding(
		key.WithKeys("e", "ctrl+e"),
		key.WithHelp("e/C-e", "Edit"),
	),
	DeleteCategory: key.NewBinding(
		key.WithKeys("d", "ctrl+d"),
		key.WithHelp("d/C-d", "Delete"),
	),
}

// ShortHelp returns the short help for this tab
func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.SelectCategory, k.NewCategory, k.EditCategory, k.DeleteCategory,
	}
}

// FullHelp returns the full help for this tab
func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.SelectCategory, k.NewCategory, k.EditCategory, k.DeleteCategory},
	}
}

// Model contains the state of this tab
type Model struct {
	colors  *colorscheme.Colorscheme
	fetcher backend.Fetcher
	title   string
	keymap  Keymap
	list    simplelist.Model
	width   int
	height  int
	loaded  bool
}

// New creates a new welcome tab with sensible defaults
func New(colors *colorscheme.Colorscheme, width, height int, title string, fetcher backend.Fetcher) Model {
	return Model{
		colors:  colors,
		width:   width,
		height:  height,
		title:   title,
		fetcher: fetcher,
		keymap:  DefaultKeymap,
	}
}

// Title returns the title of the tab
func (m Model) Title() string {
	return m.title
}

// Style returns the style of the tab
func (m Model) Style() tab.Style {
	return tab.Style{
		Color: m.colors.Color4,
		Icon:  "﫢",
		Name:  "WELCOME",
	}
}

// SetSize sets the dimensions of the tab
func (m Model) SetSize(width, height int) tab.Tab {
	if !m.loaded {
		return m
	}

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
	return m.fetcher("")
}

// Update updates the variables of the tab
func (m Model) Update(msg tea.Msg) (tab.Tab, tea.Cmd) {
	// Wait for items to be loaded
	if !m.loaded {
		_, ok := msg.(backend.FetchSuccessMsg)
		if !ok {
			return m, nil
		}

		// Initialize the list of categories, items will be set later
		m.list = simplelist.New(m.colors, "Categories", m.height, true)

		// Add the categories
		m.loaded = true
	}

	switch msg := msg.(type) {
	case backend.FetchSuccessMsg:
		// Update the list of categories
		m.list.SetItems(msg.Items)

	case popup.ChoiceResultMsg:
		if !msg.Result {
			return m, nil
		}

		// Delete the selected category
		delItemName := m.list.SelectedItem().FilterValue()
		itemCount := len(m.list.Items())

		// Move the selection to the next item
		if itemCount == 1 {
			m.list.SetIndex(0)
		} else {
			m.list.SetIndex(m.list.Index() % (itemCount - 1))
		}

		return m, backend.DeleteItem(m, delItemName)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.SelectCategory):
			// Add a new tab with the selected category
			if !m.list.IsEmpty() {
				return m, tab.NewTab(m, m.list.SelectedItem().FilterValue())
			}

			// If the list is empty, return nothing
			return m, nil

		case key.Matches(msg, m.keymap.NewCategory):
			// Add a new category
			return m, backend.NewItem(m, false, make([]string, 2))

		case key.Matches(msg, m.keymap.EditCategory):
			// Edit the selected category
			if !m.list.IsEmpty() {
				item := m.list.SelectedItem().(simplelist.Item)
				fields := []string{item.Title(), item.Description()}
				return m, backend.NewItem(m, true, fields)
			}

		case key.Matches(msg, m.keymap.DeleteCategory):
			if !m.list.IsEmpty() {
				return m, backend.MakeChoice("Delete category?", true)
			}

		default:
			// Check if we need to open a new category
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

// View returns the view for the tab
func (m Model) View() string {
	// Check if the program is loaded, if not, return a loading message
	if !m.loaded {
		return "Loading..."
	}

	// Return the view
	return m.list.View()
}
