package backend

import (
	"fmt"

	"github.com/TypicalAM/goread/internal/model/simplelist"
	"github.com/TypicalAM/goread/internal/rss"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"
)

// The Backend uses a local cache to get all the feeds and their articles
type Backend struct {
	Cache *Cache
	Rss   *rss.Rss
}

// New creates a new Cache Backend
func New(urlPath, cachePath string, resetCache bool) (Backend, error) {
	// Create the cache
	cache, err := newStore(cachePath)
	if err != nil {
		return Backend{}, err
	}

	// Try to load the cache
	if !resetCache {
		err = cache.Load()
		if err != nil {
			fmt.Printf("Failed to load the cache: %v, creating a new one", err)
		}
	}

	// Return the backend
	rss := rss.New(urlPath)
	return Backend{
		Cache: &cache,
		Rss:   &rss,
	}, nil
}

// FetchCategories returns a tea.Cmd which gets the category list
// from the backend
func (b Backend) FetchCategories() tea.Cmd {
	return func() tea.Msg {
		// Create a list of categories
		names, descs := b.Rss.GetCategories()

		// Create a list of list items
		items := make([]list.Item, len(names))
		for i := range names {
			items[i] = simplelist.NewItem(names[i], descs[i], "")
		}

		// Return the message
		return FetchSuccessMessage{Items: items}
	}
}

// FetchFeeds returns a tea.Cmd which gets the feed list from
// the backend via a string key
func (b Backend) FetchFeeds(catName string) tea.Cmd {
	return func() tea.Msg {
		// Create a list of feeds
		names, urls, err := b.Rss.GetFeeds(catName)
		if err != nil {
			return FetchErrorMessage{
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
		return FetchSuccessMessage{Items: items}
	}
}

// FetchArticles returns a tea.Cmd which gets the articles from
// the backend via a string key
func (b Backend) FetchArticles(feedName string) tea.Cmd {
	return func() tea.Msg {
		var items []gofeed.Item
		var err error

		if feedName == rss.AllFeedsName {
			// Get all the articles and fetch them
			items, err = b.Cache.GetAllArticles(b.Rss.GetAllURLs())
			if err != nil {
				return FetchErrorMessage{
					Description: "Failed to parse the article",
					Err:         err,
				}
			}
		} else {
			// Create a list of articles
			url, err := b.Rss.GetFeedURL(feedName)
			if err != nil {
				return FetchErrorMessage{
					Description: "Failed to get the article url",
					Err:         err,
				}
			}

			// Get the items from the cache
			items, err = b.Cache.GetArticle(url)
			if err != nil {
				return FetchErrorMessage{
					Description: "Failed to parse the article",
					Err:         err,
				}
			}
		}

		// Create the list of list items
		var result []list.Item
		for i, item := range items {
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
				rss.YassifyItem(&items[i]),
			))
		}

		// Return the message
		return FetchSuccessMessage{Items: result}
	}
}

// Close closes the backend
func (b Backend) Close() error {
	// Try to save the rss
	if err := b.Rss.Save(); err != nil {
		return err
	}

	// Try to save the cache
	return b.Cache.Save()
}
