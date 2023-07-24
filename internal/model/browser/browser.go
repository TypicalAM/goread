package browser

import (
	"fmt"
	"log"
	"strings"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/backend/rss"
	"github.com/TypicalAM/goread/internal/model/popup"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/TypicalAM/goread/internal/model/tab/category"
	"github.com/TypicalAM/goread/internal/model/tab/feed"
	"github.com/TypicalAM/goread/internal/model/tab/welcome"
	"github.com/TypicalAM/goread/internal/theme"

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

// SetEnabled allows to disable/enable shortcuts
func (k Keymap) SetEnabled(enabled bool) Keymap {
	k.CloseTab.SetEnabled(enabled)
	k.CycleTabs.SetEnabled(enabled)
	k.ShowHelp.SetEnabled(enabled)
	k.ToggleOfflineMode.SetEnabled(enabled)
	return k
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
	popup          popup.Popup
	backend        *backend.Backend
	help           help.Model
	style          style
	msg            string
	keymap         Keymap
	tabs           []tab.Tab
	activeTab      int
	height         int
	width          int
	waitingForSize bool
	quitting       bool
	offline        bool
}

// New returns a new model with some sensible defaults
func New(colors *theme.Colors, backend *backend.Backend) Model {
	log.Println("Initializing the browser")
	help := help.New()
	help.Styles.ShortDesc = lipgloss.NewStyle().Foreground(colors.Text)
	help.Styles.ShortKey = lipgloss.NewStyle().Foreground(colors.Text)
	help.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(colors.TextDark)
	help.ShortSeparator = " - "

	return Model{
		style:          newStyle(colors),
		backend:        backend,
		waitingForSize: true,
		keymap:         DefaultKeymap,
		help:           help,
		msg:            "Pro-tip - press [ctrl-h] to view the help page",
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles the terminal size, modifying rss items and modifying tabs
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.waitingForSize {
		return m.waitForSize(msg)
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case backend.FetchErrorMsg:
		// Update the underlying tab in case it also handles error input
		log.Printf("Error fetching data in tab %d: %v \n", m.activeTab, msg.Err)
		m.tabs[m.activeTab], _ = m.tabs[m.activeTab].Update(msg)
		m.msg = fmt.Sprintf("%s: %s", msg.Description, msg.Err.Error())
		return m, nil

	case category.ChosenCategoryMsg:
		m.popup = nil
		m.keymap = m.keymap.SetEnabled(true)

		if msg.IsEdit {
			if err := m.backend.Rss.UpdateCategory(msg.OldName, msg.Name, msg.Desc); err != nil {
				m.msg = fmt.Sprintf("Error updating category: %s", err.Error())
			} else {
				m.msg = fmt.Sprintf("Updated category %s", msg.Name)
			}
		} else {
			if err := m.backend.Rss.AddCategory(msg.Name, msg.Desc); err != nil {
				m.msg = fmt.Sprintf("Error adding category: %s", err.Error())
			} else {
				m.msg = fmt.Sprintf("Added category %s", msg.Name)
			}
		}

		log.Println(m.msg)
		return m, m.backend.FetchCategories("")

	case feed.ChosenFeedMsg:
		m.popup = nil
		m.keymap = m.keymap.SetEnabled(true)

		if msg.IsEdit {
			if err := m.backend.Rss.UpdateFeed(msg.ParentCategory, msg.OldName, msg.Name, msg.URL); err != nil {
				m.msg = fmt.Sprintf("Error updating feed: %s", err.Error())
			} else {
				m.msg = fmt.Sprintf("Updated feed %s", msg.Name)
			}
		} else {
			if err := m.backend.Rss.AddFeed(msg.ParentCategory, msg.Name, msg.URL); err != nil {
				m.msg = fmt.Sprintf("Error adding feed: %s", err.Error())
			} else {
				m.msg = fmt.Sprintf("Added feed %s", msg.Name)
			}
		}

		log.Println(m.msg)
		return m, m.backend.FetchFeeds(msg.ParentCategory)

	case tab.NewTabMessage:
		return m.createNewTab(msg)

	case backend.NewItemMsg:
		bg := m.View()
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

		m.keymap = m.keymap.SetEnabled(false)
		return m, m.popup.Init()

	case backend.DeleteItemMsg:
		return m.deleteItem(msg)

	case backend.DownloadItemMsg:
		return m.downloadItem(msg)

	case backend.MakeChoiceMsg:
		bg := m.View()
		width := m.width / 2
		m.popup = popup.NewChoice(m.style.colors, bg, width, msg.Question, msg.Default)

		m.keymap = m.keymap.SetEnabled(false)
		return m, m.popup.Init()

	case popup.ChoiceResultMsg:
		m.keymap = m.keymap.SetEnabled(true)
		m.popup = nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.msg = ""

		for i := range m.tabs {
			m.tabs[i] = m.tabs[i].SetSize(m.width, m.height-5)
		}

	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case msg.String() == "esc":
			// If we are showing a popup, close it
			if m.popup != nil {
				m.keymap = m.keymap.SetEnabled(true)
				m.popup = nil
				return m, nil
			}

			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keymap.CloseTab):
			if len(m.tabs) == 1 {
				m.quitting = true
				return m, tea.Quit
			}

			// Close the current tab
			m.tabs = append(m.tabs[:m.activeTab], m.tabs[m.activeTab+1:]...)
			m.activeTab--

			if m.activeTab < 0 {
				m.activeTab = 0
			}

			m.msg = fmt.Sprintf("Closed tab - %s", m.tabs[m.activeTab].Title())
			return m, nil

		case key.Matches(msg, m.keymap.CycleTabs):
			m.activeTab++
			if m.activeTab > len(m.tabs)-1 {
				m.activeTab = 0
			}

			m.msg = ""
			return m, nil

		case key.Matches(msg, m.keymap.ShowHelp):
			return m.showHelp()

		case key.Matches(msg, m.keymap.ToggleOfflineMode):
			return m.toggleOffline()
		}
	}

	// If we are showing a popup, we need to update the popup
	if m.popup != nil {
		m.popup, cmd = m.popup.Update(msg)
		return m, cmd
	}

	m.tabs[m.activeTab], cmd = m.tabs[m.activeTab].Update(msg)
	return m, cmd
}

