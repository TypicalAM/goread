package feed

import (
	"fmt"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/model/simplelist"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/TypicalAM/goread/internal/style"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model contains the state of this tab
type Model struct {
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

	// reader is a function which returns a tea.Cmd which will be executed
	// when the tab is initialized
	reader func(string) tea.Cmd
}

// New creates a new feed tab with sensible defaults
func New(width, height int, title string, reader func(string) tea.Cmd) Model {
	// Create a spinner for loading the data
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Foreground(style.GlobalColorscheme.Color1)

	// Create the model
	return Model{
		width:          width,
		height:         height,
		loadingSpinner: spin,
		title:          title,
		reader:         reader,
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

// Help returns the help for the tab
func (m Model) Help() tab.Help {
	return tab.Help{
		tab.KeyBind{Key: "enter", Description: "Open"},
		tab.KeyBind{Key: "left/right", Description: "Toggle focus"},
	}
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
			return m, m.reader(m.title)

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

	// Create the styles for the list items
	delegateStyles := list.NewDefaultItemStyles()
	delegateStyles.SelectedTitle = delegateStyles.SelectedTitle.Copy().
		BorderForeground(style.GlobalColorscheme.Color3).
		Foreground(style.GlobalColorscheme.Color3).
		Italic(true)

	delegateStyles.SelectedDesc = delegateStyles.SelectedDesc.Copy().
		BorderForeground(style.GlobalColorscheme.Color3).
		Foreground(style.GlobalColorscheme.Color2).
		Italic(true)

	// Create the list
	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.ShowDescription = true
	itemDelegate.Styles = delegateStyles
	itemDelegate.SetHeight(3)

	// Initialize the list
	m.list = list.New(items, itemDelegate, listWidth, m.height)

	// Set some attributes for the list
	m.list.SetShowHelp(false)
	m.list.SetShowTitle(false)
	m.list.SetShowStatusBar(false)

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
		return style.FocusedStyle.Render(m.list.View())
	}

	// If the viewport is focused, render it with the focused style
	if m.viewportFocused {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			style.ColumnStyle.Render(m.list.View()),
			style.FocusedStyle.Render(m.viewport.View()),
		)
	}

	// Otherwise render it with the default style
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		style.FocusedStyle.Render(m.list.View()),
		style.ColumnStyle.Render(m.viewport.View()),
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
		// Render the failed message with an cross mark
		errorMsgStyle := messageStyle.Copy().
			Foreground(style.GlobalColorscheme.Color4)
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
