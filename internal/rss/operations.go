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
	for i, cat := range rss.Categories {
		// Check if the category matches
		if cat.Name != name {
			continue
		}

		// Remove the category
		rss.Categories = append(rss.Categories[:i], rss.Categories[i+1:]...)
		return nil
	}

	// We couldn't remove the category
	return ErrNotFound
}

// RemoveFeed will remove a feed from the Rss structure
func (rss *Rss) RemoveFeed(category string, name string) error {
	for i, cat := range rss.Categories {
		// Check if the category matches
		if cat.Name != category {
			continue
		}

		for j, feed := range cat.Subscriptions {
			// Check if the feed matches
			if feed.Name != name {
				continue
			}

			// Remove the feed
			rss.Categories[i].Subscriptions = append(rss.Categories[i].Subscriptions[:j], rss.Categories[i].Subscriptions[j+1:]...)
			return nil
		}
	}

	// We couldn't remove the feed
	return ErrNotFound
}

// UpdateCategory will change the name/description of a category by a string key
func (rss *Rss) UpdateCategory(name, desc string, key string) error {
	// Check if the category already exists
	for _, cat := range rss.Categories {
		if cat.Name == name && cat.Description == desc {
			// FIXME: Name clash
			return ErrAlreadyExists
		}
	}

	// Find the category
	for i, cat := range rss.Categories {
		if cat.Name == key {
			// Update the category
			rss.Categories[i].Name = name
			rss.Categories[i].Description = desc
			return nil
		}
	}

	// We couldn't find the category
	return ErrNotFound
}

// UdpateFeed will change the name/url of a feed by a string key and a category
func (rss *Rss) UpdateFeed(name, url, category, key string) error {
	// Find the category
	for _, cat := range rss.Categories {
		if cat.Name == category {
			// Find the feed
			for _, feed := range cat.Subscriptions {
				if feed.Name == name {
					return ErrAlreadyExists
				}
			}
		}
	}

	// Find the category
	for i, cat := range rss.Categories {
		if cat.Name == category {
			// Find the feed
			for j, feed := range cat.Subscriptions {
				if feed.Name == key {
					// Update the feed
					rss.Categories[i].Subscriptions[j].Name = name
					rss.Categories[i].Subscriptions[j].URL = url
					return nil
				}
			}
		}
	}

	// We couldn't find the feed
	return ErrNotFound
}
