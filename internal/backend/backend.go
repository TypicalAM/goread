package backend

import (
	"errors"
	"log"
	"strconv"

	"github.com/TypicalAM/goread/internal/backend/cache"
	"github.com/TypicalAM/goread/internal/backend/rss"
	"github.com/TypicalAM/goread/internal/ui/simplelist"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"
)

// Backend uses a local cache to get all the feeds and their articles
type Backend struct {
	Rss        *rss.Rss
	Cache      *cache.Cache
	ReadStatus *cache.ReadStatus
}

// New creates a new Cache Backend
func New(urlPath, cachePath string, resetCache bool) (*Backend, error) {
	log.Println("Creating new backend")
	store, err := cache.New(cachePath)
	if err != nil {
		return nil, err
	}

	readStatus, err := cache.NewReadStatus("") // TODO: Make this configurable
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if err = rss.Load(); err != nil {
		log.Println("Rss load failed: ", err)
	}

	return &Backend{
		Cache:      store,
		Rss:        rss,
		ReadStatus: readStatus,
	}, nil
}

// FetchCategories returns a tea.Cmd which gets the category list
func (b Backend) FetchCategories(_ string) tea.Cmd {
	return func() tea.Msg {
		names, descs := b.Rss.GetCategories()

		items := make([]list.Item, len(names))
		for i := range names {
			items[i] = simplelist.NewItem(names[i], descs[i])
		}

		return FetchSuccessMsg{Items: items}
	}
}

// FetchFeeds returns a tea.Cmd which gets the feeds for a category
func (b Backend) FetchFeeds(catName string) tea.Cmd {
	return func() tea.Msg {
		names, urls, err := b.Rss.GetFeeds(catName)
		if err != nil {
			return FetchErrorMsg{
				Description: "Error while trying to get feeds",
				Err:         err,
			}
		}

		items := make([]list.Item, len(names))
		for i := range names {
			items[i] = simplelist.NewItem(names[i], urls[i])
		}

		return FetchSuccessMsg{Items: items}
	}
}

// FetchArticles returns a tea.Cmd which gets the articles
func (b Backend) FetchArticles(feedName string) tea.Cmd {
	return func() tea.Msg {
		url, err := b.Rss.GetFeedURL(feedName)
		if err != nil {
			return FetchErrorMsg{
				Description: "Error while trying to get the article url",
				Err:         err,
			}
		}

		items, err := b.Cache.GetArticles(url)
		if err != nil {
			return FetchErrorMsg{
				Description: "Error while fetching the article",
				Err:         err,
			}
		}

		result := make([]list.Item, len(items))
		contents := make([]string, len(items))

		for i, item := range items {
			if b.ReadStatus.IsRead(item) {
				item.Title = "✓ " + item.Title
			}

			result[i] = simplelist.NewItem(item.Title, betterDesc(item.Description))
			contents[i] = rss.YassifyItem(&items[i])
		}

		return FetchArticleSuccessMsg{
			Items:           result,
			ArticleContents: contents,
		}
	}
}

// FetchAllArticles returns a tea.Cmd which gets all the articles
func (b Backend) FetchAllArticles(_ string) tea.Cmd {
	return func() tea.Msg {
		items := b.Cache.GetArticlesBulk(b.Rss.GetAllURLs())

		result := make([]list.Item, len(items))
		contents := make([]string, len(items))

		for i, item := range items {
			if b.ReadStatus.IsRead(item) {
				item.Title = "✓ " + item.Title
			}

			result[i] = simplelist.NewItem(item.Title, betterDesc(item.Description))
			contents[i] = rss.YassifyItem(&items[i])
		}

		return FetchArticleSuccessMsg{
			Items:           result,
			ArticleContents: contents,
		}
	}
}

