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

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Keymap contains the key bindings for the browser
type Keymap struct {
	CloseTab          key.Binding
	CycleTabs         key.Binding
	ShowHelp          key.Binding
	ToggleOfflineMode key.Binding
}

// DefaultKeymap contains the default key bindings for the browser
var DefaultKeymap = Keymap{
	CloseTab: key.NewBinding(
		key.WithKeys("c", "ctrl+w"),
		key.WithHelp("c", "Close tab"),
	),
	CycleTabs: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "Cycle tabs"),
	),
	ShowHelp: key.NewBinding(
		key.WithKeys("h", "ctrl+h"),
		key.WithHelp("h", "Help"),
	),
	ToggleOfflineMode: key.NewBinding(
		key.WithKeys("o", "ctrl+o"),
		key.WithHelp("o", "Offline mode"),
	),
}

// ShortHelp returns the short help for this tab
func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.CloseTab, k.CycleTabs, k.ToggleOfflineMode,
	}
}

// FullHelp returns the full help for this tab
func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CloseTab, k.CycleTabs, k.ToggleOfflineMode},
	}
}

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
	width          int
	height         int

	// popups
	popupShown bool
	popup      popup.Popup

	// keys
	keymap Keymap
	help   help.Model

	// other
	msg      string
	quitting bool
	offline  bool
}

