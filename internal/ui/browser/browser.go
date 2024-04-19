package browser

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/backend/rss"
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/TypicalAM/goread/internal/ui/popup"
	"github.com/TypicalAM/goread/internal/ui/tab"
	"github.com/TypicalAM/goread/internal/ui/tab/category"
	"github.com/TypicalAM/goread/internal/ui/tab/feed"
	"github.com/TypicalAM/goread/internal/ui/tab/overview"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Keymap contains the key bindings for the browser
type Keymap struct {
	CloseTab          key.Binding
	NextTab           key.Binding
	PrevTab           key.Binding
	ShowHelp          key.Binding
	ToggleOfflineMode key.Binding
}

// DefaultKeymap contains the default key bindings for the browser
var DefaultKeymap = Keymap{
	CloseTab: key.NewBinding(
		key.WithKeys("c", "ctrl+w"),
		key.WithHelp("c", "Close tab"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "Next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("Shift+Tab", "Previous tab"),
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

// SetEnabled allows to disable/enable shortcuts
func (k *Keymap) SetEnabled(enabled bool) {
	k.CloseTab.SetEnabled(enabled)
	k.NextTab.SetEnabled(enabled)
	k.PrevTab.SetEnabled(enabled)
	k.ShowHelp.SetEnabled(enabled)
	k.ToggleOfflineMode.SetEnabled(enabled)
}

// Model is used to store the state of the application
type Model struct {
	popup          popup.Window
	overlay        popup.Overlay
	backend        *backend.Backend
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

	return Model{
		style:          newStyle(colors),
		backend:        backend,
		waitingForSize: true,
		keymap:         DefaultKeymap,
		msg:            "Pro-tip - press [ctrl+h] to view the help page",
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
	case backend.StartQuittingMsg:
		m.quitting = true
		return m, tea.Quit

	case backend.FetchErrorMsg:
		// Update the underlying tab in case it also handles error input
		log.Printf("Error fetching data in tab %d: %v \n", m.activeTab, msg.Err)
		updated, _ := m.tabs[m.activeTab].Update(msg)
		m.tabs[m.activeTab] = updated.(tab.Tab)
		m.msg = fmt.Sprintf("%s: %s", msg.Description, unwrapErrs(msg.Err))
		return m, nil

	case overview.ChosenCategoryMsg:
		m.popup = nil
		m.keymap.SetEnabled(true)

		if msg.IsEdit {
			if err := m.backend.Rss.UpdateCategory(msg.OldName, msg.Name, msg.Desc); err != nil {
				m.msg = fmt.Sprintf("Error updating category: %s", unwrapErrs(err))
			} else {
				m.msg = fmt.Sprintf("Updated category %s", msg.Name)
			}
		} else {
			if err := m.backend.Rss.AddCategory(msg.Name, msg.Desc); err != nil {
				m.msg = fmt.Sprintf("Error adding category: %s", unwrapErrs(err))
			} else {
				m.msg = fmt.Sprintf("Added category %s", msg.Name)
			}
		}

		return m, m.backend.FetchCategories("")

	case category.ChosenFeedMsg:
		m.popup = nil
		m.keymap.SetEnabled(true)

		if msg.IsEdit {
			if err := m.backend.Rss.UpdateFeed(msg.Parent, msg.OldName, msg.Name, msg.URL); err != nil {
				m.msg = fmt.Sprintf("Error updating feed: %s", unwrapErrs(err))
			} else {
				m.msg = fmt.Sprintf("Updated feed %s", msg.Name)
			}
		} else {
			if err := m.backend.Rss.AddFeed(msg.Parent, msg.Name, msg.URL); err != nil {
				m.msg = fmt.Sprintf("Error adding feed: %s", unwrapErrs(err))
			} else {
				m.msg = fmt.Sprintf("Added feed %s", msg.Name)
			}
		}

		return m, m.backend.FetchFeeds(msg.Parent)

	case tab.NewTabMsg:
		return m.createNewTab(msg)

	case backend.NewItemMsg:
		background := m.View()

		switch msg.Sender.(type) {
		case overview.Model:
			m.popup = overview.NewPopup(m.style.colors, "", "")
		case category.Model:
			m.popup = category.NewPopup(m.style.colors, "", "", msg.Sender.Title())
		case feed.Model:
		}

		width, height := m.popup.GetSize()
		m.overlay = popup.NewOverlay(background, width, height)
		m.keymap.SetEnabled(false)
		return m, m.popup.Init()

	case backend.EditItemMsg:
		oldName, oldDesc := msg.OldFields[0], msg.OldFields[1]

		switch msg.Sender.(type) {
		case overview.Model:
			m.popup = overview.NewPopup(m.style.colors, oldName, oldDesc)
		case category.Model:
			m.popup = category.NewPopup(m.style.colors, oldName, oldDesc, msg.Sender.Title())
		case feed.Model:
		}

		m.keymap.SetEnabled(false)
		return m, m.popup.Init()

	case backend.DeleteItemMsg:
		return m.deleteItem(msg)

	case backend.DownloadItemMsg:
		return m.downloadItem(msg)

	case backend.MarkAsReadMsg:
		m.backend.ReadStatus.MarkAsRead(string(msg))
		return m, nil

	case backend.MarkAsUnreadMsg:
		m.backend.ReadStatus.MarkAsUnread(string(msg))
		return m, nil

	case backend.MakeChoiceMsg:
		background := m.View()
		m.popup = popup.NewChoice(m.style.colors, msg.Question, msg.Default)
		width, height := m.popup.GetSize()
		m.overlay = popup.NewOverlay(background, width, height)
		m.keymap.SetEnabled(false)
		return m, m.popup.Init()

	case popup.ChoiceResultMsg:
		m.keymap.SetEnabled(true)
		m.popup = nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.msg = ""

		for i := range m.tabs {
			m.tabs[i] = m.tabs[i].SetSize(m.width, m.height-5)
		}

		// Delete the popup, update the overlay and rerender
		if m.popup != nil {
			popupBackup := m.popup
			m.popup = nil
			background := m.View()
			m.popup = popupBackup
			width, height := m.popup.GetSize()
			m.overlay = popup.NewOverlay(background, width, height)
		}

	case backend.SetEnableKeybindMsg:
		m.keymap.SetEnabled(bool(msg))
		log.Println("Disabling keybinds, propagating")

	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case msg.String() == "esc":
			// If we are showing a popup, close it. We leave esc handling to the model.
			if m.popup != nil {
				m.keymap.SetEnabled(true)
				m.popup = nil
				return m, nil
			}

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

		case key.Matches(msg, m.keymap.NextTab):
			m.activeTab++
			if m.activeTab > len(m.tabs)-1 {
				m.activeTab = 0
			}

			m.msg = ""
			return m, nil

		case key.Matches(msg, m.keymap.PrevTab):
			m.activeTab--
			if m.activeTab < 0 {
				m.activeTab = len(m.tabs) - 1
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
		newPopup, cmd := m.popup.Update(msg)
		m.popup = newPopup.(popup.Window)
		return m, cmd
	}

	updated, cmd := m.tabs[m.activeTab].Update(msg)
	m.tabs[m.activeTab] = updated.(tab.Tab)
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
		return m.overlay.WrapView(m.popup.View())
	}

	var b strings.Builder
	b.WriteString(m.renderTabBar())
	b.WriteRune('\n')
	constrainHeight := lipgloss.NewStyle().Height(m.height - 3).MaxHeight(m.height - 3)
	b.WriteString(constrainHeight.Render(m.tabs[m.activeTab].View()))
	b.WriteRune('\n')
	b.WriteString(m.renderStatusBar())
	b.WriteRune('\n')

	if strings.Contains(m.msg, "Error") {
		b.WriteString(m.style.errMsg.Render(m.msg))
	} else {
		b.WriteString(m.msg)
	}

	return b.String()
}

// ShortHelp returns the short help for the browser.
func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{m.keymap.CloseTab, m.keymap.NextTab, m.keymap.PrevTab, m.keymap.ToggleOfflineMode}
}

// FullHelp returns the full help for the browser.
func (m Model) FullHelp() [][]key.Binding {
	browserHelp := [][]key.Binding{m.ShortHelp()}
	childHelp := m.tabs[m.activeTab].FullHelp()
	return prettifyHelp(browserHelp, childHelp)
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

	m.tabs = append(m.tabs, overview.New(
		m.style.colors,
		m.width,
		m.height-5,
		"Welcome",
		m.backend.FetchCategories,
	))

	return m, m.tabs[0].Init()
}

// createNewTab bootstraps the new tab and adds it to the model
func (m Model) createNewTab(msg tab.NewTabMsg) (Model, tea.Cmd) {
	var newTab tab.Tab
	height := m.height - 5

	switch msg.Sender.(type) {
	case overview.Model:
		switch msg.Title {
		case rss.AllFeedsName:
			newTab = feed.New(m.style.colors, m.width, height, msg.Title, m.backend.FetchAllArticles).
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
	m.msg = fmt.Sprintf("Deleting item %s", msg.ItemName)

	// Check the type of the item
	switch msg.Sender.(type) {
	case overview.Model:
		cmd = m.backend.FetchCategories("")
		if err := m.backend.Rss.RemoveCategory(msg.ItemName); err != nil {
			m.msg = fmt.Sprintf("Error deleting category %s: %s", msg.ItemName, unwrapErrs(err))
		}

	case category.Model:
		cmd = m.backend.FetchFeeds(m.tabs[m.activeTab].Title())
		if err := m.backend.Rss.RemoveFeed(m.tabs[m.activeTab].Title(), msg.ItemName); err != nil {
			m.msg = fmt.Sprintf("Error deleting feed %s: %s", msg.ItemName, unwrapErrs(err))
		}

	case feed.Model:
		cmd = m.backend.FetchDownloadedArticles("", false)
		if msg.Sender.Title() == rss.DownloadedFeedsName {
			index, err := strconv.Atoi(msg.ItemName)
			if err != nil {
				m.msg = fmt.Sprintf("Error deleting download %s: %s", msg.ItemName, unwrapErrs(err))
			}

			if err := m.backend.Cache.RemoveFromDownloaded(index); err != nil {
				m.msg = fmt.Sprintf("Error deleting download %s: %s", msg.ItemName, unwrapErrs(err))
			}
		}
	}

	log.Println(m.msg)
	return m, cmd
}

// downloadItem downloads an item
func (m Model) downloadItem(msg backend.DownloadItemMsg) (tea.Model, tea.Cmd) {
	log.Println("Downloading item", msg.FeedName, msg.Index)
	m.msg = "Item saved! You can find it in the downloaded category"
	return m, m.backend.DownloadItem(msg.FeedName, msg.Index)
}

// showHelp shows the help menu as a popup.
func (m Model) showHelp() (tea.Model, tea.Cmd) {
	background := m.View()
	m.popup = newHelp(m.style.colors, m.FullHelp())
	width, height := m.popup.GetSize()
	m.overlay = popup.NewOverlay(background, width, height)
	m.keymap.SetEnabled(false)
	return m, nil
}

// toggleOffline toggles the offline mode
func (m Model) toggleOffline() (tea.Model, tea.Cmd) {
	m.offline = !m.offline
	m.backend.Cache.OfflineMode = m.offline

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

// unwrapErrs unwraps all errors in a chain of wrapped errors for use in a status message
func unwrapErrs(err error) error {
	for {
		unwrapErr := errors.Unwrap(err)
		if unwrapErr == nil {
			break
		}
		err = unwrapErr
	}
	return err
}

// prettifyHelp prettifies the help columns, removes the lower level bind if a higher one precedes it
func prettifyHelp(first, second [][]key.Binding) [][]key.Binding {
	toDelSecond := make([]int, 0)
	toDelThird := make([]int, 0)

	for _, elem := range first[0] {
		// Second col
		for idx2, elem2 := range second[0] {
			if strIntersect(elem.Keys(), elem2.Keys()) {
				toDelSecond = append(toDelSecond, idx2)
			}
		}

		// Third col
		for idx2, elem2 := range second[1] {
			if strIntersect(elem.Keys(), elem2.Keys()) {
				toDelThird = append(toDelThird, idx2)
			}
		}
	}

	second[0] = deleteBinds(second[0], toDelSecond)
	second[1] = deleteBinds(second[1], toDelThird)

	for _, elem := range second[0] {
		// Third col
		for idx2, elem2 := range second[1] {
			if strIntersect(elem.Keys(), elem2.Keys()) {
				toDelThird = append(toDelThird, idx2)
			}
		}
	}

	second[1] = deleteBinds(second[1], toDelThird)
	return append(first, second...)
}

// deleteBinds removes items from a bind slice via a list of indices
func deleteBinds(arr []key.Binding, indices []int) []key.Binding {
	if len(indices) == 0 {
		return arr
	}

	result := make([]key.Binding, 0, len(arr)-len(indices))
	for idx, elem := range arr {
		add := true
		for _, toDel := range indices {
			if idx == toDel {
				add = false
				break
			}
		}

		if add {
			result = append(result, elem)
		}
	}

	return result
}

// strIntersect checks if two string slices have a non-empty intersection
func strIntersect(first, second []string) bool {
	for _, elem := range first {
		for _, elem2 := range second {
			if elem == elem2 {
				return true
			}
		}
	}

	return false
}
