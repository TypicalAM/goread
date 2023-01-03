package fake

import (
	"os"
	"path/filepath"

	"github.com/TypicalAM/goread/internal/backend"
	simpleList "github.com/TypicalAM/goread/internal/list"
	"github.com/TypicalAM/goread/internal/rss"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"
)

// Create a fake backed for testing
type Backend struct{}

// Create a new fake backend
func New() Backend {
	return Backend{}
}

// Name returns the name of the backend
func (b Backend) Name() string {
	return "FakeBackend"
}

// Return some fake categories
func (b Backend) FetchCategories() tea.Cmd {
	return func() tea.Msg {
		return backend.FetchSuccessMessage{
			Items: []list.Item{
				simpleList.NewListItem("All", "All the categories", ""),
				simpleList.NewListItem("Books", "Books", "books"),
				simpleList.NewListItem("Movies", "Movies", "movies"),
				simpleList.NewListItem("Music", "Music", "music"),
				simpleList.NewListItem("Games", "Games", "games"),
				simpleList.NewListItem("Technology", "Software", "software"),
			},
		}
	}
}

// Return some fake feeds
func (b Backend) FetchFeeds(catName string) tea.Cmd {
	return func() tea.Msg {
		return backend.FetchSuccessMessage{
			Items: []list.Item{
				simpleList.NewListItem("feed 1", "cat1", "more content"),
				simpleList.NewListItem("feed 2", "cat2", "more content"),
				simpleList.NewListItem("feed 3", "cat3", "more content"),
			},
		}
	}
}

// Return some fake articles
func (b Backend) FetchArticles(feedName string) tea.Cmd {
	return func() tea.Msg {
		path := filepath.Join("test", "feeds", "test.xml")
		file, err := os.Open(path)
		if err != nil {
			return backend.FetchErrorMessage{
				Description: "Could not open file",
				Err:         err,
			}
		}

		defer file.Close()
		fp := gofeed.NewParser()
		feed, err := fp.Parse(file)
		if err != nil {
			return backend.FetchErrorMessage{
				Description: "Could not parse file",
				Err:         err,
			}
		}

		var result []list.Item
		for _, item := range feed.Items {
			result = append(result, simpleList.NewListItem(
				item.Title,
				rss.HTMLToText(item.Description),
				rss.Markdownize(*item),
			))
		}

		// Return the message
		return backend.FetchSuccessMessage{Items: result}
	}
}

// AddItem adds an item to the rss
func (b Backend) AddItem(itemType backend.ItemType, fields ...string) {}

// DeleteItem deletes an item from the rss
func (b Backend) DeleteItem(itemType backend.ItemType, key string) {}

// Close the backend
func (b Backend) Close() error {
	return nil
}
