package feed

import (
	"fmt"
	"log"
	"strings"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/theme"
	"github.com/TypicalAM/goread/internal/ui/popup"
	"github.com/TypicalAM/goread/internal/ui/tab"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
)

// Model contains the state of this tab
type Model struct {
	list            list.Model
	fetcher         backend.ArticleFetcher
	colorTr         *glamour.TermRenderer
	noColorTr       *glamour.TermRenderer
	colors          *theme.Colors
	selector        *selector
	title           string
	viewport        viewport.Model
	keymap          Keymap
	articleContent  []string
	spinner         spinner.Model
	style           style
	height          int
	width           int
	errShown        bool
	loaded          bool
	viewportOpen    bool
	viewportFocused bool
	lastFilterState list.FilterState
}

// New creates a new feed tab with sensible defaults
func New(colors *theme.Colors, width, height int, title string, fetcher backend.ArticleFetcher) Model {
	log.Println("Creating new feed tab with title", title)
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Foreground(colors.Color1)

	// Create the model
	return Model{
		colors:   colors,
		style:    newStyle(colors, width, height),
		width:    width,
		height:   height,
		selector: newSelector(colors),
		spinner:  spin,
		title:    title,
		fetcher:  fetcher,
		keymap:   DefaultKeymap,
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

// Init initializes the tab
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetcher(m.title, false))
}

