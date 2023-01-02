package rss

import "errors"

// Rss is the main object used to hold the data
type Rss map[string]Category
type Category map[string]Feed
type Feed string

var errNotFound = errors.New("Resource not found")

func New() Rss {
	rss := Rss{}
	// TODO: passing in and reading from a config file
	rss["news"] = make(Category)
	rss["news"]["cnn"] = "http://rss.cnn.com/rss/cnn_topstories.rss"

	rss["technology"] = make(Category)
	rss["technology"]["golang sub"] = "http://www.reddit.com/r/golang/.rss"
	rss["technology"]["chris titus android"] = "https://christitus.com/categories/android/index.xml"
	rss["technology"]["chris titus linux"] = "https://christitus.com/categories/linux/index.xml"
	return rss
}

// GetFeedUrl reutrns the url for a feed by a string key
func (rss Rss) GetFeedURL(name string) (string, error) {
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
	categories := make([]string, len(rss))
	count := 0

	for k := range rss {
		categories[count] = k
		count++
	}

	return categories
}
