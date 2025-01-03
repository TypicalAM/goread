package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/TypicalAM/goread/internal/backend/rss"
	"github.com/mmcdole/gofeed"
)

// DefaultCacheDuration is the default duration for which an item is cached
var DefaultCacheDuration = 24 * time.Hour

// DefaultCacheSize is the default size of the cache
var DefaultCacheSize = 100

// SortableArticles is a sortable list of articles
type SortableArticles []gofeed.Item

// Len returns the length of the item list, needed for sorting
func (sa SortableArticles) Len() int {
	return len(sa)
}

// Less returns true if the item at index i is less than the item at index j, needed for sorting
func (sa SortableArticles) Less(a, b int) bool {
	if sa[a].PublishedParsed == nil || sa[b].PublishedParsed == nil {
		return strings.Map(unicode.ToLower, sa[a].Title) < strings.Map(unicode.ToLower, sa[b].Title)
	}

	return sa[a].PublishedParsed.After(*sa[b].PublishedParsed)
}

// Swap swaps the items at index i and j, needed for sorting
func (sa SortableArticles) Swap(a, b int) {
	sa[a], sa[b] = sa[b], sa[a]
}

// Cache handles the caching of feeds and storing downloaded articles
type Cache struct {
	Content     map[string]Entry `json:"content"`
	filePath    string
	Downloaded  SortableArticles `json:"downloaded"`
	OfflineMode bool             `json:"-"`
}

// Entry is a cache entry
type Entry struct {
	Expire   time.Time        `json:"expire"`
	Articles SortableArticles `json:"articles"`
}

// New creates a new cache store.
func New(dir string) (*Cache, error) {
	log.Println("Creating new cache store")
	if dir == "" {
		defaultDir, err := getDefaultDir()
		if err != nil {
			return nil, fmt.Errorf("cache.New: %w", err)
		}

		dir = defaultDir
	}

	return &Cache{
		filePath:   filepath.Join(dir, "cache.json"),
		Content:    make(map[string]Entry),
		Downloaded: make(SortableArticles, 0),
	}, nil
}

// Load reads the cache from disk
func (c *Cache) Load() error {
	log.Println("Loading cache from", c.filePath)
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("cache.Load: %w", err)
	}

	if err = json.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("cache.Load: %w", err)
	}

	log.Println("Loaded cache entries: ", len(c.Content))
	return nil
}

// Save writes the cache to disk
func (c *Cache) Save() error {
	// Iterate over the cache and remove any expired items
	for key, value := range c.Content {
		if value.Expire.Before(time.Now()) {
			delete(c.Content, key)
		}
	}

	cacheData, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("cache.Save: %w", err)
	}

	// Try to write the data to the file
	if err = os.WriteFile(c.filePath, cacheData, 0600); err != nil {
		if err = os.MkdirAll(filepath.Dir(c.filePath), 0755); err != nil {
			return fmt.Errorf("cache.Save: %w", err)
		}

		if err = os.WriteFile(c.filePath, cacheData, 0600); err != nil {
			return fmt.Errorf("cache.Save: %w", err)
		}
	}

	return nil
}

// GetArticles returns an article list using the cache if possible
func (c *Cache) GetArticles(feed *rss.Feed, ignoreCache bool) (SortableArticles, error) {
	log.Println("Getting articles for", feed.URL, " from cache: ", !ignoreCache)

	// Delete entry if expired
	if item, ok := c.Content[feed.URL]; ok && !ignoreCache {
		if item.Expire.After(time.Now()) {
			return item.Articles, nil
		}

		delete(c.Content, feed.URL)
	}

	if c.OfflineMode {
		return nil, errors.New("offline mode")
	}

	articles, err := fetchArticles(feed.URL)
	if err != nil {
		return nil, fmt.Errorf("cache.GetArticles: %w", err)
	}

	if len(feed.BlacklistWords) != 0 {
		log.Println("Using keyword blacklist for feed", feed.Name, ":", feed.BlacklistWords)
		remaining := make([]gofeed.Item, 0)
		for i := range articles {
			if !includesKeywords(&articles[i], feed.BlacklistWords) {
				remaining = append(remaining, articles[i])
			}
		}

		articles = remaining
	}

	if len(feed.WhitelistWords) != 0 {
		log.Println("Using keyword whitelist for feed", feed.Name, ":", feed.WhitelistWords)
		remaining := make([]gofeed.Item, 0)
		for i := range articles {
			if includesKeywords(&articles[i], feed.WhitelistWords) {
				remaining = append(remaining, articles[i])
			}
		}

		articles = remaining
	}

	c.Content[feed.URL] = Entry{time.Now().Add(DefaultCacheDuration), articles}
	return articles, nil
}

