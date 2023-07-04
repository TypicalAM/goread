package category

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
	reader func(string) tea.Cmd
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

// New creates a new category tab with sensible defaults
func New(colors colorscheme.Colorscheme, width, height int, title string, reader func(string) tea.Cmd) Model {
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
	return tab.Category
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
	return m.reader(m.title)
}

// Update updates the variables of the tab
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

		case "n", "ctrl+n":
			// Add a new category
			feedPath := []string{m.title, m.list.SelectedItem().FilterValue()}
			return m, backend.NewItem(backend.Feed, true, feedPath, nil)

		case "e", "ctrl+e":
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

		case "d", "ctrl+d":
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

// View returns the view of the tab
func (m Model) View() string {
	// Check if the program is loaded, if not, return a loading message
	if !m.loaded {
		return "Loading..."
	}

	// Return the list view
	return m.list.View()
}
