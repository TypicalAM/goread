package simplelist

import (
	"strconv"

	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/ansi"
)

// Item is an item in the list
type Item struct {
	title string
	desc  string
}

// NewItem creates a new item
func NewItem(title, desc string) Item {
	return Item{
		title: title,
		desc:  desc,
	}
}

// Title returns the title of the item
func (i Item) Title() string {
	return i.title
}

// Description returns the description of the item
func (i Item) Description() string {
	return i.desc
}

// FilterValue returns the title of the item
func (i Item) FilterValue() string {
	return i.title
}

// Model contains state of the list
type Model struct {
	colors       *colorscheme.Colorscheme
	style        listStyle
	title        string
	items        []list.Item
	height       int
	page         int
	itemsPerPage int
	selected     int
	showDesc     bool
}

// New creates a new list
func New(colors *colorscheme.Colorscheme, title string, height int, showDesc bool) Model {
	style := newListStyle(colors)
	var itemsPerPage int
	if showDesc {
		itemsPerPage = (height - lipgloss.Height(style.titleStyle.Render(""))) / 2
	} else {
		itemsPerPage = height - lipgloss.Height(style.titleStyle.Render(""))
	}

	return Model{
		colors:       colors,
		title:        title,
		height:       height,
		itemsPerPage: itemsPerPage,
		showDesc:     showDesc,
		style:        style,
	}
}

// Init initializes the tab
func (m Model) Init() tea.Cmd {
	return nil
}

// Update updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up", "k":
			m.selected--
			if m.selected < 0 {
				m.selected = len(m.items) - 1
				m.page = len(m.items) / m.itemsPerPage
			}

			// Check if the page needs to be changed
			if m.selected < m.page*m.itemsPerPage {
				m.page--
			}

		case "down", "j":
			m.selected++
			if m.selected >= len(m.items) {
				m.selected = 0
				m.page = 0
			}

			if m.selected >= (m.page+1)*m.itemsPerPage {
				m.page++
			}

		case "shift+up", "K":
			m.selected = 0
			m.page = 0

		case "shift+down", "J":
			m.selected = len(m.items) - 1
			m.page = len(m.items) / m.itemsPerPage
		}
	}

	return m, nil
}

// View returns the view of the list
func (m Model) View() string {
	sections := make([]string, 1)

	sections = append(sections, m.style.titleStyle.Render(m.title))

	if len(m.items) == 0 {
		sections = append(sections, m.style.noItemsStyle.Render("<no items>"))
		return lipgloss.JoinVertical(lipgloss.Top, sections...)
	}

	for i := m.itemsPerPage * m.page; i < m.itemsPerPage*(m.page+1); i++ {
		if i >= len(m.items) {
			break
		}

		isSelected := i == m.selected
		sections = append(sections, lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.style.styleIndex(i, isSelected),
			m.style.itemStyle.Render(m.items[i].FilterValue()),
		))

		if m.showDesc {
			if item, ok := m.items[i].(list.DefaultItem); ok {
				if ansi.PrintableRuneWidth(item.Description()) != 0 {
					sections = append(sections, m.style.styleDescription(item.Description()))
				} else {
					sections = append(sections, "")
				}
			}
		}
	}

	sections = append(sections, "")
	return lipgloss.JoinVertical(lipgloss.Top, sections...)
}

// SetHeight sets the height of the list
func (m *Model) SetHeight(height int) {
	if m.showDesc {
		m.itemsPerPage = (height - lipgloss.Height(m.style.titleStyle.Render(""))) / 2
	} else {
		m.itemsPerPage = height - lipgloss.Height(m.style.titleStyle.Render(""))
	}

	m.height = height
}

// Items returns the items in the list
func (m Model) Items() []list.Item {
	return m.items
}

// SetItems sets the items in the list
func (m *Model) SetItems(items []list.Item) {
	if len(items) > 36 {
		panic("list: too many items")
	}

	m.items = items
}

// IsEmpty checks if the list is empty
func (m *Model) IsEmpty() bool {
	return len(m.items) == 0
}

// SelectedItem returns the selected item
func (m Model) SelectedItem() list.Item {
	return m.items[m.selected]
}

// GetItem checks if the list has an item and returns it
func (m Model) GetItem(text string) (list.Item, bool) {
	index, err := strconv.Atoi(text)
	if err != nil {
		return nil, false
	}

	if index >= len(m.items) || index < 0 {
		return nil, false
	}

	return m.items[index], true
}

// Index returns the index of the selected item
func (m Model) Index() int {
	return m.selected
}

// SetIndex sets the index of the selected item
func (m *Model) SetIndex(index int) {
	m.selected = index
}
