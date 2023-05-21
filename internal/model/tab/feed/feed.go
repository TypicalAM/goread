package feed

import (
	"fmt"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/model/simplelist"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model contains the state of this tab
type Model struct {
	colors          colorscheme.Colorscheme
	style           style
	width           int
	height          int
	title           string
	loaded          bool
	loadingSpinner  spinner.Model
	fetchFailed     bool
	list            list.Model
	isViewportOpen  bool
	viewport        viewport.Model
	viewportFocused bool
	keymap          keymap
	help            help.Model

	// reader is a function which returns a tea.Cmd which will be executed
	// when the tab is initialized
	reader func(string) tea.Cmd
}

type keymap struct {
	CloseTab    key.Binding
	CycleTabs   key.Binding
	Open        key.Binding
	ToggleFocus key.Binding
	Refresh     key.Binding
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.CloseTab, k.CycleTabs, k.Open, k.ToggleFocus, k.Refresh,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CloseTab, k.CycleTabs, k.Open, k.ToggleFocus, k.Refresh},
	}
}

// New creates a new feed tab with sensible defaults
func New(colors colorscheme.Colorscheme, width, height int, title string, reader func(string) tea.Cmd) Model {
	// Create a spinner for loading the data
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Foreground(colors.Color1)

	help := help.New()
	help.Styles.ShortDesc = lipgloss.NewStyle().Foreground(colors.Text)
	help.Styles.ShortKey = lipgloss.NewStyle().Foreground(colors.Text)
	help.Styles.Ellipsis = lipgloss.NewStyle().Foreground(colors.BgDark)

	// Create the model
	return Model{
		colors:         colors,
		style:          newStyle(colors),
		width:          width,
		height:         height,
		loadingSpinner: spin,
		title:          title,
		reader:         reader,
		help:           help,
		keymap: keymap{
			CloseTab: key.NewBinding(
				key.WithKeys("ctrl+w"),
				key.WithHelp("ctrl+w", "Close tab"),
			),
			CycleTabs: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "Cycle tabs"),
			),
			Open: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "Open"),
			),
			ToggleFocus: key.NewBinding(
				key.WithKeys("left", "right"),
				key.WithHelp("left/right", "Toggle focus"),
			),
			Refresh: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "Refresh"),
			),
		},
	}
}

// Title returns the title of the tab
func (m Model) Title() string {
	return m.title
}

// Type returns the type of the tab
func (m Model) Type() tab.Type {
	return tab.Feed
}

// SetSize sets the dimensions of the tab
func (m Model) SetSize(width, height int) tab.Tab {
	listWidth := width / 4
	viewportWidth := width - listWidth - 2
	m.list.SetSize(listWidth, height)
	m.viewport.Width = viewportWidth
	m.viewport.Height = height
	m.width = width
	m.height = height
	newTab, _ := m.updateViewport()
	return newTab
}

func (m Model) ShowHelp() string {
	return m.help.View(m.keymap)
}

// Init initializes the tab
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.reader(m.title), m.loadingSpinner.Tick)
}

