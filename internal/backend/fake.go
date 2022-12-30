package backend

import (
	"os"
	"time"

	simpleList "github.com/TypicalAM/goread/internal/list"
	"github.com/TypicalAM/goread/internal/style"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
)

// Create a fake backed for testing
type FakeBackend struct{}

func (b FakeBackend) Name() string {
	return "FakeBackend"
}

// Return some fake categories
func (b FakeBackend) FetchCategories() tea.Cmd {
	return func() tea.Msg {
		return FetchSuccessMessage{
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
func (b FakeBackend) FetchFeeds(catName string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		return FetchSuccessMessage{
			Items: []list.Item{
				simpleList.NewListItem("feed 1", "cat1", "more content"),
				simpleList.NewListItem("feed 2", "cat2", "more content"),
				simpleList.NewListItem("feed 3", "cat3", "more content"),
			},
		}
	}
}

// Return some fake articles
func (b FakeBackend) FetchArticles(feedName string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		file, err := os.Open("rss.rss")
		if err != nil {
			return FetchErrorMessage{
				description: "Could not open file",
				Err:         err,
			}
		}

		defer file.Close()
		fp := gofeed.NewParser()
		feed, err := fp.Parse(file)
		if err != nil {
			return FetchErrorMessage{
				description: "Could not parse file",
				Err:         err,
			}
		}

		var result []list.Item
		for i := range feed.Items {
			content := CreateFakeContent(i, feed)
			result = append(result, simpleList.NewListItem(
				content.Title,
				content.Description,
				content.MoreContent(),
			))
		}

		// Return the message
		return FetchSuccessMessage{Items: result}
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
	item.Description = parseHTML(feed.Items[id].Description)
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
	var sections []string

	titleTextStyle := lipgloss.NewStyle().
		Foreground(style.BasicColorscheme.Color1).
		Bold(true)

	descTextStyle := lipgloss.NewStyle().
		Foreground(style.BasicColorscheme.Color2).
		Bold(true).
		Width(70)

	fixedWidth := lipgloss.NewStyle().
		Width(70)

	sections = append(
		sections,
		titleTextStyle.Render(i.Title), "",
		descTextStyle.Render(parseHTML(i.Description)),
		fixedWidth.Render(parseHTML(i.Content)),
	)

	if len(i.Links) > 0 {
		sections = append(sections, titleTextStyle.Render("Links"), "")
		for _, link := range i.Links {
			sections = append(sections, fixedWidth.Render(link))
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		sections...,
	)
}
