package backend

import (
	"github.com/mmcdole/gofeed"
)

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
