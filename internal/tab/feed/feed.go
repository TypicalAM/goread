package feed

import (
	"fmt"
	"strings"

	"github.com/TypicalAM/goread/internal/backend"
	simpleList "github.com/TypicalAM/goread/internal/list"
	"github.com/TypicalAM/goread/internal/style"
	"github.com/TypicalAM/goread/internal/tab"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SelectedPane int

const (
	articlesList SelectedPane = iota
	articlesPreview
)

// RssFeedTab is a tab that displays a list of RSS feeds
type RssFeedTab struct {
	title  string
	loaded bool

	loadingSpinner spinner.Model
	list           list.Model
	isViewportOpen bool
	viewport       viewport.Model
	selected       SelectedPane
	failed         bool

	readerFunc func(string) tea.Cmd

	availableWidth  int
	availableHeight int
}

// New creates a new RssFeedTab with sensible defaults
func New(availableWidth, availableHeight int, title string, readerFunc func(string) tea.Cmd) RssFeedTab {
	// Create a spinner for loading the data
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Foreground(style.GlobalColorscheme.Color1)

	return RssFeedTab{
		availableWidth:  availableWidth,
		availableHeight: availableHeight,
		loadingSpinner:  spin,
		title:           title,
		readerFunc:      readerFunc,
	}
}

// Return the title of the tab
func (r RssFeedTab) Title() string {
	return r.title
}

// Initialize the tab
func (r RssFeedTab) Init() tea.Cmd {
	return tea.Batch(r.readerFunc(r.title), r.loadingSpinner.Tick)
}

// loadTab is fired when the items are retrieved and we know the
// dimensions of the window. It initializes the list and the viewport
func (r *RssFeedTab) loadTab(items []list.Item) {
	// Set the width and the height of the components
	listWidth := r.availableWidth / 4
	viewportWidth := r.availableHeight - listWidth - 4 // 4 is the padding

	// Get the default styles for the list items
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

	// Now that we know the width of the list we can wrap the descriptions
	// to match it
	for i := range items {
		items[i] = items[i].(simpleList.Item)
	}

	// Initialize the list
	r.list = list.New(items, itemDelegate, listWidth, r.availableHeight-5)

	// Set some attributes for the list
	r.list.SetShowHelp(false)
	r.list.SetShowTitle(false)
	r.list.SetShowStatusBar(false)

	// Initialize the viewport
	r.viewport = viewport.New(viewportWidth, r.availableHeight-5)

	// We are locked and loaded
	r.loaded = true
}

// Update the tab
func (r RssFeedTab) Update(msg tea.Msg) (tab.Tab, tea.Cmd) {
	switch msg := msg.(type) {
	case backend.FetchErrorMessage:
		if !r.loaded {
			r.failed = true
			return r, nil
		}

	case backend.FetchSuccessMessage:
		if !r.loaded {
			r.loadTab(msg.Items)
			return r, nil
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// If the tab is not loaded, return
			if !r.loaded {
				return r, nil
			}

			// Set the width of the styled content for word wrapping
			contentWidth := r.availableWidth - r.availableWidth/4 - 4

			// Get the content of the selected item
			r.viewport.SetContent(
				r.list.SelectedItem().(simpleList.Item).StyleContent(contentWidth),
			)

			// Set the view as open if it isn't
			if !r.isViewportOpen {
				r.isViewportOpen = true
			}

			// We don't need to update the list or the viewport
			return r, nil

		case "left", "right":
			// If the viewport isn't open, don't do anything
			if !r.isViewportOpen {
				return r, nil
			}

			// If the viewport is open, switch the selected pane
			if r.selected == articlesPreview {
				r.selected = articlesList
			} else {
				r.selected = articlesPreview
			}

			// We don't need to update the list or the viewport
			return r, nil
		}

	default:
		// If the model is not loaded, update the loading spinner
		if !r.loaded {
			var cmd tea.Cmd
			r.loadingSpinner, cmd = r.loadingSpinner.Update(msg)
			return r, cmd
		}
	}

	// Update the selected item from the pane
	var cmd tea.Cmd
	if r.selected == articlesList {
		// Prevent the list from updating if we are not loaded yet
		if r.loaded {
			r.list, cmd = r.list.Update(msg)
		}
	} else {
		r.viewport, cmd = r.viewport.Update(msg)
	}

	return r, cmd
}

func (r RssFeedTab) View() string {
	if !r.loaded {
		// The style of the message
		messageStyle := lipgloss.NewStyle().
			MarginLeft(3).
			MarginTop(1)

		var loadingMsg string
		if r.failed {
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
				fmt.Sprintf("%s Loading feed %s", r.loadingSpinner.View(), r.title),
			)
		}

		padding := r.availableHeight - 3 - lipgloss.Height(loadingMsg)
		return loadingMsg + strings.Repeat("\n", padding)
	}

	rssList := r.list.View()
	rssViewport := r.viewport.View()

	// If the view is not open show just the rss list
	if !r.isViewportOpen {
		return style.FocusedStyle.Render(rssList)
	}

	// If the viewport is open and the list is selected, highlight the list
	if r.selected == articlesList {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			style.FocusedStyle.Render(rssList),
			style.ColumnStyle.Render(rssViewport),
		)
	}

	// Highlight the viewport
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		style.ColumnStyle.Render(rssList),
		style.FocusedStyle.Render(rssViewport),
	)
}

// Return the type of the tab
func (r RssFeedTab) Type() tab.Type {
	return tab.Feed
}
