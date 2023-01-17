package browser

import (
	"fmt"
	"strings"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/config"
	"github.com/TypicalAM/goread/internal/model/input"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/TypicalAM/goread/internal/model/tab/category"
	"github.com/TypicalAM/goread/internal/model/tab/feed"
	"github.com/TypicalAM/goread/internal/model/tab/welcome"
	"github.com/TypicalAM/goread/internal/rss"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is used to store the state of the application
type Model struct {
	// config is the config of the application
	config config.Config
	style  style

	// managing tabs
	tabs      []tab.Tab
	activeTab int

	// window size
	waitingForSize bool
	windowWidth    int
	windowHeight   int

	// creating items
	newItem bool
	input   input.Model

	// other
	message  string
	quitting bool
}

// New returns a new model with some sensible defaults
func New(cfg config.Config) Model {
	return Model{
		config:         cfg,
		style:          newStyle(cfg.Colors),
		waitingForSize: true,
		message:        "Pro-tip - press [ctrl-h] to view the help page",
	}
}

// Init initializes the model, there are no I/O operations needed
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles the terminal size, modifying rss items
// and modifying tabs
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check if we have the window size, if not, we wait for it
	if m.waitingForSize {
		return m.waitForSize(msg)
	}

	// If we are creating new items, we need to update the inputs
	if m.newItem {
		return m.updateItemCreation(msg)
	}

	switch msg := msg.(type) {
	case backend.FetchErrorMessage:
		// If there is an error, display it on the status bar
		// the error message will be cleared when the user closes the tab
		m.message = fmt.Sprintf("%s - %s", msg.Description, msg.Err.Error())

		// Update the underlying tab in case it also handles error input
		m.tabs[m.activeTab], _ = m.tabs[m.activeTab].Update(msg)
		return m, nil

	case tab.NewTabMessage:
		// Create the new tab
		m.createNewTab(msg.Title, msg.Type)
		m.message = ""
		return m, m.tabs[m.activeTab].Init()

	case backend.NewItemMessage:
		// Initialize the new Item model
		m.input = input.New(msg.Type, msg.New, msg.Fields, msg.ItemPath, msg.OldFields)
		m.newItem = true
		return m, m.input.Init()

	case backend.DeleteItemMessage:
		// Delete the item
		return m.deleteItem(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			// Quit the program
			m.quitting = true
			return m, tea.Quit

		case "tab":
			// Cycle through the tabs
			m.activeTab++
			if m.activeTab > len(m.tabs)-1 {
				m.activeTab = 0
			}

			// Clear the message
			m.message = ""
			return m, nil

		case "shift+tab":
			// Cycle through the tabs
			m.activeTab--
			if m.activeTab < 0 {
				m.activeTab = len(m.tabs) - 1
			}

			// Clear the current message
			m.message = ""
			return m, nil

		case "ctrl+w":
			// If there is only one tab, quit
			if len(m.tabs) == 1 {
				m.quitting = true
				return m, tea.Quit
			}

			// Close the current tab
			m.tabs = append(m.tabs[:m.activeTab], m.tabs[m.activeTab+1:]...)
			m.activeTab--

			// Wrap around
			if m.activeTab < 0 {
				m.activeTab = 0
			}

			// Set the message
			m.message = fmt.Sprintf("Closed tab - %s", m.tabs[m.activeTab].Title())
			return m, nil

		case "ctrl+h":
			// View the help page
			return m.showHelp()
		}
	}

	// Call the tab model and update its variables
	var cmd tea.Cmd
	m.tabs[m.activeTab], cmd = m.tabs[m.activeTab].Update(msg)
	return m, tea.Batch(cmd)
}