// New returns a new model with some sensible defaults
func New(cfg config.Config) Model {
	help := help.New()
	help.Styles.ShortDesc = lipgloss.NewStyle().Foreground(cfg.Colors.Text)
	help.Styles.ShortKey = lipgloss.NewStyle().Foreground(cfg.Colors.Text)
	help.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(cfg.Colors.TextDark)
	help.ShortSeparator = " - "

	return Model{
		config:         cfg,
		style:          newStyle(cfg.Colors),
		waitingForSize: true,
		keymap:         DefaultKeymap,
		help:           help,
		msg:            "Pro-tip - press [ctrl-h] to view the help page",
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
		m.msg = fmt.Sprintf("%s: %s", msg.Description, msg.Err.Error())

		// Update the underlying tab in case it also handles error input
		m.tabs[m.activeTab], _ = m.tabs[m.activeTab].Update(msg)
		return m, nil

	case category.ChosenCategoryMsg:
		m.popupShown = false

		if msg.IsEdit {
			if err := m.config.Backend.Rss.UpdateCategory(msg.OldName, msg.Name, msg.Desc); err != nil {
				m.msg = fmt.Sprintf("Error updating category: %s", err.Error())
			} else {
				m.msg = fmt.Sprintf("Updated category %s", msg.Name)
			}
		} else {
			if err := m.config.Backend.Rss.AddCategory(msg.Name, msg.Desc); err != nil {
				m.msg = fmt.Sprintf("Error adding category: %s", err.Error())
			} else {
				m.msg = fmt.Sprintf("Added category %s", msg.Name)
			}
		}

		return m, m.config.Backend.FetchCategories()

	case feed.ChosenFeedMsg:
		m.popupShown = false

		if msg.IsEdit {
			if err := m.config.Backend.Rss.UpdateFeed(msg.ParentCategory, msg.OldName, msg.Name, msg.URL); err != nil {
				m.msg = fmt.Sprintf("Error updating feed: %s", err.Error())
			} else {
				m.msg = fmt.Sprintf("Updated feed %s", msg.Name)
			}
		} else {
			if err := m.config.Backend.Rss.AddFeed(msg.ParentCategory, msg.Name, msg.URL); err != nil {
				m.msg = fmt.Sprintf("Error adding feed: %s", err.Error())
			} else {
				m.msg = fmt.Sprintf("Added feed %s", msg.Name)
			}
		}

		return m, m.config.Backend.FetchFeeds(msg.ParentCategory)

	case tab.NewTabMessage:
		// Create the new tab
		return m.createNewTab(msg)

	case backend.NewItemMessage:
		bg := lipgloss.NewStyle().Width(m.width).Height((m.height)).Render(m.View())
		width := m.width / 2
		height := 17

		oldName, oldDesc := msg.OldFields[0], msg.OldFields[1]

		// Open a new popup
		switch msg.Sender.(type) {
		case welcome.Model:
			m.popup = category.NewPopup(m.style.colors, bg, width, height, oldName, oldDesc)
		case category.Model:
			m.popup = feed.NewPopup(m.style.colors, bg, width, height, oldName, oldDesc, msg.Sender.Title())
		case feed.Model:
		}

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
		m.width = msg.Width
		m.height = msg.Height
		m.msg = ""

		// Resize every tab
		for i := range m.tabs {
			m.tabs[i] = m.tabs[i].SetSize(m.width, m.height-5)
		}

	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			// Quit the program
			m.quitting = true
			return m, tea.Quit

		case msg.String() == "esc":
			// If we are showing a popup, close it
			if m.popupShown {
				m.popupShown = false
				return m, nil
			}

			// Quit the program
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keymap.CloseTab):
			if m.popupShown {
				break
			}

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
			m.msg = fmt.Sprintf("Closed tab - %s", m.tabs[m.activeTab].Title())
			return m, nil

		case key.Matches(msg, m.keymap.CycleTabs):
			// Cycle through the tabs
			m.activeTab++
			if m.activeTab > len(m.tabs)-1 {
				m.activeTab = 0
			}

			// Clear the message
			m.msg = ""
			return m, nil

		case key.Matches(msg, m.keymap.ShowHelp):
			// View the help page
			if !m.popupShown {
				return m.showHelp()
			}

		case key.Matches(msg, m.keymap.ToggleOfflineMode):
			// Toggle offline mode
			m.offline = !m.offline
			m.config.Backend.SetOfflineMode(m.offline)

			if m.offline {
				m.msg = "Offline mode enabled"
			} else {
				m.msg = "Offline mode disabled"
			}

			return m, nil
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
	constrainHeight := lipgloss.NewStyle().Height(m.height - 3)
	sections = append(sections, constrainHeight.Render(m.tabs[m.activeTab].View()))

	// Render the status bar
	sections = append(sections, m.renderStatusBar())

	// Render the message bar
	if strings.Contains(m.msg, "Error") {
		errStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f08ca8")).
			Italic(true)
		sections = append(sections, errStyle.Render(m.msg))
	} else {
		sections = append(sections, m.msg)
	}

	// Join all the sections
	return lipgloss.JoinVertical(lipgloss.Top, sections...)
}

