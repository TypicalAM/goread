package list

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

// List contains the list of items
type List struct {
	// title is the title of the list
	title string
	// height is the height of the list
	height int
	// items is the list of items
	items []list.Item
	// selected is the index of the selected item
	selected int
	// titleStyle is the style of the title
	titleStyle lipgloss.Style
	// itemStyle is the style of an item
	itemStyle lipgloss.Style
}

// Create a new list with sane defaults
func NewList(title string, height int) List {
	return List{
		title:  title,
		height: height,
		titleStyle: lipgloss.NewStyle().
			Foreground(style.BasicColorscheme.Color1).
			PaddingBottom(1),
		itemStyle: lipgloss.NewStyle().
			Foreground(style.BasicColorscheme.Color2),
	}
}

// Set the list of items
func (l *List) SetItems(items []list.Item) {
	// TODO: keep calm, don't panic, handle this in a status message
	if len(items) > 36 {
		panic("List: too many items")
	}

	// Set the items
	l.items = items
}

// Get the item at a specified index
func (l List) GetItem(index int) list.Item {
	return l.items[index]
}

// Init the list
func (l List) Init() tea.Cmd {
	return nil
}

// View the list
func (l List) View() string {
	// We will store the sections here and join them later
	sections := []string{""}

	// Create the title
	titleText := l.titleStyle.MarginLeft(3).Render(l.title)
	sections = append(sections, titleText)

	// If the list is empty return early
	if len(l.items) == 0 {
		noItemsText := lipgloss.NewStyle().
			MarginLeft(3).
			Foreground(style.BasicColorscheme.Color2).
			Italic(true)
		sections = append(sections, noItemsText.Render("<no items>"))
		return lipgloss.JoinVertical(lipgloss.Top, sections...)
	}

	// If the list has items style it
	for i, item := range l.items {
		isSelected := i == l.selected
		itemLine := lipgloss.NewStyle().
			MarginLeft(3).
			MarginRight(3).
			Render(fmt.Sprintf("%s  %s",
				style.StyleIndex(i, isSelected),
				l.itemStyle.Render(item.FilterValue())),
			)
		sections = append(sections, itemLine)
	}

	// Append a blank line at the end
	sections = append(sections, "")

	// Add padding
	sections = append(sections, strings.Repeat("\n", l.height-len(sections)))
	return lipgloss.JoinVertical(lipgloss.Top, sections...)
}

// Update the list, things like scrolling
func (l List) Update(msg tea.Msg) (List, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if l.selected > 0 {
				l.selected--
			}
		case "down":
			if l.selected < len(l.items)-1 {
				l.selected++
			}
		}
	}

	return l, nil
}

// Return the selected item
func (l List) SelectedItem() list.Item {
	return l.items[l.selected]
}

// Check if an item by this index is in the list
func (l List) HasItem(text string) (int, bool) {
	// Convert the text to an integer and
	// check if the index is in the list
	if index, err := strconv.Atoi(text); err == nil {
		return index, index < len(l.items)
	}

	// If the text is not an integer, check if
	// it is a lowercase letter
	if !unicode.IsLower(rune(text[0])) {
		return 0, false
	}

	// Convert the letter to an index since 97
	// is the ASCII code for 'a' and we have 10 digits
	index := int(text[0]) - 97 + 10
	return index, index < len(l.items)
}
