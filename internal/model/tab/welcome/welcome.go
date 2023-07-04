package welcome

import (
	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/model/simplelist"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model contains the state of this tab
type Model struct {
	colors colorscheme.Colorscheme
	width  int
	height int
	title  string
	loaded bool
	list   simplelist.Model
	keymap keymap
	help   help.Model

	// reader is a function which returns a tea.Cmd which will be executed
	// when the tab is initialized
	reader func() tea.Cmd
}

// keymap contains the key bindings for this tab
type keymap struct {
	CloseTab  key.Binding
	CycleTabs key.Binding
	Enter     key.Binding
	New       key.Binding
	Edit      key.Binding
	Delete    key.Binding
}

// ShortHelp returns the short help for this tab
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.CloseTab, k.CycleTabs, k.Enter, k.New, k.Edit, k.Delete,
	}
}

// FullHelp returns the full help for this tab
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CloseTab, k.CycleTabs, k.Enter, k.New, k.Edit, k.Delete},
	}
}

// New creates a new welcome tab with sensible defaults
func New(colors colorscheme.Colorscheme, width, height int, title string, reader func() tea.Cmd) Model {
	help := help.New()
	help.Styles.ShortDesc = lipgloss.NewStyle().Foreground(colors.Text)
	help.Styles.ShortKey = lipgloss.NewStyle().Foreground(colors.Text)
	help.Styles.Ellipsis = lipgloss.NewStyle().Foreground(colors.BgDark)

	return Model{
		colors: colors,
		width:  width,
		height: height,
		title:  title,
		reader: reader,
		help:   help,
		keymap: keymap{
			CloseTab: key.NewBinding(
				key.WithKeys("c/ctrl+w"),
				key.WithHelp("c/ctrl+w", "Close tab"),
			),
			CycleTabs: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "Cycle tabs"),
			),
			Enter: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "Open"),
			),
			New: key.NewBinding(
				key.WithKeys("n/ctrl+n"),
				key.WithHelp("n/ctrl+n", "New"),
			),
			Edit: key.NewBinding(
				key.WithKeys("e/ctrl+e"),
				key.WithHelp("e/ctrl+e", "Edit"),
			),
			Delete: key.NewBinding(
				key.WithKeys("d/ctrl+d"),
				key.WithHelp("d/ctrl+d", "Delete"),
			),
		},
	}
}

// Title returns the title of the tab
func (m Model) Title() string {
	return m.title
}

// Type returns the type of the tab
func (m Model) Type() tab.Type {
	return tab.Welcome
}

// SetSize sets the dimensions of the tab
func (m Model) SetSize(width, height int) tab.Tab {
	m.width = width
	m.height = height
	m.list.SetHeight(m.height)
	return m
}

// ShowHelp shows the help for this tab
func (m Model) ShowHelp() string {
	return m.help.View(m.keymap)
}

// Init initializes the tab
func (m Model) Init() tea.Cmd {
	return m.reader()
}

// Update updates the variables of the tab
func (m Model) Update(msg tea.Msg) (tab.Tab, tea.Cmd) {
	// Wait for items to be loaded
	if !m.loaded {
		_, ok := msg.(backend.FetchSuccessMessage)
		if !ok {
			return m, nil
		}

		// Initialize the list of categories, items will be set later
		m.list = simplelist.New(m.colors, "Categories", m.height, true)

		// Add the categories
		m.loaded = true
	}

	switch msg := msg.(type) {
	case backend.FetchSuccessMessage:
		// Update the list of categories
		m.list.SetItems(msg.Items)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Add a new tab with the selected category
			if !m.list.IsEmpty() {
				return m, tab.NewTab(m.list.SelectedItem().FilterValue(), tab.Category)
			}

			// If the list is empty, return nothing
			return m, nil

		case "n", "ctrl+n":
			// Add a new category
			return m, backend.NewItem(backend.Category, true, nil, nil)

		case "e", "ctrl+e":
			// Edit the selected category
			if !m.list.IsEmpty() {
				categoryPath := []string{m.list.SelectedItem().FilterValue()}
				item := m.list.SelectedItem().(simplelist.Item)
				return m, backend.NewItem(
					backend.Category,
					false,
					categoryPath,
					[]string{item.Title(), item.Description()},
				)
			}

		case "d", "ctrl+d":
			// Delete the selected category
			if !m.list.IsEmpty() {
				return m, backend.DeleteItem(backend.Category, m.list.SelectedItem().FilterValue())
			}

		default:
			// Check if we need to open a new category
			if item, ok := m.list.GetItem(msg.String()); ok {
				return m, tab.NewTab(item.FilterValue(), tab.Category)
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
