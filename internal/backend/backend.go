package backend

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"

	"github.com/TypicalAM/goread/internal/backend/cache"
	"github.com/TypicalAM/goread/internal/backend/rss"
	"github.com/TypicalAM/goread/internal/ui/simplelist"
)

// Backend provides a way of fetching data from the cache and the RSS feed.
type Backend struct {
	Rss        *rss.Rss
	Cache      *cache.Cache
	ReadStatus *cache.ReadStatus
}

// New creates a new backend and its components.
func New(urlPath, cacheDir string, resetCache bool) (*Backend, error) {
	log.Println("Creating new backend")
	store, err := cache.New(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("backend.New: %w", err)
	}

	readStatus, err := cache.NewReadStatus(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("backend.New: %w", err)
	}

	if !resetCache {
		if err = store.Load(); err != nil {
			log.Println("Cache load failed: ", err)
		}

		if err = readStatus.Load(); err != nil {
			log.Println("Read status load failed: ", err)
		}
	}

	rss, err := rss.New(urlPath)
	if err != nil {
		return nil, fmt.Errorf("backend.New: %w", err)
	}

	if err = rss.Load(); err != nil {
		log.Println("Rss load failed: ", err)
		return nil, fmt.Errorf("backend.New: %w", err)
	}

	return &Backend{rss, store, readStatus}, nil
}

// FetchCategories gets the categories.
func (b Backend) FetchCategories(_ string) tea.Cmd {
	return func() tea.Msg {
		items := make([]list.Item, len(b.Rss.Categories))
		for i, cat := range b.Rss.Categories {
			items[i] = simplelist.NewItem(cat.Name, cat.Description)
		}

		return FetchSuccessMsg{Items: items}
	}
}

// FetchFeeds gets the feeds from a category.
func (b Backend) FetchFeeds(catname string) tea.Cmd {
	return func() tea.Msg {
		feeds, err := b.Rss.GetFeeds(catname)
		if err != nil {
			return FetchErrorMsg{err, "Error while trying to get feeds"}
		}

		items := make([]list.Item, len(feeds))
		for i, feed := range feeds {
			items[i] = simplelist.NewItem(feed.Name, feed.URL)
		}

		return FetchSuccessMsg{items}
	}
}

// FetchArticles gets the articles from a feed.
func (b Backend) FetchArticles(feedname string, refresh bool) tea.Cmd {
	return func() tea.Msg {
		feed, err := b.Rss.GetFeed(feedname)
		if err != nil {
			return FetchErrorMsg{err, "Error while trying to get the article url"}
		}

		items, err := b.Cache.GetArticles(feed, refresh)
		if err != nil {
			return FetchErrorMsg{err, "Error while fetching the article"}
		}

		return b.articlesToSuccessMsg(items)
	}
}

// FetchAllArticles gets all the articles from all the feeds.
func (b Backend) FetchAllArticles(_ string, refresh bool) tea.Cmd {
	return func() tea.Msg {
		return b.articlesToSuccessMsg(b.Cache.GetArticlesBulk(b.Rss.GetAllFeeds(), refresh))
	}
}

// FetchDownloaded gets the downloaded articles.
func (b Backend) FetchDownloadedArticles(_ string, _ bool) tea.Cmd {
	return func() tea.Msg {
		return b.articlesToSuccessMsg(b.Cache.GetDownloaded())
	}
}

// DownloadItem downloads an article.
func (b Backend) DownloadItem(feedName string, index int) tea.Cmd {
	return func() tea.Msg {
		item, err := b.indexToItem(feedName, index)
		if err != nil {
			return FetchErrorMsg{err, "Error while getting the article"}
		}

		b.Cache.AddToDownloaded(*item)
		return nil
	}
}

// Close closes the backend and saves its components.
func (b Backend) Close(urlsReadOnly bool) error {
	if !urlsReadOnly {
		if err := b.Rss.Save(); err != nil {
			return fmt.Errorf("backend.Close: %w", err)
		}
	}
	if err := b.Cache.Save(); err != nil {
		return fmt.Errorf("backend.Close: %w", err)
	}

	if err := b.ReadStatus.Save(); err != nil {
		return fmt.Errorf("backend.Close: %w", err)
	}

	return nil
}

// articlesToSuccessMsg converts a list of items to a FetchArticleSuccessMsg.
func (b Backend) articlesToSuccessMsg(items cache.SortableArticles) FetchArticleSuccessMsg {
	sort.Sort(items)
	result := make([]list.Item, len(items))

	savedArticles := b.Cache.GetDownloaded()
	sort.Sort(savedArticles)

	for i, item := range items {
		alreadySaved := false
		for j, saved := range savedArticles {
			if item.Title+item.Link == item.Title+saved.Link {
				alreadySaved = true
				break
			}

			if j > 25 {
				break // NOTE: We don't want to do this for long
			}
		}

		if alreadySaved {
			item.Title = "↓ " + item.Title
		} else if b.ReadStatus.IsRead(item.Link) {
			item.Title = "✓ " + item.Title
		}

		result[i] = ArticleItem{
			ArtTitle:        item.Title,
			RawDesc:         betterDesc(item.Description),
			MarkdownContent: rss.YassifyItem(&items[i]),
			FeedURL:         item.Link,
		}
	}

	return FetchArticleSuccessMsg{result}
}

// indexToItem resolves an index to an item.
func (b Backend) indexToItem(feedName string, index int) (*gofeed.Item, error) {
	var articles cache.SortableArticles

	switch feedName {
	case rss.AllFeedsName:
		articles = b.Cache.GetArticlesBulk(b.Rss.GetAllFeeds(), false)

	case rss.DownloadedFeedsName:
		articles = b.Cache.GetDownloaded()

	default:
		feed, err := b.Rss.GetFeed(feedName)
		if err != nil {
			return nil, errors.New("getting the article url")
		}

		articles, err = b.Cache.GetArticles(feed, false)
		if err != nil {
			return nil, errors.New("fetching the article")
		}
	}

	sort.Sort(articles)
	return &articles[index], nil
}

// betterDesc returns a styled item description.
func betterDesc(rawDesc string) string {
	desc := rawDesc
	text, err := rss.HTMLToText(rawDesc)
	if err == nil {
		desc = text
	}

	return desc
}
