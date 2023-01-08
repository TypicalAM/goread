package web

import (
	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/model/simplelist"
	"github.com/TypicalAM/goread/internal/rss"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"
)

// The Web Backend uses the internet to get all the feeds and their articles
type Backend struct {
	rss *rss.Rss
}

// New returns a new WebBackend
func New(urlFilePath string) Backend {
	rss := rss.New(urlFilePath)
	return Backend{rss: &rss}
}

// Name returns the name of the backend
func (b Backend) Name() string {
	return "WebBackend"
}

// FetchCategories returns a tea.Cmd which gets the category list
// fron the backend
func (b Backend) FetchCategories() tea.Cmd {
	return func() tea.Msg {
		// Create a list of categories
		names, descs := b.rss.GetCategories()

		// Create a list of list items
		items := make([]list.Item, len(names))
		for i := range names {
			items[i] = simplelist.NewItem(names[i], descs[i], "")
		}

		// Return the message
		return backend.FetchSuccessMessage{Items: items}
	}
}

// FetchFeeds returns a tea.Cmd which gets the feed list from
// the backend via a string key
func (b Backend) FetchFeeds(catName string) tea.Cmd {
	return func() tea.Msg {
		// Create a list of feeds
		names, urls, err := b.rss.GetFeeds(catName)
		if err != nil {
			return backend.FetchErrorMessage{
				Description: "Failed to get feeds",
				Err:         err,
			}
		}

		// Create a list of list items
		items := make([]list.Item, len(names))
		for i := range names {
			items[i] = simplelist.NewItem(names[i], urls[i], "")
		}

		// Return the message
		return backend.FetchSuccessMessage{Items: items}
	}
}

// FetchArticles returns a tea.Cmd which gets the articles from
// the backend via a string key
func (b Backend) FetchArticles(feedName string) tea.Cmd {
	return func() tea.Msg {
		// Create a list of articles
		url, err := b.rss.GetFeedURL(feedName)
		if err != nil {
			return backend.FetchErrorMessage{
				Description: "Failed to get articles",
				Err:         err,
			}
		}

		// Get the articles and parse them using gofeed
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(url)
		if err != nil {
			return backend.FetchErrorMessage{
				Description: "Failed to parse the articles",
				Err:         err,
			}
		}

		// Create the list of list items
		var result []list.Item
		for _, item := range feed.Items {
			// Check if the description can be converted to a string
			var description string
			text, err := rss.HTMLToText(item.Description)
			if err != nil {
				description = item.Description
			} else {
				description = text
			}

			// Create the list item
			result = append(result, simplelist.NewItem(
				item.Title,
				description,
				rss.YassifyItem(item),
			))
		}

		// Return the message
		return backend.FetchSuccessMessage{Items: result}
	}
}

// Rss returns the rss object
func (b Backend) Rss() *rss.Rss {
	return b.rss
}

// Close closes the backend
func (b Backend) Close() error {
	// Try to save the rss
	return b.rss.Save()
}