// Update the variables of the tab
func (m Model) Update(msg tea.Msg) (tab.Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case backend.FetchErrorMessage:
		// If the fetch failed, we need to display an error message
		m.fetchFailed = true
		return m, nil

	case backend.FetchSuccessMessage:
		// If the fetch succeeded, we need to load the tab
		return m.loadTab(msg.Items)

	case tea.KeyMsg:
		// If the tab is not loaded, return
		if !m.loaded {
			return m, nil
		}

		// Handle the key message
		switch msg.String() {
		case "enter":
			// Update the viewport
			return m.updateViewport()

		case "r":
			// Refresh the contents of the tab
			m.isViewportOpen = false
			m.loaded = false
			m.viewportFocused = false

			// Rerun with data fetching and loading
			return m, tea.Batch(m.reader(m.title), m.loadingSpinner.Tick)

		case "left", "right":
			// If the viewport isn't open, don't do anything
			if !m.isViewportOpen {
				return m, nil
			}

			// Toggle the viewport focus
			m.viewportFocused = !m.viewportFocused
			return m, nil
		}

	default:
		// If the model is not loaded, update the loading spinner
		if !m.loaded {
			var cmd tea.Cmd
			m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
			return m, cmd
		}
	}

	// Update the selected item from the pane
	var cmd tea.Cmd
	if m.viewportFocused {
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	// Prevent the list from updating if we are not loaded yet
	if m.loaded {
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	// Return no commands
	return m, nil
}

// loadTab is fired when the items are retrieved from the backend, it
// initializes the list and the viewport
func (m Model) loadTab(items []list.Item) (tab.Tab, tea.Cmd) {
	// Set the width and the height of the components
	listWidth := m.width / 4
	viewportWidth := m.width - listWidth - 2

	// Check if the items are simplelist.Item compliant
	wrappedItems := make([]list.Item, len(items))
	for i, item := range items {
		if _, ok := item.(simplelist.Item); ok {
			wrappedItems[i] = item.(simplelist.Item).WrapDescription(m.width/4 - 3)
		}
	}

	// Create the styles for the list items
	delegateStyles := list.NewDefaultItemStyles()
	delegateStyles.SelectedTitle = delegateStyles.SelectedTitle.Copy().
		BorderForeground(m.colors.Color3).
		Foreground(m.colors.Color3).
		Italic(true)

	delegateStyles.SelectedDesc = delegateStyles.SelectedDesc.Copy().
		BorderForeground(m.colors.Color3).
		Foreground(m.colors.Color2).
		Height(2).
		Italic(true)

	delegateStyles.NormalDesc = delegateStyles.NormalDesc.Copy().Height(2)

	// Create the list
	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.ShowDescription = true
	itemDelegate.Styles = delegateStyles
	itemDelegate.SetHeight(3)

	// Initialize the list
	m.list = list.New(wrappedItems, itemDelegate, listWidth, m.height)

	// Set some attributes for the list
	m.list.SetShowHelp(false)
	m.list.SetShowTitle(false)
	m.list.SetShowStatusBar(false)
	m.list.DisableQuitKeybindings()

	// Initialize the viewport
	m.viewport = viewport.New(viewportWidth, m.height)

	// We are locked and loaded
	m.loaded = true
	return m, nil
}

// updateViewport is fired when the user presses enter, it updates the
// viewport with the selected item
func (m Model) updateViewport() (tab.Tab, tea.Cmd) {
	// Set the view as open if it isn't
	if !m.isViewportOpen {
		m.isViewportOpen = true
	}

	// Set the width of the styled content for word wrapping
	contentWidth := m.width - m.width/4 - 2

	// Get the content of the selected item
	content, err := m.list.SelectedItem().(simplelist.Item).StyleContent(contentWidth)
	if err != nil {
		m.viewport.SetContent(
			fmt.Sprintf("We have encountered an error styling the content: %s", err),
		)
		return m, nil
	}

	// Set the content of the viewport
	m.viewport.SetContent(content)
	m.viewport.SetYOffset(0)
	return m, nil
}

// View the tab
func (m Model) View() string {
	if !m.loaded {
		// Show the loading message
		return m.showLoading()
	}

	// If the view is not open show just the rss list
	if !m.isViewportOpen {
		return m.style.focusedStyle.Render(m.list.View())
	}

	// If the viewport is focused, render it with the focused style
	if m.viewportFocused {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.style.columnStyle.Render(m.list.View()),
			m.style.focusedStyle.Render(m.viewport.View()),
		)
	}

	// Otherwise render it with the default style
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.style.focusedStyle.Render(m.list.View()),
		m.style.columnStyle.Render(m.viewport.View()),
	)
}

// showLoading shows the loading message or the error message
func (m Model) showLoading() string {
	// The style of the message
	messageStyle := lipgloss.NewStyle().
		MarginLeft(3).
		MarginTop(1)

	var loadingMsg string
	if m.fetchFailed {
		// Render the failed message with a cross mark
		errorMsgStyle := messageStyle.Copy().
			Foreground(m.colors.Color4)
		loadingMsg = lipgloss.JoinHorizontal(
			lipgloss.Top,
			errorMsgStyle.Render(" "),
			messageStyle.Render("Failed to load the articles"),
		)
	} else {
		// Render the loading message with a spinner
		loadingMsg = messageStyle.Render(
			fmt.Sprintf("%s Loading feed %s", m.loadingSpinner.View(), m.title),
		)
	}

	// Render the message
	return loadingMsg
}