// View renders the tab bar, the active tab and the status bar
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!"
	}

	if m.waitingForSize {
		return "Loading..."
	}

	if m.popup != nil {
		return m.popup.View()
	}

	// TODO: refactor
	var sections []string
	sections = append(sections, m.renderTabBar())
	constrainHeight := lipgloss.NewStyle().Height(m.height - 3).MaxHeight(m.height - 3)
	sections = append(sections, constrainHeight.Render(m.tabs[m.activeTab].View()))
	sections = append(sections, m.renderStatusBar())

	if strings.Contains(m.msg, "Error") {
		sections = append(sections, m.style.errMsg.Render(m.msg))
	} else {
		sections = append(sections, m.msg)
	}

	return strings.Join(sections, "\n")
}

// waitForSize waits for the window size to be set and loads the tab
func (m Model) waitForSize(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.WindowSizeMsg); !ok {
		return m, nil
	}

	sizeMsg := msg.(tea.WindowSizeMsg)
	m.width = sizeMsg.Width
	m.height = sizeMsg.Height
	m.waitingForSize = false

	m.tabs = append(m.tabs, welcome.New(
		m.style.colors,
		m.width,
		m.height-5,
		"Welcome",
		m.backend.FetchCategories,
	))

	return m, m.tabs[0].Init()
}

// createNewTab bootstraps the new tab and adds it to the model
func (m Model) createNewTab(msg tab.NewTabMessage) (Model, tea.Cmd) {
	var newTab tab.Tab
	height := m.height - 5

	switch msg.Sender.(type) {
	case welcome.Model:
		switch msg.Title {
		case rss.AllFeedsName:
			newTab = feed.New(m.style.colors, m.width, height, msg.Title, m.backend.FetchAllArticles).
				DisableSaving().
				DisableDeleting()

		case rss.DownloadedFeedsName:
			newTab = feed.New(m.style.colors, m.width, height, msg.Title, m.backend.FetchDownloadedArticles).
				DisableSaving()

		default:
			newTab = category.New(m.style.colors, m.width, height, msg.Title, m.backend.FetchFeeds)
		}

	case category.Model:
		newTab = feed.New(m.style.colors, m.width, height, msg.Title, m.backend.FetchArticles).
			DisableDeleting()
	}

	// Insert the tab after the active tab
	m.tabs = append(m.tabs[:m.activeTab+1], append([]tab.Tab{newTab}, m.tabs[m.activeTab+1:]...)...)
	m.activeTab++
	m.msg = ""

	return m, newTab.Init()
}