// View renders the tab bar, the active tab and the status bar
func (m Model) View() string {
	// If we are quitting, render the quit message
	if m.quitting {
		return "Goodbye!"
	}

	// If we are not loaded, render the loading message
	if m.waitingForSize {
		return "Loading..."
	}

	// Hold the sections of the screen
	var sections []string

	// Do not render the tab bar if there is only one tab
	sections = append(sections, m.renderTabBar())

	// Render the tab content and the status bar
	constrainHeight := lipgloss.NewStyle().Height(m.windowHeight - 3)
	sections = append(sections, constrainHeight.Render(m.tabs[m.activeTab].View()))

	// Render the status bar
	sections = append(sections, m.renderStatusBar())

	// If we are typing, shift the focus onto the text-field
	var messageBar string
	if m.newItem {
		messageBar = m.input.View()
	} else {
		messageBar = m.message
	}

	// Render the message bar
	sections = append(sections, messageBar)

	// Join all the sections
	return lipgloss.JoinVertical(lipgloss.Top, sections...)
}

// waitForSize waits for the window size to be set and loads the tab
// if it receives it
func (m Model) waitForSize(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Wait for the window size message
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		// Initialize the window height and width
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.waitingForSize = false

		// Append a new welcome tab
		m.tabs = append(m.tabs, welcome.New(
			m.config.Colors,
			m.windowWidth,
			m.windowHeight-5,
			"Welcome",
			m.config.Backend.FetchCategories,
		))

		// Return the init of the tab
		return m, m.tabs[0].Init()
	}

	// Return nothing if we didn't get size yet
	return m, nil
}

// updateItemCreation updates the child model for creating items
func (m Model) updateItemCreation(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	// If the child model is done, add the item
	switch m.input.State {
	case input.NotEnoughText:
		m.message = "Fields cannot be blank!"
		m.newItem = false
		return m, cmd

	case input.Cancel:
		m.message = "Cancelled adding or editing item"
		m.newItem = false
		return m, cmd

	case input.Finished:
		return m.addItem()

	default:
		return m, cmd
	}
}

// Create the new tab and add it to the model
func (m *Model) createNewTab(title string, tabType tab.Type) {
	// Create and add the new tab
	var newTab tab.Tab

	// Create a new tab based on the type
	switch tabType {
	case tab.Category:
		newTab = category.New(
			m.config.Colors,
			m.windowWidth,
			m.windowHeight-5,
			title,
			m.config.Backend.FetchFeeds,
		)
	case tab.Feed:
		newTab = feed.New(
			m.config.Colors,
			m.windowWidth,
			m.windowHeight-5,
			title,
			m.config.Backend.FetchArticles,
		)
	}

	// If it is a category tab with the name of rss.AllFeedsName, create a new feed tab
	if tabType == tab.Category && title == rss.AllFeedsName {
		newTab = feed.New(
			m.config.Colors,
			m.windowWidth,
			m.windowHeight-5,
			title,
			m.config.Backend.FetchArticles,
		)
	}

	// Insert the tab after the active tab
	m.tabs = append(m.tabs[:m.activeTab+1], append([]tab.Tab{newTab}, m.tabs[m.activeTab+1:]...)...)

	// Increase the active tab count
	m.activeTab++
}

// addItem gets the data from the child model and adds it to the rss
func (m Model) addItem() (tea.Model, tea.Cmd) {
	// End creating new items
	m.newItem = false
	values := m.input.GetValues()

	// Check if we are creating or editing the item
	if m.input.Creating {
		m.message = "Adding an item - " + strings.Join(values, " ")
	} else {
		m.message = "Editing an item - " + strings.Join(values, " ")
	}

	// Check if the values are valid
	if m.input.Type == backend.Category {
		var err error
		if m.input.Creating {
			err = m.config.Backend.Rss.AddCategory(values[0], values[1])
		} else {
			err = m.config.Backend.Rss.UpdateCategory(m.input.Path[0], values[0], values[1])
		}

		// Check if there was an error
		if err != nil {
			m.message = "Error adding or updating category: " + err.Error()
			return m, nil
		}

		// Refresh the categories
		return m, m.config.Backend.FetchCategories()
	}

	// Check if the feed already exists
	var err error
	if m.input.Creating {
		err = m.config.Backend.Rss.AddFeed(m.tabs[m.activeTab].Title(), values[0], values[1])
	} else {
		err = m.config.Backend.Rss.UpdateFeed(m.input.Path[0], m.input.Path[1], values[0], values[1])
	}

	// Check if there was an error
	if err != nil {
		m.message = "Error adding or updating feed: " + err.Error()
		return m, nil
	}

	// Refresh the feeds
	return m, m.config.Backend.FetchFeeds(m.tabs[m.activeTab].Title())
}

