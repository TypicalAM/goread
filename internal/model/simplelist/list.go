package simplelist

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/TypicalAM/goread/internal/style"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model contains state of the list
type Model struct {
	title        string
	height       int
	items        []list.Item
	selected     int
	titleStyle   lipgloss.Style
	noItemsStyle lipgloss.Style
	itemStyle    lipgloss.Style
}

// New creates a new list
func New(title string, height int) Model {
	return Model{
		title:  title,
		height: height,
		titleStyle: lipgloss.NewStyle().
			Foreground(style.GlobalColorscheme.Color1).
			MarginLeft(3).
			PaddingBottom(1),
		noItemsStyle: lipgloss.NewStyle().
			MarginLeft(3).
			Foreground(style.GlobalColorscheme.Color2).
			Italic(true),
		itemStyle: lipgloss.NewStyle().
			MarginLeft(3).
			MarginRight(3).
			Foreground(style.GlobalColorscheme.Color2),
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
		case "j", "down":
			m.selected++
			if m.selected >= len(m.items) {
				m.selected = 0
			}
		case "k", "up":
			m.selected--
			if m.selected < 0 {
				m.selected = len(m.items) - 1
			}
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
	sections = append(sections, m.titleStyle.Render(m.title))

	// If the list is empty show a message
	if len(m.items) == 0 {
		sections = append(sections, m.noItemsStyle.Render("<no items>"))
		return lipgloss.JoinVertical(lipgloss.Top, sections...)
	}

	// If the list has items, style them
	for i, item := range m.items {
		isSelected := i == m.selected
		itemText := fmt.Sprintf("%s  %s", style.Index(i, isSelected), item.FilterValue())
		itemLine := m.itemStyle.Render(itemText)
		sections = append(sections, itemLine)
	}

	// Append a blank line at the end
	sections = append(sections, "")

	// Add padding
	sections = append(sections, strings.Repeat("\n", m.height-len(sections)))
	return lipgloss.JoinVertical(lipgloss.Top, sections...)
}

// SetItems sets the items in the list
func (m *Model) SetItems(items []list.Item) {
	// FIXME: Error propagation
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
	// Convert the text to an integer and
	// check if the index is in the list
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