// deleteItem deletes the focused item from the backend
func (m Model) deleteItem(msg backend.DeleteItemMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.msg = fmt.Sprintf("Deleting item %s", msg.Key)

	// Check the type of the item
	switch msg.Sender.(type) {
	case welcome.Model:
		cmd = m.backend.FetchCategories("")
		if err := m.backend.Rss.RemoveCategory(msg.Key); err != nil {
			m.msg = fmt.Sprintf("Error deleting category %s: %s", msg.Key, err.Error())
		}

	case category.Model:
		cmd = m.backend.FetchFeeds(m.tabs[m.activeTab].Title())
		if err := m.backend.Rss.RemoveFeed(m.tabs[m.activeTab].Title(), msg.Key); err != nil {
			m.msg = fmt.Sprintf("Error deleting feed %s: %s", msg.Key, err.Error())
		}

	case feed.Model:
		cmd = m.backend.FetchDownloadedArticles("")
		if msg.Sender.Title() == rss.DownloadedFeedsName {
			if err := m.backend.RemoveDownload(msg.Key); err != nil {
				m.msg = fmt.Sprintf("Error deleting download %s: %s", msg.Key, err.Error())
			}
		}
	}

	log.Println(m.msg)
	return m, cmd
}

// downloadItem downloads an item
func (m Model) downloadItem(msg backend.DownloadItemMsg) (tea.Model, tea.Cmd) {
	log.Println("Downloading item", msg.Key, msg.Index)
	m.msg = "Item saved! You can find it in the downloaded category"
	return m, m.backend.DownloadItem(msg.Key, msg.Index)
}

// showHelp shows the help menu at the bottom of the screen
func (m Model) showHelp() (tea.Model, tea.Cmd) {
	// Extend the bindings with the tab specific bindings
	bindings := append(m.keymap.ShortHelp(), m.tabs[m.activeTab].GetKeyBinds()...)
	m.msg = m.help.ShortHelpView(bindings)
	log.Println("Showing help with n entries: ", len(bindings))
	return m, nil
}

// toggleOffline toggles the offline mode
func (m Model) toggleOffline() (tea.Model, tea.Cmd) {
	m.offline = !m.offline
	m.backend.SetOfflineMode(m.offline)

	if m.offline {
		m.msg = "Offline mode enabled"
	} else {
		m.msg = "Offline mode disabled"
	}

	log.Println(m.msg)
	return m, nil
}

// renderTabBar renders the tab bar at the top of the screen
func (m Model) renderTabBar() string {
	tabs := make([]string, len(m.tabs))
	for i := range m.tabs {
		tabs[i] = m.style.attachIcon(m.tabs[i], m.tabs[i].Title(), i == m.activeTab)
	}

	if lipgloss.Width(strings.Join(tabs, "")) > m.width {
		tabs = tabs[m.activeTab:]
	}

	row := strings.Join(tabs, "")

	var gapAmount int
	if m.width-lipgloss.Width(row) < 0 {
		gapAmount = 0
	} else {
		gapAmount = m.width - lipgloss.Width(row)
	}

	gap := m.style.tabGap.Render(strings.Repeat(" ", gapAmount))
	return lipgloss.JoinHorizontal(lipgloss.Left, row, gap)
}

// renderStatusBar is used to render the status bar at the bottom of the screen
func (m Model) renderStatusBar() string {
	row := m.style.styleStatusBarCell(m.tabs[m.activeTab], m.offline)

	var gapAmount int
	if m.width-lipgloss.Width(row) < 0 {
		gapAmount = 0
	} else {
		gapAmount = m.width - lipgloss.Width(row)
	}

	gap := m.style.statusBarGap.Render(strings.Repeat(" ", gapAmount))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
}
