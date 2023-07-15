package feed

import (
	"fmt"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/TypicalAM/goread/internal/model/simplelist"
	"github.com/TypicalAM/goread/internal/model/tab"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
)

// Keymap contains the key bindings for this tab
type Keymap struct {
	OpenArticle     key.Binding
	ToggleFocus     key.Binding
	RefreshArticles key.Binding
	SaveArticle     key.Binding
	DeleteFromSaved key.Binding
}

// DefaultKeymap contains the default key bindings for this tab
var DefaultKeymap = Keymap{
	OpenArticle: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Open"),
	),
	ToggleFocus: key.NewBinding(
		key.WithKeys("left", "right", "h", "l"),
		key.WithHelp("←/→", "Move left/right"),
	),
	RefreshArticles: key.NewBinding(
		key.WithKeys("r", "ctrl+r"),
		key.WithHelp("r/C-r", "Refresh"),
	),
	SaveArticle: key.NewBinding(
		key.WithKeys("s", "ctrl+s"),
		key.WithHelp("s/C-s", "Save"),
	),
	DeleteFromSaved: key.NewBinding(
		key.WithKeys("d", "ctrl+d"),
		key.WithHelp("d/C-d", "Delete from saved"),
	),
}

// ShortHelp returns the short help for the tab
func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.OpenArticle, k.ToggleFocus, k.RefreshArticles, k.SaveArticle, k.DeleteFromSaved,
	}
}

// FullHelp returns the full help for the tab
func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.OpenArticle, k.ToggleFocus, k.RefreshArticles, k.SaveArticle, k.DeleteFromSaved},
	}
}

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
	articleContent  []string
	termRenderer    glamour.TermRenderer
	isViewportOpen  bool
	viewport        viewport.Model
	viewportFocused bool
	keymap          Keymap

	// reader is a function which returns a tea.Cmd which will be executed
	// when the tab is initialized
	reader func(string) tea.Cmd
}

// New creates a new feed tab with sensible defaults
func New(colors colorscheme.Colorscheme, width, height int, title string, reader func(string) tea.Cmd) Model {
	// Create a spinner for loading the data
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Foreground(colors.Color1)

	// Create the model
	return Model{
		colors:         colors,
		style:          newStyle(colors, width, height),
		width:          width,
		height:         height,
		loadingSpinner: spin,
		title:          title,
		reader:         reader,
		keymap:         DefaultKeymap,
	}
}

// Title returns the title of the tab
func (m Model) Title() string {
	return m.title
}

// Style returns the style of the tab
func (m Model) Style() tab.Style {
	return tab.Style{
		Color: m.colors.Color3,
		Icon:  "",
		Name:  "FEED",
	}
}

// SetSize sets the dimensions of the tab
func (m Model) SetSize(width, height int) tab.Tab {
	m.style = m.style.setSize(width, height)
	m.list.SetSize(m.style.listWidth, height)
	m.viewport.Width = m.style.viewportWidth
	m.viewport.Height = height
	m.width = width
	m.height = height
	newTab, _ := m.updateViewport()
	return newTab
}

