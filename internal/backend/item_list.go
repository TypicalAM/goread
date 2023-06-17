package backend

import (
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
