package cache

import (
	"fmt"

	"github.com/TypicalAM/goread/internal/backend"
	simpleList "github.com/TypicalAM/goread/internal/list"
	"github.com/TypicalAM/goread/internal/rss"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// The Cache Backend uses a local cache to get all the feeds and their articles
type Backend struct {
	Cache Cache
	rss   rss.Rss
}

// New creates a new Cache Backend
func New(urlFilePath string) backend.Backend {
	// TODO: Make the path configurable
	cache := newCache()

	// Save the cache if it doesn't exist (crate the file)
	if err := cache.Load(); err != nil {
		// TODO: Logging
		fmt.Println("Cache doesn't exist ", err)
	}

	// Return the backend
	return Backend{
		Cache: cache,
		rss:   rss.New(urlFilePath),
	}
}

// Name returns the name of the backend
func (b Backend) Name() string {
	return "CacheBackend"
}

// FetchCategories returns a tea.Cmd which gets the category list
// fron the backend
func (b Backend) FetchCategories() tea.Cmd {
	return func() tea.Msg {
		// Create a list of categories
		categories := b.rss.GetCategories()

		// Create a list of list items
		items := make([]list.Item, len(categories))
		for i, cat := range categories {
			items[i] = simpleList.NewListItem(cat, "", "")
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
		feeds, err := b.rss.GetFeeds(catName)
		if err != nil {
			return backend.FetchErrorMessage{
				Description: "Failed to get feeds",
				Err:         err,
			}
		}

		// Create a list of list items
		items := make([]list.Item, len(feeds))
		for i, feed := range feeds {
			items[i] = simpleList.NewListItem(feed, "", "")
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
				Description: "Failed to get the article url",
				Err:         err,
			}
		}

		// Get the items from the cache
		items, err := b.Cache.GetArticle(url)
		if err != nil {
			return backend.FetchErrorMessage{
				Description: "Failed to parse the article",
				Err:         err,
			}
		}

		// Create the list of list items
		var result []list.Item
		for _, item := range items {
			result = append(result, simpleList.NewListItem(
				item.Title,
				rss.HTMLToText(item.Description),
				rss.Markdownize(item),
			))
		}

		// Return the message
		return backend.FetchSuccessMessage{Items: result}
	}
}

// Close closes the backend
func (b Backend) Close() error {
	// Try to save the rss
	if err := b.rss.Save(); err != nil {
		return err
	}

	// Try to save the cache
	return b.Cache.Save()
}
