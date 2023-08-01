package category

import (
	"log"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/TypicalAM/goread/internal/ui/popup"
	"github.com/TypicalAM/goread/internal/ui/simplelist"
	"github.com/TypicalAM/goread/internal/ui/tab"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Model contains the state of this tab
type Model struct {
	colors *theme.Colors
	reader backend.Fetcher
	title  string
	keymap Keymap
	list   simplelist.Model
	width  int
	height int
	loaded bool
}

// New creates a new category tab with sensible defaults
func New(colors *theme.Colors, width, height int, title string, fetcher backend.Fetcher) Model {
	log.Println("Creating new category tab with title", title)

	return Model{
		colors: colors,
		width:  width,
		height: height,
		title:  title,
		reader: fetcher,
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
	if !m.loaded {
		return m
	}

	m.width = width
	m.height = height
	m.list.SetHeight(m.height)
	return m
}

// Init initializes the tab
func (m Model) Init() tea.Cmd {
	return m.reader(m.title)
}

// Update updates the variables of the tab
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case backend.FetchSuccessMsg:
		if !m.loaded {
			m.list = simplelist.New(m.colors, m.title, m.height, false)
			m.loaded = true
		}

		m.list.SetItems(msg.Items)
		return m, nil

	case popup.ChoiceResultMsg:
		if !msg.Result {
			return m, nil
		}

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
		if !m.loaded {
			return m, nil
		}

		switch {
		case key.Matches(msg, m.list.Keymap.Open):
			if !m.list.IsEmpty() {
				return m, tab.NewTab(m, m.list.SelectedItem().FilterValue())
			}

			return m, nil

		case key.Matches(msg, m.keymap.NewFeed):
			return m, backend.NewItem(m)

		case key.Matches(msg, m.keymap.EditFeed):
			// If the list is empty, return nothing
			if m.list.IsEmpty() {
				return m, nil
			}

			item := m.list.SelectedItem().(simplelist.Item)
			fields := []string{item.Title(), item.Description()}
			return m, backend.EditItem(m, fields)

		case key.Matches(msg, m.keymap.DeleteFeed):
			if !m.list.IsEmpty() {
				return m, backend.MakeChoice("Delete this feed?", true)
			}

		default:
			if item, ok := m.list.GetItem(msg.String()); ok {
				return m, tab.NewTab(m, item.FilterValue())
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View returns the view of the tab
func (m Model) View() string {
	if !m.loaded {
		return "Loading..."
	}

	return m.list.View()
}

// ShortHelp returns the short help for this tab
func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{m.keymap.NewFeed, m.keymap.EditFeed, m.keymap.DeleteFeed}
}

// FullHelp returns the full help for this tab
func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp(), m.list.ShortHelp()}
}