// GetArticlesBulk returns a sorted list of articles from all the given urls, ignoring any errors
func (c *Cache) GetArticlesBulk(feeds []*rss.Feed, ignoreCache bool) SortableArticles {
	var result SortableArticles

	for i, feed := range feeds {
		if items, err := c.GetArticles(feeds[i], ignoreCache); err == nil {
			result = append(result, items...)
		} else {
			// NOTE: Let's say you have 50 feeds and 5 fail, we don't want to keep trying failed feeds
			// so we just fill the cache with an empty item. That way load for bulk feeds is faster next time.
			log.Println("Error getting articles for", feed.URL, err, "filling with empty item")
			c.Content[feed.URL] = Entry{time.Now().Add(DefaultCacheDuration), SortableArticles{}}
		}
	}

	return result
}

// GetDownloaded returns a list of downloaded items
func (c *Cache) GetDownloaded() SortableArticles {
	return c.Downloaded
}

// AddToDownloaded adds an item to the downloaded list
func (c *Cache) AddToDownloaded(item gofeed.Item) {
	c.Downloaded = append(c.Downloaded, item)
}

// RemoveFromDownloaded removes an item from the downloaded list
func (c *Cache) RemoveFromDownloaded(index int) error {
	if index < 0 || index >= len(c.Downloaded) {
		return errors.New("index out of range")
	}

	c.Downloaded = append(c.Downloaded[:index], c.Downloaded[index+1:]...)
	return nil
}

// fetchArticles fetches articles from the internet and returns them
func fetchArticles(url string) (SortableArticles, error) {
	log.Println("Fetching articles from", url)
	feed, err := parseFeed(url)
	if err != nil {
		return nil, fmt.Errorf("cache.fetchArticles: %w", err)
	}

	items := make(SortableArticles, len(feed.Items))
	for i, item := range feed.Items {
		items[i] = *item
	}

	return items, nil
}

// parseFeed parses a url and attempts to return a parsed feed
// authors note: this is was because the gofeed parser did not support reddit
func parseFeed(url string) (*gofeed.Feed, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("cache.parseFeed: %w", err)
	}
	req.Header.Set("User-Agent", "goread (by /u/TypicalAM)")

	client := http.Client{
		Transport: &http.Transport{
			Proxy:        http.ProxyFromEnvironment,
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("cache.parseFeed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, gofeed.HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	feed, err := gofeed.NewParser().Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cache.parseFeed: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("cache.parseFeed: %w", err)
	}

	if feed.PublishedParsed == nil && strings.TrimSpace(feed.Published) != "" {
		// Fix for #64 since I don't want to bother gofeed's maintainer with date format additions
		if t, err := time.Parse("Fri, 03 Jan 2025 17:20:38 GMT", strings.TrimSpace(feed.Published)); err == nil {
			feed.PublishedParsed = &t
		}
	}

	return feed, nil
}

// getDefaultDir returns the default cache directory
func getDefaultDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("cache.getDefaultDir: %w", err)
	}

	return filepath.Join(dir, "goread"), nil
}

// includesKeywords checks if an article contains any specified keyword from a slice
func includesKeywords(feed *gofeed.Item, keywords []string) bool {
	for _, keyword := range keywords {
		lowerKeyword := strings.ToLower(keyword)
		if strings.Contains(strings.ToLower(feed.Title), lowerKeyword) ||
			strings.Contains(strings.ToLower(feed.Description), lowerKeyword) ||
			strings.Contains(strings.ToLower(feed.Content), lowerKeyword) {
			return true
		}
	}

	return false
}
