package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/mmcdole/gofeed"
)

// Cache is a basic cache to read and write gofeed.Items based on the URL
type Cache struct {
	path    string
	Content map[string]Item
}

// Item is an item in the cache
type Item struct {
	Expire time.Time
	Items  []gofeed.Item
}

// newCache creates a new cache
func newCache() Cache {
	// Get the path to the cache file
	path, err := getCachePath()
	// TODO: Handle error
	if err != nil {
		panic(err)
	}

	return Cache{
		path:    path,
		Content: make(map[string]Item),
	}
}

// Load reads the cache from disk
func (c *Cache) Load() error {
	// Load the cache from the file
	file, err := os.ReadFile(c.path)
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
	f, err := os.Create(c.path)
	if err != nil {
		// Try to create the directory
		err = os.MkdirAll(filepath.Dir(c.path), 0755)
		if err != nil {
			return err
		}

		// Try to create the file again
		f, err = os.Create(c.path)
		if err != nil {
			return err
		}
	}
	defer f.Close()

	// Try to encode the cache
	content, err := json.Marshal(c.Content)
	if err != nil {
		return err
	}

	// Try to write the cache to disk
	_, err = f.Write(content)
	if err != nil {
		return err
	}

	// Writing was successful
	return nil
}

// GetArticle returns an article list from the cache or
// fetches it from the internet if it is not in the cache
func (c *Cache) GetArticle(url string) ([]gofeed.Item, error) {
	// Check if the cache contains the url
	if item, ok := c.Content[url]; ok {
		// Check if the item is expired
		if item.Expire.After(time.Now()) {
			// Return the items
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

	// Add the item to the cache
	c.Content[url] = cacheItem
	return cacheItem.Items, nil
}

// fetchArticle fetches an article list from the internet and
// reutrns a slice of gofeed.Items
func fetchArticle(url string) (Item, error) {
	// Create a new feed parser
	fp := gofeed.NewParser()

	// Parse the feed
	feed, err := fp.ParseURL(url)
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
		Expire: time.Now().Add(24 * time.Hour),
		Items:  items,
	}, nil
}

// getCachedPath returns the path to the cache file
func getCachePath() (string, error) {
	// Get the temporary directory
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	// Join the path
	return filepath.Join(dir, "goread", "cache.json"), nil
}
