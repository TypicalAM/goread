package simplelist

import (
	"strconv"

	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/ansi"
)

// Model contains state of the list
type Model struct {
	colors       colorscheme.Colorscheme
	title        string
	height       int
	page         int
	itemsPerPage int
	items        []list.Item
	selected     int
	showDesc     bool
	style        listStyle
}

// New creates a new list
func New(colors colorscheme.Colorscheme, title string, height int, showDesc bool) Model {
	// Calculate the items per page
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
	// Handle key presses
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up", "k":
			// Select the previous item
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
			// Select the next item
			m.selected++
			if m.selected >= len(m.items) {
				m.selected = 0
				m.page = 0
			}

			// Check if the page needs to be changed
			if m.selected >= (m.page+1)*m.itemsPerPage {
				m.page++
			}

		case "shift+up", "K":
			// Select the first item
			m.selected = 0
			m.page = 0

		case "shift+down", "J":
			// Select the last item
			m.selected = len(m.items) - 1
			m.page = len(m.items) / m.itemsPerPage
		}
	}

	// Return the updated model
	return m, nil
}

// View returns the view of the list
func (m Model) View() string {
	// Sections will be used to build the view
	sections := make([]string, 1)

	// Create the title
	sections = append(sections, m.style.titleStyle.Render(m.title))

	// If the list is empty show a message
	if len(m.items) == 0 {
		sections = append(sections, m.style.noItemsStyle.Render("<no items>"))
		return lipgloss.JoinVertical(lipgloss.Top, sections...)
	}

	// If the list has items, style them
	for i := m.itemsPerPage * m.page; i < m.itemsPerPage*(m.page+1); i++ {
		// Check if the index is in the list
		if i >= len(m.items) {
			break
		}

		// Render the item
		isSelected := i == m.selected
		sections = append(sections, lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.style.styleIndex(i, isSelected),
			m.style.itemStyle.Render(m.items[i].FilterValue()),
		))

		// If the description is shown add the description
		if m.showDesc {
			if ansi.PrintableRuneWidth(m.items[i].(Item).Description()) != 0 {
				sections = append(sections, m.style.styleDescription(m.items[i].(Item).Description()))
			} else {
				sections = append(sections, "")
			}
		}
	}

	// Append a blank line at the end
	sections = append(sections, "")

	// Join the sections and return the view
	return lipgloss.JoinVertical(lipgloss.Top, sections...)
}

// SetHeight sets the height of the list
func (m *Model) SetHeight(height int) {
	// Calculate the items per page
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
		panic("List: too many items")
	}

	// Set the items
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

// GetItem checks if the list has an item by a [0-9] index and also
// a [a-z] index if the list has more than 10 elements
func (m Model) GetItem(text string) (list.Item, bool) {
	// Check if the string is more than one character (left, right, up, down, etc)
	if len(text) > 1 {
		return nil, false
	}

	// Convert the text to an integer and check if the index is in the list
	if index, err := strconv.Atoi(text); err == nil {
		inList := index < len(m.items)
		if !inList {
			return nil, false
		}

		// Return the item
		return m.items[index], true
	}

	// We cannot find the item
	return nil, false
}
