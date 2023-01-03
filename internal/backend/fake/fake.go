package fake

import (
	"os"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/TypicalAM/goread/internal/backend"
	simpleList "github.com/TypicalAM/goread/internal/list"
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
		file, err := os.Open("rss.rss")
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
		for i := range feed.Items {
			content := CreateFakeContent(i, feed)
			result = append(result, simpleList.NewListItem(
				content.Title,
				strings.Join(content.Categories, ", "),
				content.MoreContent(),
			))
		}

		// Return the message
		return backend.FetchSuccessMessage{Items: result}
	}
}

// Fake backend logic
type Item struct {
	Title           string           `json:"title,omitempty"`
	Description     string           `json:"description,omitempty"`
	Content         string           `json:"content,omitempty"`
	Links           []string         `json:"links,omitempty"`
	UpdatedParsed   *time.Time       `json:"updatedParsed,omitempty"`
	PublishedParsed *time.Time       `json:"publishedParsed,omitempty"`
	Authors         []*gofeed.Person `json:"authors,omitempty"`
	Image           *gofeed.Image    `json:"image,omitempty"`
	Categories      []string         `json:"categories,omitempty"`
}

// A function to create fake text content
func CreateFakeContent(id int, feed *gofeed.Feed) *Item {
	item := &Item{}

	item.Title = feed.Items[id].Title
	item.Description = feed.Items[id].Description
	item.Content = feed.Items[id].Content
	item.Links = feed.Items[id].Links
	item.UpdatedParsed = feed.Items[id].UpdatedParsed
	item.PublishedParsed = feed.Items[id].PublishedParsed
	item.Authors = feed.Items[id].Authors
	item.Image = feed.Items[id].Image
	item.Categories = feed.Items[id].Categories

	return item
}

func (i Item) MoreContent() string {
	var mdown string

	// Add the title
	mdown += "# " + i.Title + "\n "

	// If there are no authors, then we don't want to add a newline
	if i.Authors != nil {
		mdown += i.Authors[0].Name + "\n"
	}

	// Show when the article was published if available
	if i.PublishedParsed != nil {
		mdown += "\n"
		mdown += "Published: " + i.PublishedParsed.Format("2006-01-02 15:04:05")
	}

	mdown += "\n\n"
	mdown += htmlToMarkdown(i.Description)

	return mdown
}

// htmlToMarkdown converts html to markdown using the html-to-markdown library
func htmlToMarkdown(content string) string {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(content)
	if err != nil {
		panic(err)
	}

	return markdown
}

// Close the backend
func (b Backend) Close() error {
	return nil
}
