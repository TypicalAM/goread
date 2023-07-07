package backend

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/TypicalAM/goread/internal/model/simplelist"
	"github.com/TypicalAM/goread/internal/rss"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Backend uses a local cache to get all the feeds and their articles
type Backend struct {
	Rss   *rss.Rss
	Cache *Cache
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
		Cache: cache,
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
				Description: "Error while trying to get feeds",
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
		// Create a list of articles
		url, err := b.Rss.GetFeedURL(feedName)
		if err != nil {
			return FetchErrorMessage{
				Description: "Error while trying to get the article url",
				Err:         err,
			}
		}

		// Get the items from the cache
		items, err := b.Cache.GetArticles(url)
		if err != nil {
			return FetchErrorMessage{
				Description: "Error while fetching the article",
				Err:         err,
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

// FetchAllArticles returns a tea.Cmd which gets all the articles from
// the backend
func (b Backend) FetchAllArticles(_ string) tea.Cmd {
	return func() tea.Msg {
		// Get all the articles and fetch them
		items := b.Cache.GetArticlesBulk(b.Rss.GetAllURLs())

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

// FetchDownloaded returns a tea.Cmd which gets all the downloaded
// articles from the backend
func (b Backend) FetchDownloadedArticles(_ string) tea.Cmd {
	return func() tea.Msg {
		// Get all the downloaded articles
		items := b.Cache.GetDownloaded()

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

// DownloadItem returns a tea.Cmd which downloads an item
func (b Backend) DownloadItem(key string, index int) tea.Cmd {
	return func() tea.Msg {
		// Get the url for the item
		url, err := b.Rss.GetFeedURL(key)
		if err != nil {
			return FetchErrorMessage{
				Description: "Error while getting the article url",
				Err:         err,
			}
		}

		// Get the items from the cache
		err = b.Cache.AddToDownloaded(url, index)
		if err != nil {
			return FetchErrorMessage{
				Description: "Error while downloading the article",
				Err:         err,
			}
		}

		// Return nothing
		return nil
	}
}

// RemoveDownload tries to remove a download from the backend
func (b Backend) RemoveDownload(key string) error {
	index, err := strconv.Atoi(key)
	if err != nil {
		return errors.New("Invalid key")
	}

	return b.Cache.RemoveFromDownloaded(index)
}

// SetOfflineMode sets the offline mode of the backend
func (b *Backend) SetOfflineMode(mode bool) {
	b.Cache.SetOfflineMode(mode)
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
