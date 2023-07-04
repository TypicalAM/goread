package browser

import (
	"fmt"
	"strings"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/config"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/TypicalAM/goread/internal/model/tab/category"
	"github.com/TypicalAM/goread/internal/model/tab/feed"
	"github.com/TypicalAM/goread/internal/model/tab/welcome"
	"github.com/TypicalAM/goread/internal/popup"
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

	// other
	message  string
	quitting bool

	popupShown bool
	popup      popup.Popup
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

	switch msg := msg.(type) {
	case backend.FetchErrorMessage:
		// If there is an error, display it on the status bar
		// the error message will be cleared when the user closes the tab
		m.message = fmt.Sprintf("%s - %s", msg.Description, msg.Err.Error())

		// Update the underlying tab in case it also handles error input
		m.tabs[m.activeTab], _ = m.tabs[m.activeTab].Update(msg)
		return m, nil

	case category.ChosenCategoryMsg:
		m.popupShown = false
		if err := m.config.Backend.Rss.AddCategory(msg.Name, ""); err != nil {
			m.message = fmt.Sprintf("Error adding category: %s", err.Error())
		} else {
			m.message = fmt.Sprintf("Added category %s", msg.Name)
		}

		return m, m.config.Backend.FetchCategories()

	case tab.NewTabMessage:
		// Create the new tab
		m.createNewTab(msg.Title, msg.Type)
		m.message = ""
		return m, m.tabs[m.activeTab].Init()

	case backend.NewItemMessage:
		// Open a new popup
		bg := lipgloss.NewStyle().Width(m.windowWidth).Height((m.windowHeight))
		m.popup = category.NewPopup(m.style.colors, bg.Render(m.View()), m.windowWidth/2, m.windowHeight/2+m.windowHeight/4)
		m.popupShown = true
		return m, m.popup.Init()

	case backend.DeleteItemMessage:
		// Delete the item
		return m.deleteItem(msg)

	case backend.DownloadItemMessage:
		// Download the item
		return m.downloadItem(msg)

	case tea.WindowSizeMsg:
		// Resize the window
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height

		// Resize every tab
		for i := range m.tabs {
			m.tabs[i] = m.tabs[i].SetSize(m.windowWidth, m.windowHeight-5)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// Quit the program
			m.quitting = true
			return m, tea.Quit

		case "esc":
			// If we are showing a popup, close it
			if m.popupShown {
				m.popupShown = false
				return m, nil
			}

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

	// If we are showing a popup, we need to update the popup
	if m.popupShown {
		var cmd tea.Cmd
		m.popup, cmd = m.popup.Update(msg)
		return m, cmd
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

	// If we are showing a popup, render the popup
	if m.popupShown {
		return m.popup.View()
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

	// Render the message bar
	sections = append(sections, m.message)

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

// createNewTab bootstraps the new tab and adds it to the model
func (m *Model) createNewTab(title string, tabType tab.Type) {
	// Create and add the new tab
	var newTab tab.Tab

	// Create a new tab based on the type
	switch tabType {
	case tab.Category:
		switch title {
		case rss.AllFeedsName:
			newTab = feed.New(
				m.config.Colors,
				m.windowWidth,
				m.windowHeight-5,
				title,
				m.config.Backend.FetchAllArticles,
			)

		case rss.DownloadedFeedsName:
			newTab = feed.New(
				m.config.Colors,
				m.windowWidth,
				m.windowHeight-5,
				title,
				m.config.Backend.FetchDownloadedArticles,
			)

		default:
			newTab = category.New(
				m.config.Colors,
				m.windowWidth,
				m.windowHeight-5,
				title,
				m.config.Backend.FetchFeeds,
			)
		}

	case tab.Feed:
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

// downloadItem downloads an item
func (m Model) downloadItem(msg backend.DownloadItemMessage) (tea.Model, tea.Cmd) {
	m.message = fmt.Sprintf("Saving item from feed %s", msg.Key)
	return m, m.config.Backend.DownloadItem(msg.Key, msg.Index)
}

// showHelp() shows the help menu at the bottom of the screen
func (m Model) showHelp() (tea.Model, tea.Cmd) {
	m.message = m.tabs[m.activeTab].ShowHelp()
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
