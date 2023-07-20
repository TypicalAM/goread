package backend

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

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
	return sa[a].PublishedParsed.Before(
		*sa[b].PublishedParsed,
	)
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
	offlineMode bool
}

// Entry is a cache entry
type Entry struct {
	Expire   time.Time        `json:"expire"`
	Articles SortableArticles `json:"articles"`
}

// newStore creates a new cache
func newStore(path string) (*Cache, error) {
	if path == "" {
		defaultPath, err := getDefaultPath()
		if err != nil {
			return nil, err
		}

		path = defaultPath
	}

	return &Cache{
		filePath:   path,
		Content:    make(map[string]Entry),
		Downloaded: make(SortableArticles, 0),
	}, nil
}

// load reads the cache from disk
func (c *Cache) load() error {
	if _, err := os.Stat(c.filePath); err != nil && os.IsNotExist(err) {
		return nil
	}

	file, err := os.ReadFile(c.filePath)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, &c); err != nil {
		return err
	}

	// Iterate over the cache and remove any expired items
	for key, value := range c.Content {
		if value.Expire.Before(time.Now()) {
			delete(c.Content, key)
		}
	}

	return nil
}

// save writes the cache to disk
func (c *Cache) save() error {
	cacheData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// Try to write the data to the file
	if err = os.WriteFile(c.filePath, cacheData, 0600); err != nil {
		if err = os.MkdirAll(filepath.Dir(c.filePath), 0755); err != nil {
			return err
		}

		if err = os.WriteFile(c.filePath, cacheData, 0600); err != nil {
			return err
		}
	}

	return nil
}

// getArticles returns an article list using the cache if possible
func (c *Cache) getArticles(url string) (SortableArticles, error) {
	// Delete entry if expired
	if item, ok := c.Content[url]; ok {
		if item.Expire.After(time.Now()) {
			return item.Articles, nil
		}

		delete(c.Content, url)
	}

	if c.offlineMode {
		return nil, fmt.Errorf("offline mode")
	}

	articles, err := fetchArticles(url)
	if err != nil {
		return nil, err
	}

	// Delete oldest item if cache is full
	if len(c.Content) >= DefaultCacheSize {
		var oldestKey string
		var oldestTime time.Time
		for key, value := range c.Content {
			if oldestTime.IsZero() || value.Expire.Before(oldestTime) {
				oldestKey = key
				oldestTime = value.Expire
			}
		}

		delete(c.Content, oldestKey)
	}

	entry := Entry{
		Expire:   time.Now().Add(DefaultCacheDuration),
		Articles: articles,
	}

	c.Content[url] = entry
	return entry.Articles, nil
}

// getArticlesBulk returns a sorted list of articles from all the given urls, ignoring any errors
func (c *Cache) getArticlesBulk(urls []string) SortableArticles {
	var result SortableArticles

	for _, url := range urls {
		if items, err := c.getArticles(url); err == nil {
			result = append(result, items...)
		}
	}

	sort.Sort(result)
	return result
}

// getDownloaded returns a list of downloaded items
func (c *Cache) getDownloaded() SortableArticles {
	sort.Sort(c.Downloaded)
	return c.Downloaded
}

// addToDownloaded adds an item to the downloaded list
func (c *Cache) addToDownloaded(url string, index int) error {
	articles, err := c.getArticles(url)
	if err != nil {
		return err
	}

	if index < 0 || index >= len(articles) {
		return fmt.Errorf("index out of range")
	}

	c.Downloaded = append(c.Downloaded, articles[index])
	return nil
}

// removeFromDownloaded removes an item from the downloaded list
func (c *Cache) removeFromDownloaded(index int) error {
	if index < 0 || index >= len(c.Downloaded) {
		return fmt.Errorf("index out of range")
	}

	c.Downloaded = append(c.Downloaded[:index], c.Downloaded[index+1:]...)
	return nil
}

// fetchArticles fetches articles from the internet and returns them
func fetchArticles(url string) (SortableArticles, error) {
	feed, err := parseFeed(url)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	req.Header.Set("User-Agent", "goread:v1.3.4 (by /u/TypicalAM)")

	client := http.Client{
		Transport: &http.Transport{
			Proxy:        http.ProxyFromEnvironment,
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, gofeed.HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	feed, err := gofeed.NewParser().Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return feed, nil
}

// getDefaultPath returns the default path to the cache file
func getDefaultPath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "goread", "cache.json"), nil
}
