package web

import (
	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/backend/fake"
	simpleList "github.com/TypicalAM/goread/internal/list"
	"github.com/TypicalAM/goread/internal/rss"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"
)

// Backend uses the internet to get all the feeds and their articles
type Backend struct {
	rss rss.Rss
}

// New returns a new WebBackend
func New() Backend {
	return Backend{rss: rss.New()}
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
		for i := range feed.Items {
			content := fake.CreateFakeContent(i, feed)
			result = append(result, simpleList.NewListItem(
				content.Title,
				content.Description,
				content.MoreContent(),
			))
		}

		// Return the message
		return backend.FetchSuccessMessage{Items: result}
	}
}
