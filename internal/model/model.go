package model

import (
	"fmt"
	"strings"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/style"
	"github.com/TypicalAM/goread/internal/tab"
	"github.com/TypicalAM/goread/internal/tab/category"
	"github.com/TypicalAM/goread/internal/tab/feed"
	"github.com/TypicalAM/goread/internal/tab/welcome"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tabs      []tab.Tab
	backend   backend.Backend
	activeTab int
	message   string
	quitting  bool
	loaded    bool

	creatingItem bool
	createItem   createItem
}

// NewModel returns a new model with some sensible defaults
func New(backend backend.Backend) Model {
	model := Model{}
	model.backend = backend
	model.message = fmt.Sprintf("Using backend - %s", backend.Name())
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Create the command array to pass it when updating if there are more than one model
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Check if we have the window size, if not, we wait for it
	if !m.loaded {
		if msg, ok := msg.(tea.WindowSizeMsg); ok {
			style.WindowWidth = msg.Width
			style.WindowHeight = msg.Height
			m.tabs = append(m.tabs, welcome.New("Welcome", m.backend.FetchCategories))
			m.loaded = true
			cmds = append(cmds, m.tabs[0].Init())
		} else {
			return m, nil
		}
	}

	// If we are creating new items, we need to update the inputs
	if m.creatingItem {
		m.createItem, cmd = m.createItem.Update(msg)
		cmds = append(cmds, cmd)

		if m.createItem.Index() == -1 {
			m.creatingItem = false
			m.message = "Ended " + strings.Join(m.createItem.GetValues(), " ")
		}

		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case backend.FetchErrorMessage:
		// If there is an error, display it on the status bar
		// the error message will be cleared when the user closes the tab
		m.message = fmt.Sprintf("%s - %s", msg.Description, msg.Err.Error())

	case tab.NewTabMessage:
		// Create the new tab
		m.createNewTab(msg.Title, msg.Type)

		// Initialize the tab and do the first update
		cmds = append(cmds, m.tabs[m.activeTab].Init())
		m.tabs[m.activeTab], cmd = m.tabs[m.activeTab].Update(msg)
		cmds = append(cmds, cmd)

		// Clear the message
		m.message = ""

	case backend.NewItemMessage:
		// Initialize the textfields and move into the
		// textbox view
		m.createItem = newItemCreation(msg.Fields, msg.TabType)
		m.creatingItem = true

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			// Quit the program
			m.quitting = true
			return m, tea.Quit

		case "tab":
			m.activeTab++
			// Wrap around
			if m.activeTab > len(m.tabs)-1 {
				m.activeTab = 0
			}

			// If it is not loaded, initialize it
			if !m.tabs[m.activeTab].Loaded() {
				cmds = append(cmds, m.tabs[m.activeTab].Init())
			}

			// Clear the message
			m.message = ""

		case "shift+tab":
			m.activeTab--
			if m.activeTab < 0 {
				m.activeTab = len(m.tabs) - 1
			}

			// Clear the current message
			m.message = ""

		case "ctrl+w":
			// If there is only one tab, quit
			if len(m.tabs) == 1 {
				m.quitting = true
				return m, tea.Quit
			}

			// Close the current tab
			m.message = fmt.Sprintf("Closed tab - %s", m.tabs[m.activeTab].Title())
			m.tabs = append(m.tabs[:m.activeTab], m.tabs[m.activeTab+1:]...)
			m.activeTab--

			// Wrap around
			if m.activeTab < 0 {
				m.activeTab = 0
			}
		}
	}

	// Call the tab model and update its variables
	m.tabs[m.activeTab], cmd = m.tabs[m.activeTab].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) RenderTabBar() string {
	// Render the tab bar at the top of the screen
	tabs := make([]string, len(m.tabs))
	for i, tabObj := range m.tabs {
		tabs[i] = tab.AttachIconToTab(tabObj.Title(), tabObj.Type(), i == m.activeTab)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	gap := style.TabGap.Render(strings.Repeat(" ", style.Max(0, style.WindowWidth-lipgloss.Width(row))))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
}

func (m *Model) RenderStatusBar() string {
	// Render the status bar at the bottom of the screen
	row := lipgloss.JoinHorizontal(lipgloss.Top, tab.StyleStatusBarCell(m.tabs[m.activeTab].Type()))
	gap := style.StatusBarGap.Render(strings.Repeat(" ", style.Max(0, style.WindowWidth-lipgloss.Width(row))))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
}

func (m Model) View() string {
	// If we are quitting, render the quit message
	if m.quitting {
		return "Goodbye!"
	}

	// If we are not loaded, render the loading message
	if !m.loaded {
		return "Loadings..."
	}

	// Hold the sections of the screen
	var sections []string

	// Do not render the tab bar if there is only one tab
	sections = append(sections, m.RenderTabBar())

	// Render the tab content and the status bar
	constrainHeight := lipgloss.NewStyle().Height(style.WindowHeight - 3)
	sections = append(sections, constrainHeight.Render(m.tabs[m.activeTab].View()))
	sections = append(sections, m.RenderStatusBar())

	// If we are typing, shift the focus onto the textfield
	if m.creatingItem {
		// FIXME: Section name?
		sections = append(sections, m.createItem.View())
	} else {
		sections = append(sections, m.message)
	}

	// Render the message if there is one
	return lipgloss.JoinVertical(lipgloss.Top, sections...)
}

// Create the new tab and add it to the model
func (m *Model) createNewTab(title string, tabType tab.Type) {
	// Create and add the new tab
	var newTab tab.Tab

	// Create a new tab based on the type
	switch tabType {
	case tab.Category:
		newTab = category.New(
			title,
			m.backend.FetchFeeds,
		)
	case tab.Feed:
		newTab = feed.New(
			title,
			m.backend.FetchArticles,
		)
	}

	// Insert the tab after the active tab
	m.tabs = append(m.tabs[:m.activeTab+1], append([]tab.Tab{newTab}, m.tabs[m.activeTab+1:]...)...)

	// Increase the active tab count
	m.activeTab++
}
