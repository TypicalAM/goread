package rss

import "errors"

// ErrAlreadyExists is returned when a category/feed already exists
var ErrAlreadyExists = errors.New("already exists")
var ErrTooManyItems = errors.New("too many items")
var ErrReservedName = errors.New("reserved name")

// AddCategory will add a category to the Rss structure
func (rss *Rss) AddCategory(name string, description string) error {
	// Check if there are too many categories
	if len(rss.Categories) >= 36 {
		return ErrTooManyItems
	}

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
	// Check if the name is reserved
	if name == AllFeedsName {
		return ErrReservedName
	}

	// Check if the feed already exists
	for _, cat := range rss.Categories {
		if cat.Name == category {
			// Check if there are too many feeds
			if len(cat.Subscriptions) >= 36 {
				return ErrTooManyItems
			}

			// Check if the feed already exists
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
			return nil
		}
	}

	// We couldn't find the category
	return ErrNotFound
}

// RemoveCategory will remove a category from the Rss structure
func (rss *Rss) RemoveCategory(name string) error {
	for i, cat := range rss.Categories {
		// Check if the category matches
		if cat.Name != name {
			continue
		}

		// Check if the category cannot be removed
		if cat.Name == AllFeedsName {
			return ErrReservedName
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
func (rss *Rss) UpdateCategory(key, name, desc string) error {
	// Check if the name is reserved
	if name == AllFeedsName {
		return ErrReservedName
	}

	// Check if the category already exists
	for _, cat := range rss.Categories {
		if cat.Name == name && name != key {
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

// UpdateFeed will change the name/url of a feed by a string key and a category
func (rss *Rss) UpdateFeed(category, key, name, url string) error {
	// Check if the name is reserved
	if name == AllFeedsName {
		return ErrReservedName
	}

	// Find the category
	for _, cat := range rss.Categories {
		if cat.Name == category {
			// Find the feed
			for _, feed := range cat.Subscriptions {
				if feed.Name == name && name != key {
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