// waitForSize waits for the window size to be set and loads the tab
// if it receives it
func (m Model) waitForSize(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Wait for the window size message
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		// Initialize the window height and width
		m.width = msg.Width
		m.height = msg.Height
		m.waitingForSize = false

		// Append a new welcome tab
		m.tabs = append(m.tabs, welcome.New(
			m.config.Colors,
			m.width,
			m.height-5,
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
func (m Model) createNewTab(msg tab.NewTabMessage) (Model, tea.Cmd) {
	var newTab tab.Tab
	height := m.height - 5

	switch msg.Sender.(type) {
	case welcome.Model:
		switch msg.Title {
		case rss.AllFeedsName:
			newTab = feed.New(m.config.Colors, m.width, height, msg.Title, m.config.Backend.FetchAllArticles).
				DisableSaving().
				DisableDeleting()

		case rss.DownloadedFeedsName:
			newTab = feed.New(m.config.Colors, m.width, height, msg.Title, m.config.Backend.FetchDownloadedArticles).
				DisableSaving()

		default:
			newTab = category.New(m.config.Colors, m.width, height, msg.Title, m.config.Backend.FetchFeeds)
		}

	case category.Model:
		newTab = feed.New(m.config.Colors, m.width, height, msg.Title, m.config.Backend.FetchArticles).
			DisableDeleting()
	}

	// Insert the tab after the active tab
	m.tabs = append(m.tabs[:m.activeTab+1], append([]tab.Tab{newTab}, m.tabs[m.activeTab+1:]...)...)
	m.activeTab++
	m.msg = ""

	return m, newTab.Init()
}

// deleteItem deletes the focused item from the backend
func (m Model) deleteItem(msg backend.DeleteItemMessage) (tea.Model, tea.Cmd) {
	m.msg = fmt.Sprintf("Deleting item %s", msg.Key)

	// Check the type of the item
	switch msg.Sender.(type) {
	case welcome.Model:
		if err := m.config.Backend.Rss.RemoveCategory(msg.Key); err != nil {
			m.msg = fmt.Sprintf("Error deleting category %s: %s", msg.Key, err.Error())
		}

		return m, m.config.Backend.FetchCategories()

	case category.Model:
		if err := m.config.Backend.Rss.RemoveFeed(m.tabs[m.activeTab].Title(), msg.Key); err != nil {
			m.msg = fmt.Sprintf("Error deleting feed %s: %s", msg.Key, err.Error())
		}

		return m, m.config.Backend.FetchFeeds(m.tabs[m.activeTab].Title())

	case feed.Model:
		if msg.Sender.Title() == rss.DownloadedFeedsName {
			if err := m.config.Backend.RemoveDownload(msg.Key); err != nil {
				m.msg = fmt.Sprintf("Error deleting download %s: %s", msg.Key, err.Error())
			}
		}

		return m, m.config.Backend.FetchDownloadedArticles("")
	}

	return m, nil
}

// downloadItem downloads an item
func (m Model) downloadItem(msg backend.DownloadItemMessage) (tea.Model, tea.Cmd) {
	m.msg = fmt.Sprintf("Saving item from feed %s", msg.Key)
	return m, m.config.Backend.DownloadItem(msg.Key, msg.Index)
}

// showHelp shows the help menu at the bottom of the screen
func (m Model) showHelp() (tea.Model, tea.Cmd) {
	// Extend the bindings with the tab specific bindings
	bindings := append(m.keymap.ShortHelp(), m.tabs[m.activeTab].GetKeyBinds()...)
	m.msg = m.help.ShortHelpView(bindings)
	return m, nil
}

// renderTabBar renders the tab bar at the top of the screen
func (m Model) renderTabBar() string {
	// Render the tab bar at the top of the screen
	tabs := make([]string, len(m.tabs))
	for i := range m.tabs {
		tabs[i] = m.style.attachIcon(m.tabs[i], m.tabs[i].Title(), i == m.activeTab)
	}

	// Check if the row exceeds the width of the screen
	if lipgloss.Width(strings.Join(tabs, "")) > m.width {
		tabs = tabs[m.activeTab:]
	}

	// Create the row
	row := strings.Join(tabs, "")

	// Calculate the gap amount
	var gapAmount int
	if m.width-lipgloss.Width(row) < 0 {
		gapAmount = 0
	} else {
		gapAmount = m.width - lipgloss.Width(row)
	}

	// Create the gap on the right
	gap := m.style.tabGap.Render(strings.Repeat(" ", gapAmount))
	return lipgloss.JoinHorizontal(lipgloss.Left, row, gap)
}

// renderStatusBar is used to render the status bar at the bottom of the screen
func (m Model) renderStatusBar() string {
	// Render the status bar at the bottom of the screen
	row := m.style.styleStatusBarCell(m.tabs[m.activeTab], m.offline)

	// Calculate the gap amount
	var gapAmount int
	if m.width-lipgloss.Width(row) < 0 {
		gapAmount = 0
	} else {
		gapAmount = m.width - lipgloss.Width(row)
	}

	// Render the gap on the right
	gap := m.style.statusBarGap.Render(strings.Repeat(" ", gapAmount))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
}
