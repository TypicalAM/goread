package welcome

import (
	"github.com/TypicalAM/goread/internal/backend"
	simpleList "github.com/TypicalAM/goread/internal/list"
	"github.com/TypicalAM/goread/internal/style"
	"github.com/TypicalAM/goread/internal/tab"
	tea "github.com/charmbracelet/bubbletea"
)

type Welcome struct {
	// General fields
	title  string
	index  int
	loaded bool

	// The list of categorie
	list       simpleList.List
	readerFunc func() tea.Cmd
}

// New creates a new RssFeedTab with sensible defautls
func New(title string, index int, readerFunc func() tea.Cmd) Welcome {
	return Welcome{
		title:      title,
		index:      index,
		readerFunc: readerFunc,
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
func (w Welcome) Type() tab.TabType {
	return tab.Welcome
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
	w.list.SetItems(w.readerFunc()().(backend.FetchSuccessMessage).Items)
}

// Implement the bubbletea.Model interface
func (w Welcome) Init() tea.Cmd {
	return nil
}

// Update the variables
func (w Welcome) Update(msg tea.Msg) (tab.Tab, tea.Cmd) {
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
			cmds = append(cmds, tab.NewTab(w.list.GetItem(index).FilterValue(), tab.Category))
		}

		// Check if the user has pressed enter
		if msg.String() == "enter" {
			cmds = append(cmds, tab.NewTab(w.list.SelectedItem().FilterValue(), tab.Category))
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
func (w Welcome) SetIndex(index int) tab.Tab {
	w.index = index
	return w
}
