package simplelist

import (
	"strconv"
	"unicode"

	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// Update the variables in the list
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Handle key presses
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up":
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

		case "down":
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

		case "shift+up":
			// Select the first item
			m.selected = 0
			m.page = 0

		case "shift+down":
			// Select the last item
			m.selected = len(m.items) - 1
			m.page = len(m.items) / m.itemsPerPage
		}
	}

	// Return the updated model
	return m, nil
}

// View the list
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
			sections = append(sections,
				m.style.styleDescription(m.items[i].(Item).Description()),
			)
		}
	}

	// Append a blank line at the end
	sections = append(sections, "")

	// Join the sections and return the view
	return lipgloss.JoinVertical(lipgloss.Top, sections...)
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

	// If the text is not an integer, check if
	// it is a lowercase letter
	if !unicode.IsLower(rune(text[0])) {
		return nil, false
	}

	// Convert the letter to an index since 97
	// is the ASCII code for 'a' and we have 10 digits
	index := int(text[0]) - 97 + 10
	inList := index < len(m.items)
	if !inList {
		return nil, false
	}

	// Return the item
	return m.items[index], true
}