// GetKeyBinds returns the key bindings of the tab
func (m Model) GetKeyBinds() []key.Binding {
	return m.keymap.ShortHelp()
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

	case backend.FetchArticleSuccessMessage:
		// If the fetch succeeded, we need to load the tab
		return m.loadTab(msg.Items, msg.ArticleContents)

	case tea.KeyMsg:
		// If the tab is not loaded, return
		if !m.loaded {
			return m, nil
		}

		// Handle the key message
		switch {
		case key.Matches(msg, m.keymap.OpenArticle):
			// If there are no items, don't do anything
			if m.list.SelectedItem() == nil {
				return m, nil
			}

			// Set the view as open if it isn't
			if !m.isViewportOpen {
				m.isViewportOpen = true
			}

			// Update the viewport
			return m.updateViewport()

		case key.Matches(msg, m.keymap.ToggleFocus):
			// If the viewport isn't open, don't do anything
			if !m.isViewportOpen {
				return m, nil
			}

			// Toggle the viewport focus
			m.viewportFocused = !m.viewportFocused
			return m, nil

		case key.Matches(msg, m.keymap.RefreshArticles):
			// Refresh the contents of the tab
			m.isViewportOpen = false
			m.loaded = false
			m.viewportFocused = false

			// Rerun with data fetching and loading
			return m, tea.Batch(m.reader(m.title), m.loadingSpinner.Tick)

		case key.Matches(msg, m.keymap.SaveArticle):
			// Tell the main model to download the item
			return m, backend.DownloadItem(m.title, m.list.Index())

		case key.Matches(msg, m.keymap.DeleteFromSaved):
			// Tell the main model to delete the item
			return m, backend.DeleteItem(m, fmt.Sprintf("%d", m.list.Index()))
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
// TODO: Wrap descs to m.width/4 - 3
func (m Model) loadTab(items []list.Item, articleContents []string) (tab.Tab, tea.Cmd) {
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

	delegateStyles.NormalDesc = delegateStyles.NormalDesc.Copy().
		Foreground(m.colors.TextDark).
		Height(2)

	// Create the list
	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.ShowDescription = true
	itemDelegate.Styles = delegateStyles
	itemDelegate.SetHeight(3)

	// Wrap the descs, it's better to do it upfront then to rely on the list pagination
	for i := range items {
		item := items[i].(list.DefaultItem)
		items[i] = simplelist.NewItem(item.Title(), wrap.String(item.Description(), m.style.listWidth-4))
	}

	// Initialize the list
	m.list = list.New(items, itemDelegate, m.style.listWidth, m.height)

	// Set some attributes for the list
	m.list.SetShowHelp(false)
	m.list.SetShowTitle(false)
	m.list.SetShowStatusBar(false)
	m.list.DisableQuitKeybindings()

	// Initialize the viewport
	m.viewport = viewport.New(m.style.viewportWidth, m.height)

	// Create the renderer for the viewport
	m.articleContent = articleContents
	termRenderer, err := glamour.NewTermRenderer(
		// TODO: Auto gen ansi.StyleConfig
		glamour.WithStylePath("dracula"),
		glamour.WithWordWrap(m.style.viewportWidth),
	)

	// TODO: Infinite loop?
	if err != nil {
		m.loaded = false
		return m, nil
	}

	// Locked and loaded
	m.termRenderer = *termRenderer
	m.loaded = true
	return m, nil
}

// updateViewport is fired when the user presses enter, it updates the
// viewport with the selected item
func (m Model) updateViewport() (tab.Tab, tea.Cmd) {
	// If the viewport isn't open, don't do anything
	if !m.isViewportOpen {
		return m, nil
	}

	// If therer are no items selected, don't do anything
	if m.list.SelectedItem() == nil {
		return m, nil
	}

	// Get the content of the selected item
	content, err := m.termRenderer.Render(m.articleContent[m.list.Index()])
	if err != nil {
		m.viewport.SetContent(fmt.Sprintf("We have encountered an error styling the content: %s", err))
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
		return m.style.focusedList.Render(m.list.View())
	}

	// If the viewport is focused, render it with the focused style
	if m.viewportFocused {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.style.idleList.Render(m.list.View()),
			m.style.focusedViewport.Render(m.viewport.View()),
		)
	}

	// Otherwise render it with the default style
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.style.focusedList.Render(m.list.View()),
		m.style.idleViewport.Render(m.viewport.View()),
	)
}

// DisableSaving disables the saving of the article
func (m Model) DisableSaving() Model {
	m.keymap.SaveArticle.SetEnabled(false)
	return m
}

// DisableDeleting disables the deleting of the article
func (m Model) DisableDeleting() Model {
	m.keymap.DeleteFromSaved.SetEnabled(false)
	return m
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
