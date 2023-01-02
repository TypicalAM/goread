package rss

import (
	"fmt"

	"github.com/TypicalAM/goread/internal/tab"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// CreateItem is used as a model when creating a new item
type CreateItem struct {
	activeInput int

	fields []string
	inputs []textinput.Model
	Type   tab.Type
}

// New creates a new instance of the create item model
func NewItemCreation(fields []string, tabType tab.Type) CreateItem {
	// Create an empty instance
	c := CreateItem{}

	// Set the fields
	c.fields = fields

	// Create the textfields
	c.inputs = make([]textinput.Model, len(fields))

	// Set the textfields
	for i := range c.inputs {
		t := textinput.New()
		t.Focus()
		t.Prompt = fmt.Sprintf("Enter %s: ", fields[i])
		c.inputs[i] = t
	}

	// Set the active textbox
	c.activeInput = 0
	return c
}

// Init initializes the model
func (c CreateItem) Init() tea.Cmd {
	return nil
}

// Update the model
func (c CreateItem) Update(msg tea.Msg) (CreateItem, tea.Cmd) {
	// Update the textfields
	var cmd tea.Cmd
	c.inputs[c.activeInput], cmd = c.inputs[c.activeInput].Update(msg)

	// Check if we pressed enter, if yes, switch to the next input
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			c.activeInput++
			if c.activeInput >= len(c.inputs) {
				c.activeInput = -1
				return c, nil
			}

		case "esc":
			c.activeInput = -1
			return c, nil
		}
	}

	// If we are not creating new items, we need to update the tabs
	return c, cmd
}

// View the selected input
func (c CreateItem) View() string {
	return c.inputs[c.activeInput].View()
}

// Index() returns the index of the active input
func (c CreateItem) Index() int {
	return c.activeInput
}

// GetValues returns the values of the inputs
func (c CreateItem) GetValues() []string {
	values := make([]string, len(c.inputs))
	for i := range c.inputs {
		values[i] = c.inputs[i].Value()
	}
	return values
}
