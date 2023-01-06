package input

import (
	"fmt"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	Normal State = iota
	Cancel
	Finished
	NotEnoughText
)

// Model contains the state of this tab
type Model struct {
	State       State
	activeInput int
	fields      []string
	inputs      []textinput.Model
	Creating    bool
	Type        backend.ItemType
}

// New creates a new instance of the create item model
func New(itemType backend.ItemType, creating bool, fields []string) Model {
	// Create an empty instance
	c := Model{}

	// Set the type
	c.Type = itemType

	// Are we creating a new item?
	c.Creating = creating

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
func (m Model) Init() tea.Cmd {
	return nil
}

// Update the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Update the textfields
	var cmd tea.Cmd
	m.inputs[m.activeInput], cmd = m.inputs[m.activeInput].Update(msg)

	// Check if we pressed enter, if yes, switch to the next input
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			m.activeInput++
			if m.activeInput >= len(m.inputs) {
				// If any of the inputs are empty, return
				for i := range m.inputs {
					if m.inputs[i].Value() == "" {
						m.State = NotEnoughText
						return m, nil
					}
				}

				// If we are here, all inputs are filled
				m.State = Finished
				return m, nil
			}

		case "esc":
			m.State = Cancel
			return m, nil
		}
	}

	// If we are not creating new items, we need to update the tabs
	return m, cmd
}

// View the selected input
func (m Model) View() string {
	return m.inputs[m.activeInput].View()
}

// Index returns the index of the active input
func (m Model) Index() int {
	return m.activeInput
}

// GetValues returns the values of the inputs
func (m Model) GetValues() []string {
	values := make([]string, len(m.inputs))
	for i := range m.inputs {
		values[i] = m.inputs[i].Value()
	}
	return values
}
