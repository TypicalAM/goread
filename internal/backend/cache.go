package backend

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/mmcdole/gofeed"
)

// DefaultCacheDuration is the default duration for which an item is cached
var DefaultCacheDuration = 24 * time.Hour

// DefaultCacheSize is the default size of the cache
var DefaultCacheSize = 100

// Cache is a basic cache to read and write gofeed.Items based on the URL
type Cache struct {
	filePath string
	Content  map[string]Item
}

// Item is an item in the cache
type Item struct {
	Expire time.Time
	Items  []gofeed.Item
}

// newStore creates a new cache
func newStore(path string) (Cache, error) {
	// Get the path to the cache file
	if path == "" {
		defaultPath, err := getDefaultPath()
		if err != nil {
			return Cache{}, err
		}

		path = defaultPath
	}

	// Create the cache
	return Cache{
		filePath: path,
		Content:  make(map[string]Item),
	}, nil
}

// Load reads the cache from disk
func (c *Cache) Load() error {
	// Load the cache from the file
	file, err := os.ReadFile(c.filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &c.Content)
	if err != nil {
		return err
	}

	// Iterate over the cache and remove any expired items
	for key, value := range c.Content {
		if value.Expire.Before(time.Now()) {
			delete(c.Content, key)
		}
	}

	// Return no errors
	return nil
}

// Save writes the cache to disk
func (c *Cache) Save() error {
	// Try to encode the cache
	cacheData, err := json.Marshal(c.Content)
	if err != nil {
		return err
	}

	// Try to write the data to the file
	if err = os.WriteFile(c.filePath, cacheData, 0600); err != nil {
		// Try to create the directory
		err = os.MkdirAll(filepath.Dir(c.filePath), 0755)
		if err != nil {
			return err
		}

		// Try to write to the file again
		err = os.WriteFile(c.filePath, cacheData, 0600)
		if err != nil {
			return err
		}
	}

	// Writing was successful
	return nil
}

// GetArticle returns an article list from the cache or fetches it from the internet
// if it is not cached and updates the cache, it also updates expired items
func (c *Cache) GetArticle(url string) ([]gofeed.Item, error) {
	// Check if the cache contains the url
	if item, ok := c.Content[url]; ok {
		if item.Expire.After(time.Now()) {
			return item.Items, nil
		}

		// Fetch the cacheItem from the internet
		cacheItem, err := fetchArticle(url)
		if err != nil {
			return nil, err
		}

		// Add the item to the cache
		c.Content[url] = cacheItem
		return cacheItem.Items, nil
	}

	// Fetch the cacheItem from the internet
	cacheItem, err := fetchArticle(url)
	if err != nil {
		return nil, err
	}

	// Check if the cache is full
	if len(c.Content) >= DefaultCacheSize {
		// Find the oldest item
		var oldestKey string
		var oldestTime time.Time
		for key, value := range c.Content {
			if oldestTime.IsZero() || value.Expire.Before(oldestTime) {
				oldestKey = key
				oldestTime = value.Expire
			}
		}

		// Remove the item
		delete(c.Content, oldestKey)
	}

	// Add the item to the cache
	c.Content[url] = cacheItem
	return cacheItem.Items, nil
}

// fetchArticle fetches an article list from the internet and returns a slice of items
func fetchArticle(url string) (Item, error) {
	// Parse the url
	feed, err := parseURL(url)
	if err != nil {
		return Item{}, err
	}

	// Parse the items
	items := make([]gofeed.Item, len(feed.Items))
	for i, item := range feed.Items {
		items[i] = *item
	}

	// Return the items
	return Item{
		Expire: time.Now().Add(DefaultCacheDuration),
		Items:  items,
	}, nil
}

// parseURL parses a url and attempts to return a parsed feed
// authors note: this is made because the gofeed parser does not support some feeds, namely the ones from reddit
func parseURL(feedURL string) (*gofeed.Feed, error) {
	// Create a new client
	var client = http.Client{
		Transport: &http.Transport{
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		},
	}

	// Create a new request with our user agent
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "goread:v1.0.0 (by /u/TypicalAM)")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check the response
	if resp != nil {
		defer func() {
			ce := resp.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
	}

	// Check the status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, gofeed.HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	// Try to parse the body
	return gofeed.NewParser().Parse(resp.Body)
}

// getDefaultPath returns the default path to the cache file
func getDefaultPath() (string, error) {
	// Get the temporary directory
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	// Join the path
	return filepath.Join(dir, "goread", "cache.json"), nil
}
