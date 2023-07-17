package feed

import (
	"fmt"
	"strings"

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
	"mvdan.cc/xurls/v2"
)

// Keymap contains the key bindings for this tab
type Keymap struct {
	OpenArticle     key.Binding
	ToggleFocus     key.Binding
	RefreshArticles key.Binding
	SaveArticle     key.Binding
	DeleteFromSaved key.Binding
	CycleSelection  key.Binding
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
	CycleSelection: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "Cycle selection"),
	),
}

// ShortHelp returns the short help for the tab
func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.OpenArticle, k.ToggleFocus, k.RefreshArticles, k.SaveArticle, k.DeleteFromSaved, k.CycleSelection,
	}
}

// FullHelp returns the full help for the tab
func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.OpenArticle, k.ToggleFocus, k.RefreshArticles, k.SaveArticle, k.DeleteFromSaved, k.CycleSelection},
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
	errShown        bool
	list            list.Model
	articleContent  []string
	selCandidates   [][]int
	selIndex        int
	selActive       bool
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
	if !m.loaded {
		return m
	}

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
	case backend.FetchErrorMsg:
		// If the fetch failed, we need to display an error message
		m.errShown = true
		return m, nil

	case backend.FetchArticleSuccessMsg:
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

		case key.Matches(msg, m.keymap.CycleSelection):
			// TODO: We do not display styles if we are selecting, this is because
			// selection doesn't play well with line breaks. This is a problem for
			// future me.
			return m.cycleSelection(m.articleContent[m.list.Index()])
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
func (m Model) loadTab(items []list.Item, articleContents []string) (tab.Tab, tea.Cmd) {
	// Create the list
	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.ShowDescription = true
	itemDelegate.Styles = m.style.listItems
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
		glamour.WithStyles(m.colors.MarkdownStyle),
		glamour.WithWordWrap(m.style.viewportWidth),
	)

	if err != nil {
		m.errShown = true
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
	content := m.articleContent[m.list.Index()]
	text, err := m.termRenderer.Render(content)
	if err != nil {
		m.viewport.SetContent(fmt.Sprintf("We have encountered an error styling the content: %s", err))
		return m, nil
	}

	// Find all the selectable URLs
	m.selCandidates = xurls.Strict().FindAllStringIndex(content, -1)
	m.selIndex = 0
	m.selActive = false

	// Set the content of the viewport
	m.viewport.SetContent(text)
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
	var loadingMsg string

	if m.errShown {
		loadingMsg = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.style.errIcon,
			m.style.loadingMsg.Render("Failed to load the tab"),
		)
	} else {
		loadingMsg = m.style.loadingMsg.Render(
			fmt.Sprintf("%s Loading feed %s", m.loadingSpinner.View(), m.title),
		)
	}

	return loadingMsg
}

// cycleSelection highlights the link in the viewport
func (m Model) cycleSelection(content string) (tab.Tab, tea.Cmd) {
	start, end := m.selCandidates[m.selIndex][0], m.selCandidates[m.selIndex][1]
	ogOffset := m.viewport.YOffset

	var b strings.Builder
	b.WriteString(content[:start])
	b.WriteString(m.style.link.Render(content[start:end]))
	b.WriteString(fmt.Sprintf("%d", len(content[start:end])))
	b.WriteString(content[end:])

	m.viewport.SetContent(b.String())
	m.viewport.SetYOffset(ogOffset)
	m.selActive = true
	m.selIndex = (m.selIndex + 1) % len(m.selCandidates)

	return m, nil
}
