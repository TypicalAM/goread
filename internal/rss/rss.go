package rss

import "errors"

// Rss is the main object used to hold the data
type Rss map[string]Category
type Category map[string]Feed
type Feed string

var errNotFound = errors.New("Resource not found")

func New() Rss {
	rss := Rss{}
	rss["news"] = make(Category)
	rss["news"]["cnn"] = "http://rss.cnn.com/rss/cnn_topstories.rss"
	rss["news"]["reddit"] = "http://www.reddit.com/r/news/.rss"

	rss["technology"] = make(Category)
	rss["technology"]["golang sub"] = "http://www.reddit.com/r/golang/.rss"
	return rss
}

// GetFeedUrl reutrns the url for a feed by a string key
func (rss Rss) GetFeedUrl(name string) (string, error) {
	// Find the category
	for _, cat := range rss {
		// Find the feed
		for k, v := range cat {
			if k == name {
				return string(v), nil
			}
		}
	}

	// The feed was not found
	return "", errNotFound
}

// GetCategory returns a category by name
func (rss Rss) GetCategory(name string) (Category, error) {
	for k, v := range rss {
		if k == name {
			return v, nil
		}
	}

	// The category was not found
	return nil, errNotFound
}

// GetCategories returns a list of categories
func (rss Rss) GetCategories() []string {
	var categories []string
	for k := range rss {
		categories = append(categories, k)
	}
	return categories
}
