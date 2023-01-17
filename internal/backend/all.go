package backend

import (
	"sort"

	"github.com/mmcdole/gofeed"
)

type itemList []gofeed.Item

// Len returns the length of the item list, needed for sorting
func (i itemList) Len() int {
	return len(i)
}

// Less returns true if the item at index i is less than the item at index j, needed for sorting
func (i itemList) Less(a, b int) bool {
	return i[a].PublishedParsed.Before(
		*i[b].PublishedParsed,
	)
}

// Swap swaps the items at index i and j, needed for sorting
func (i itemList) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

// GetAllArticles returns an article list from the cache or fetches it from the internet
// if it is not cached and updates the cache, it also updates expired items and sorts
// the items by publish date
func (c *Cache) GetAllArticles(urls []string) ([]gofeed.Item, error) {
	// Create the result slice
	var result []gofeed.Item

	// Iterate over the urls
	for _, url := range urls {
		// Get the article
		items, err := c.GetArticle(url)
		if err != nil {
			return nil, err
		}

		// Add the items to the result
		result = append(result, items...)
	}

	// Sort the items
	sort.Sort(itemList(result))

	// Return the result
	return result, nil
}
