package tab

import (
	simpleList "github.com/TypicalAM/goread/internal/list"
	"github.com/TypicalAM/goread/internal/style"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Welcome struct {
	// General fields
	title  string
	index  int
	loaded bool

	// The list of categorie
	list simpleList.List
}

// New creates a new RssFeedTab with sensible defautls
func NewWelcomeTab(title string, index int) Welcome {
	return Welcome{
		title: title,
		index: index,
	}
}

// Return the title of the tab
func (w Welcome) Title() string {
	return w.title
}

// Return the index of the tab
func (w Welcome) Index() int {
	return w.index
}

// Return the type of the tab
func (w Welcome) Type() TabType {
	return WelcomeTab
}

// Return the load state
func (w Welcome) Loaded() bool {
	return w.loaded
}

// Create the list of categories
func (w *Welcome) initList(height int) {
	defaultList := simpleList.NewList("Categories", height-5)
	w.list = defaultList

	// Add the categories
	w.list.SetItems(
		[]list.Item{
			simpleList.NewListItem("All", "All the categories", ""),
			simpleList.NewListItem("Books", "Books", "books"),
			simpleList.NewListItem("Movies", "Movies", "movies"),
			simpleList.NewListItem("Music", "Music", "music"),
			simpleList.NewListItem("Games", "Games", "games"),
			simpleList.NewListItem("Technology", "Software", "software"),
		},
	)
}

// Implement the bubbletea.Model interface
func (w Welcome) Init() tea.Cmd {
	return nil
}

// Update the variables
func (w Welcome) Update(msg tea.Msg) (Tab, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Check if the program is loaded, if not, load it
	if !w.loaded && style.WindowWidth > 0 && style.WindowHeight > 0 {
		w.initList(style.WindowHeight)
		w.loaded = true
		return w, nil
	}

	// Update the list
	w.list, cmd = w.list.Update(msg)
	cmds = append(cmds, cmd)

	// Check the message type
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check if we need to open a new tab
		if index, ok := w.list.HasItem(msg.String()); ok {
			cmds = append(cmds, NewTab(w.list.GetItem(index).FilterValue(), CategoryTab))
		}

		// Check if the user has pressed enter
		if msg.String() == "enter" {
			cmds = append(cmds, NewTab(w.list.SelectedItem().FilterValue(), CategoryTab))
		}
	}

	return w, tea.Batch(cmds...)
}

// View the tab
func (w Welcome) View() string {
	// Check if the program is loaded, if not, return an empty string
	if !w.loaded {
		return "Loading..."
	}

	// Return the view
	return w.list.View()
}

// Set the index of the tab
func (w Welcome) SetIndex(index int) Tab {
	w.index = index
	return w
}