// deleteItem deletes the focused item from the backend
func (m Model) deleteItem(msg backend.DeleteItemMessage) (tea.Model, tea.Cmd) {
	m.message = fmt.Sprintf("Deleting item %s", msg.Key)

	// Check the type of the item
	if msg.Type == backend.Category {
		err := m.config.Backend.Rss.RemoveCategory(msg.Key)
		if err != nil {
			m.message = fmt.Sprintf("Error deleting category %s - %s", msg.Key, err.Error())
		}

		// Refresh the categories
		return m, m.config.Backend.FetchCategories()
	}

	// Delete the feed
	err := m.config.Backend.Rss.RemoveFeed(m.tabs[m.activeTab].Title(), msg.Key)
	if err != nil {
		m.message = fmt.Sprintf("Error deleting feed %s - %s", msg.Key, err.Error())
	}

	// Fetch the feeds again to update the list
	return m, m.config.Backend.FetchFeeds(m.tabs[m.activeTab].Title())
}

// showHelp() shows the help menu at the bottom of the screen
func (m Model) showHelp() (tea.Model, tea.Cmd) {
	// Create the help menu
	message := "Help: [ctrl+w] Close tab, [Tab] Cycle tabs, "
	for _, keyBind := range m.tabs[m.activeTab].Help() {
		message += fmt.Sprintf("[%s] %s, ", keyBind.Key, keyBind.Description)
	}

	// Set the message
	m.message = message[:len(message)-2]
	return m, nil
}

// renderTabBar renders the tab bar at the top of the screen
func (m *Model) renderTabBar() string {
	// Render the tab bar at the top of the screen
	tabs := make([]string, len(m.tabs))
	for i, tabObj := range m.tabs {
		tabs[i] = m.style.attachIconToTab(tabObj.Title(), tabObj.Type(), i == m.activeTab)
	}

	// Check if the row exceeds the width of the screen
	if lipgloss.Width(strings.Join(tabs, "")) > m.windowWidth {
		// Trim the tabs to fit the screen
		tabs = tabs[m.activeTab:]
	}

	// Create the row
	row := strings.Join(tabs, "")

	// Calculate the gap amount
	var gapAmount int
	if m.windowWidth-lipgloss.Width(row) < 0 {
		gapAmount = 0
	} else {
		gapAmount = m.windowWidth - lipgloss.Width(row)
	}

	// Create the gap on the right
	gap := m.style.tabGap.Render(strings.Repeat(" ", gapAmount))
	return lipgloss.JoinHorizontal(lipgloss.Left, row, gap)
}

// renderStatusBar is used to render the status bar at the bottom of the screen
func (m *Model) renderStatusBar() string {
	// Render the status bar at the bottom of the screen
	row := lipgloss.JoinHorizontal(lipgloss.Top, m.style.styleStatusBarCell(m.tabs[m.activeTab].Type()))

	// Calculate the gap amount
	var gapAmount int
	if m.windowWidth-lipgloss.Width(row) < 0 {
		gapAmount = 0
	} else {
		gapAmount = m.windowWidth - lipgloss.Width(row)
	}

	// Render the gap on the right
	gap := m.style.statusBarGap.Render(strings.Repeat(" ", gapAmount))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
}