// Update the variables of the tab
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Allow quitting when fetching failed
	if m.errShown {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
			return m, backend.StartQuitting()
		}
	}

	switch msg := msg.(type) {
	case backend.FetchErrorMsg:
		m.errShown = true
		return m, nil

	case backend.FetchArticleSuccessMsg:
		return m.loadTab(msg.Items), nil

	case backend.SetEnableKeybindMsg:
		m.keymap.SetEnabled(bool(msg))
		return m, nil

	case popup.ChoiceResultMsg:
		if !msg.Result {
			return m, nil
		}

		_ = m.selector.open()
		return m, nil

	case tea.KeyMsg:
		if !m.loaded {
			return m, nil
		}

		switch {
		case msg.String() == "esc":
			if m.list.FilterState() == list.Unfiltered {
				return m, backend.StartQuitting()
			}

			// There is no way to call `list.resetFiltering` since it's not exported
			m.list.SetFilteringEnabled(false)
			m.list.SetFilteringEnabled(true)
			return m, nil

		case key.Matches(msg, m.keymap.Open):
			if m.viewportFocused && m.selector.active {
				return m, backend.MakeChoice("Open in browser?", true)
			}

			if m.list.SelectedItem() == nil {
				return m, nil
			}

			if !m.viewportOpen {
				m.viewportOpen = true
			}

			return m.updateViewport()

		case key.Matches(msg, m.keymap.ToggleFocus):
			if !m.viewportOpen {
				return m, nil
			}

			m.viewportFocused = !m.viewportFocused
			return m, nil

		case key.Matches(msg, m.keymap.RefreshArticles):
			m.viewportOpen = false
			m.loaded = false
			m.viewportFocused = false

			return m, tea.Batch(m.spinner.Tick, m.fetcher(m.title, true))

		case key.Matches(msg, m.keymap.SaveArticle):
			return m, backend.DownloadItem(m.title, m.list.Index())

		case key.Matches(msg, m.keymap.DeleteFromSaved):
			return m, backend.DeleteItem(m, fmt.Sprintf("%d", m.list.Index()))

		case key.Matches(msg, m.keymap.MarkAsUnread):
			item := m.list.SelectedItem().(backend.ArticleItem)
			if strings.HasPrefix(item.ArtTitle, "✓ ") {
				item.ArtTitle = strings.Join(strings.Split(item.ArtTitle, " ")[1:], " ")
				m.list.SetItem(m.list.Index(), item)
			}

			return m, backend.MarkAsUnread(item.FeedURL)

		case key.Matches(msg, m.keymap.CycleSelection):
			if !m.viewportFocused {
				return m, nil
			}

			m.viewport.SetContent(m.selector.cycle())
			return m, nil
		}

	default:
		if !m.loaded {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if !m.loaded {
		return m, nil
	}

	var cmd tea.Cmd
	if m.viewportFocused {
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)
	if m.list.FilterState() == m.lastFilterState {
		return m, cmd
	}

	keysEnabled := m.list.FilterState() != list.Filtering
	m.lastFilterState = m.list.FilterState()
	return m, tea.Sequence(cmd, backend.SetEnableKeybind(keysEnabled))
}

// loadTab is fired when the items are retrieved from the backend
func (m Model) loadTab(items []list.Item) tab.Tab {
	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.ShowDescription = true
	itemDelegate.Styles = m.style.listItems
	itemDelegate.SetHeight(3)

	// Wrap the descs, it's better to do it upfront then to rely on the list pagination
	for i := range items {
		item := items[i].(backend.ArticleItem)
		item.Desc = wrap.String(item.Desc, m.style.listWidth-4)
		items[i] = item
	}

	m.list = list.New(items, itemDelegate, m.style.listWidth, m.height)

	m.list.SetShowHelp(false)
	m.list.SetShowTitle(false)
	m.list.SetShowStatusBar(false)
	m.list.DisableQuitKeybindings()
	m.list.KeyMap.NextPage.SetEnabled(false)
	m.list.KeyMap.PrevPage.SetEnabled(false)
	m.list.KeyMap.CloseFullHelp.SetEnabled(false)

	m.viewport = viewport.New(m.style.viewportWidth, m.height)

	colorTr, err := glamour.NewTermRenderer(
		glamour.WithStyles(m.colors.MarkdownStyle),
		glamour.WithWordWrap(m.style.viewportWidth-2),
	)

	if err != nil {
		m.errShown = true
		m.loaded = false
		return m
	}

	noColorTr, err := glamour.NewTermRenderer(
		glamour.WithStyles(glamour.NoTTYStyleConfig),
		glamour.WithWordWrap(m.style.viewportWidth-2),
	)

	if err != nil {
		m.errShown = true
		m.loaded = false
		return m
	}

	// Locked and loaded
	m.colorTr = colorTr
	m.noColorTr = noColorTr
	m.loaded = true
	return m
}

// updateViewport is fired when the user presses enter, it updates the
// viewport with the selected item
func (m Model) updateViewport() (tab.Tab, tea.Cmd) {
	if !m.viewportOpen {
		return m, nil
	}

	if m.list.SelectedItem() == nil {
		return m, nil
	}

	rawText := m.list.SelectedItem().(backend.ArticleItem).MarkdownContent
	styledText, err := m.colorTr.Render(rawText)
	if err != nil {
		m.viewport.SetContent(fmt.Sprintf("We have encountered an error styling the content: %s", err))
		return m, nil
	}

	noColorText, err := m.noColorTr.Render(rawText)
	if err != nil {
		m.viewport.SetContent(fmt.Sprintf("We have encountered an error styling the content: %s", err))
		return m, nil
	}

	m.selector.newArticle(&rawText, &noColorText)
	m.viewport.SetContent(styledText)
	m.viewport.SetYOffset(0)

	// Mark this item as read and prepend a ✓
	item := m.list.SelectedItem().(backend.ArticleItem)
	if !strings.HasPrefix(item.ArtTitle, "✓ ") {
		item.ArtTitle = "✓ " + item.ArtTitle
		m.list.SetItem(m.list.Index(), item)
	}
	return m, backend.MarkAsRead(item.FeedURL)
}

// View the tab
func (m Model) View() string {
	if !m.loaded {
		return m.showLoading()
	}

	if !m.viewportOpen {
		return m.style.focusedList.Render(m.list.View())
	}

	if m.viewportFocused {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.style.idleList.Render(m.list.View()),
			m.style.focusedViewport.Render(m.viewport.View()),
		)
	}

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

// ShortHelp returns the short help for the tab
func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keymap.Open, m.keymap.ToggleFocus, m.keymap.RefreshArticles,
		m.keymap.SaveArticle, m.keymap.DeleteFromSaved, m.keymap.CycleSelection,
		m.keymap.MarkAsUnread,
	}
}

// FullHelp returns the full help for the tab
func (m Model) FullHelp() [][]key.Binding {
	if !m.viewportOpen && m.viewportFocused {
		result := [][]key.Binding{m.ShortHelp()}
		result = append(result, m.list.FullHelp()...)
		return result
	}

	return [][]key.Binding{m.ShortHelp(), {
		m.viewport.KeyMap.PageDown,
		m.viewport.KeyMap.PageUp,
		m.viewport.KeyMap.HalfPageDown,
		m.viewport.KeyMap.HalfPageUp,
		m.viewport.KeyMap.Down,
		m.viewport.KeyMap.Up,
	}}
}

// showLoading shows the loading message or the error message
func (m Model) showLoading() string {
	if m.errShown {
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.style.errIcon,
			m.style.loadingMsg.Render("Failed to load the tab"),
		)
	}

	return m.style.loadingMsg.Render(
		fmt.Sprintf("%s Loading feed %s", m.spinner.View(), m.title),
	)
}
