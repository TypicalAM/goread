package rss

import "errors"

var ErrAlreadyExists = errors.New("category already exists")

// AddCategory will add a category to the Rss structure
func (rss *Rss) AddCategory(name string, description string) error {
	// Check if the category already exists
	for _, cat := range rss.Categories {
		if cat.Name == name {
			return ErrAlreadyExists
		}
	}

	// Add the category
	rss.Categories = append(rss.Categories, Category{
		Name:        name,
		Description: description,
	})

	// Return no errors
	return nil
}

// AddFeed will add a feed to the Rss structure
func (rss *Rss) AddFeed(category string, name string, url string) error {
	// Check if the feed already exists
	for _, cat := range rss.Categories {
		if cat.Name == category {
			for _, feed := range cat.Subscriptions {
				if feed.Name == name {
					return ErrAlreadyExists
				}
			}
		}
	}

	// Add the feed
	for i, cat := range rss.Categories {
		if cat.Name == category {
			rss.Categories[i].Subscriptions = append(rss.Categories[i].Subscriptions, Feed{
				Name: name,
				URL:  url,
			})
		}
	}

	// Return no errors
	return nil
}

// RemoveCategory will remove a category from the Rss structure
func (rss *Rss) RemoveCategory(name string) error {
	// Remove the category
	for i, cat := range rss.Categories {
		if cat.Name == name {
			rss.Categories = append(rss.Categories[:i], rss.Categories[i+1:]...)
		}
	}

	// Return no errors
	return nil
}