// FetchDownloaded returns a tea.Cmd which gets the downloaded articles
func (b Backend) FetchDownloadedArticles(_ string) tea.Cmd {
	return func() tea.Msg {
		items := b.Cache.GetDownloaded()

		result := make([]list.Item, len(items))
		contents := make([]string, len(items))

		for i, item := range items {
			if b.ReadStatus.IsRead(item) {
				item.Title = "✓ " + item.Title
			}

			result[i] = simplelist.NewItem(item.Title, betterDesc(item.Description))
			contents[i] = rss.YassifyItem(&items[i])
		}

		return FetchArticleSuccessMsg{
			Items:           result,
			ArticleContents: contents,
		}
	}
}

// DownloadItem returns a tea.Cmd which downloads an item
func (b Backend) DownloadItem(key string, index int) tea.Cmd {
	return func() tea.Msg {
		url, err := b.Rss.GetFeedURL(key)
		if err != nil {
			return FetchErrorMsg{
				Description: "Error while getting the article url",
				Err:         err,
			}
		}

		if err = b.Cache.AddToDownloaded(url, index); err != nil {
			return FetchErrorMsg{
				Description: "Error while downloading the article",
				Err:         err,
			}
		}

		return nil
	}
}

// RemoveDownload tries to remove a download from the backend
func (b Backend) RemoveDownload(key string) error {
	index, err := strconv.Atoi(key)
	if err != nil {
		return errors.New("invalid key")
	}

	return b.Cache.RemoveFromDownloaded(index)
}

// MarkAsRead marks an article as read
func (b Backend) MarkAsRead(key string, index int) tea.Cmd {
	return func() tea.Msg {
		var item gofeed.Item

		switch key {
		case rss.AllFeedsName:
			item = b.Cache.GetArticlesBulk(rss.Default.GetAllURLs())[index]
		case rss.DownloadedFeedsName:
			item = b.Cache.GetDownloaded()[index]
		default:
			url, err := b.Rss.GetFeedURL(key)
			if err != nil {
				return FetchErrorMsg{
					Description: "Error while getting the article url",
					Err:         err,
				}
			}

			items, err := b.Cache.GetArticles(url)
			if err != nil {
				return FetchErrorMsg{
					Description: "Error while fetching the article",
					Err:         err,
				}
			}

			item = items[index]
		}

		log.Println("Marking as read:", item.Title)
		b.ReadStatus.MarkAsRead(item)
		return nil
	}
}

// MarkAsUnread marks an article as unread
func (b Backend) MarkAsUnread(key string, index int) tea.Cmd {
	return func() tea.Msg {
		var item gofeed.Item

		switch key {
		case rss.AllFeedsName:
			item = b.Cache.GetArticlesBulk(rss.Default.GetAllURLs())[index]
		case rss.DownloadedFeedsName:
			item = b.Cache.GetDownloaded()[index]
		default:
			url, err := b.Rss.GetFeedURL(key)
			if err != nil {
				return FetchErrorMsg{
					Description: "Error while getting the article url",
					Err:         err,
				}
			}

			items, err := b.Cache.GetArticles(url)
			if err != nil {
				return FetchErrorMsg{
					Description: "Error while fetching the article",
					Err:         err,
				}
			}

			item = items[index]
		}

		log.Println("Marking as unread:", item.Title)
		b.ReadStatus.MarkAsUnread(item)
		return nil
	}
}

// SetOfflineMode sets the offline mode of the backend
func (b *Backend) SetOfflineMode(mode bool) {
	b.Cache.OfflineMode = mode
}

// Close closes the backend
func (b Backend) Close() error {
	if err := b.Rss.Save(); err != nil {
		return err
	}

	if err := b.Cache.Save(); err != nil {
		return err
	}

	return b.ReadStatus.Save()
}

// betterDesc returns a styled description
func betterDesc(rawDesc string) string {
	desc := rawDesc
	text, err := rss.HTMLToText(rawDesc)
	if err == nil {
		desc = text
	}

	return desc
}
