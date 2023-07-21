package rss

import (
	"strconv"
	"testing"
)

func getRss(t *testing.T) *Rss {
	myRss, err := New("../../test/data/urls.yml")
	if err != nil {
		t.Errorf("error creating rss object: %v", err)
	}

	if err = myRss.Load(); err != nil {
		t.Errorf("error loading file: %v", err)
	}

	return myRss
}

// TestRssLoadNoFile if we get an error then the urls file is not loaded correctly
func TestRssLoadNoFile(t *testing.T) {
	myRss, err := New("non-existent")
	if err != nil {
		t.Errorf("error creating rss object: %v", err)
	}

	if err = myRss.Load(); err != nil {
		t.Errorf("no error returned when loading non-existent file")
	}
}

// TestRssLoadFile if we get an error then the urls file is not loaded correctly
func TestRssLoadFile(t *testing.T) {
	myRss := getRss(t)
	if len(myRss.Categories) != 2 {
		t.Errorf("incorrect number of categories, expected 2, got %d", len(myRss.Categories))
	}

	if myRss.Categories[0].Subscriptions[0].Name != "Primordial soup" {
		t.Errorf("incorrect name, expected Primordial soup, got %s", myRss.Categories[0].Subscriptions[0].Name)
	}
}

// TestRssGetCategories if we get an error then the rss categories are not retrieved correctly
func TestRssGetCategories(t *testing.T) {
	myRss := getRss(t)
	names, descs := myRss.GetCategories()
	if len(names) != len(descs) {
		t.Errorf("incorrect number of descriptions, expected %d, got %d", len(names), len(descs))
	}

	if len(names) != 2 {
		t.Errorf("incorrect number of categories, expected 2, got %d", len(names))
	}

	if names[0] != "News" {
		t.Errorf("incorrect category, expected News, got %s", names[0])
	}

	if descs[1] != "Discover a new use for your spare transistors!" {
		t.Errorf("incorrect description, expected Discover a new use for your spare transistors!, got %s", descs[1])
	}
}

// TestRssGetFeeds if we get an error then the rss feeds are not retrieved correctly
func TestRssGetFeeds(t *testing.T) {
	myRss := getRss(t)

	names, urls, err := myRss.GetFeeds(myRss.Categories[0].Name)
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(names) != len(urls) {
		t.Errorf("incorrect number of urls, expected %d, got %d", len(names), len(urls))
	}

	if len(names) != 1 {
		t.Errorf("incorrect number of feeds, expected 1, got %d", len(names))
	}

	if names[0] != "Primordial soup" {
		t.Errorf("incorrect feed, expected Primordial soup, got %s", names[0])
	}

	if _, _, err = myRss.GetFeeds("Non-existent"); err == nil || err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %s", err)
	}
}

// TestRssGetFeedURL if we get an error then the rss feed url is not retrieved correctly
func TestRssGetFeedURL(t *testing.T) {
	myRss := getRss(t)
	url, err := myRss.GetFeedURL(myRss.Categories[0].Subscriptions[0].Name)
	if err != nil {
		t.Errorf("failed to get feed url, %s", err)
	}

	if url != "https://primordialsoup.info/feed" {
		t.Errorf("incorrect url, expected https://primordialsoup.info/feed, got %s", url)
	}

	if _, err = myRss.GetFeedURL("Non-existent"); err == nil || err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %s", err)
	}
}

// TestRssGetAllURLs if we get an error then all the urls are not retrieved correctly
func TestRssGetAllURLs(t *testing.T) {
	myRss := getRss(t)
	urls := myRss.GetAllURLs()

	if len(urls) != 3 {
		t.Errorf("incorrect number of urls, expected 2, got %d", len(urls))
	}

	if urls[0] != "https://primordialsoup.info/feed" {
		t.Errorf("incorrect url, expected https://primordialsoup.info/feed, got %s", urls[0])
	}
}

// TestRssCategoryAdd if we get an error adding a category doesn't work
func TestRssCategoryAdd(t *testing.T) {
	myRss := getRss(t)
	if err := myRss.AddCategory("New", "New category"); err != nil {
		t.Errorf("failed to add category, %s", err)
	}

	names, descs := myRss.GetCategories()
	if len(names) != 3 {
		t.Errorf("incorrect number of categories, expected 3, got %d", len(names))
	}

	if names[2] != "New" {
		t.Errorf("incorrect category, expected New, got %s", names[2])
	}

	if descs[2] != "New category" {
		t.Errorf("incorrect description, expected New category, got %s", descs[2])
	}

	if err := myRss.AddCategory("New", "Some other new category"); err == nil || err != ErrAlreadyExists {
		t.Errorf("expected an error (ErrAlreadyExists), got nil")
	}

	// Check if we can add a new category if there are more than 36 already
	for i := 0; i < 36; i++ {
		_ = myRss.AddCategory(strconv.Itoa(i), "Some other new category")
	}

	if err := myRss.AddCategory("37", "Some other new category"); err == nil {
		t.Errorf("expected an error, got nil")
	}
}

// TestRssCategoryUpdate if we get an error updating a category doesn't work
func TestRssCategoryUpdate(t *testing.T) {
	myRss := getRss(t)
	if err := myRss.UpdateCategory("News", "New", "New category"); err != nil {
		t.Errorf("failed to update category, %s", err)
	}

	names, descs := myRss.GetCategories()
	if len(names) != 2 {
		t.Errorf("incorrect number of categories, expected 2, got %d", len(names))
	}

	if names[0] != "New" {
		t.Errorf("incorrect category, expected New, got %s", names[0])
	}

	if descs[0] != "New category" {
		t.Errorf("incorrect description, expected New category, got %s", descs[0])
	}

	if err := myRss.UpdateCategory("New", AllFeedsName, "New category"); err == nil || err != ErrReservedName {
		t.Errorf("expected an error (ErrReservedName)")
	}

	err := myRss.UpdateCategory("Non-existent", "A little bit of trolling", "Some other new category")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %s", err)
	}

	if err = myRss.UpdateCategory("New", "New", "Some other new category"); err != nil {
		t.Errorf("failed to update category, %s", err)
	}

	names, descs = myRss.GetCategories()
	if len(names) != 2 {
		t.Errorf("incorrect number of categories, expected 2, got %d", len(names))
	}

	if names[0] != "New" {
		t.Errorf("incorrect category, expected New, got %s", names[0])
	}

	if descs[0] != "Some other new category" {
		t.Errorf("incorrect description, expected Some other new category, got %s", descs[0])
	}

	// Check if we can update a category with the same name as an existing category
	if err = myRss.UpdateCategory("New", "Technology", "Some other new category"); err == nil || err != ErrAlreadyExists {
		t.Errorf("expected an error (ErrAlreadyExists)")
	}
}

// TestRssCategoryRemove if we get an error removing a category doesn't work
func TestRssCategoryRemove(t *testing.T) {
	myRss := getRss(t)
	if err := myRss.RemoveCategory("News"); err != nil {
		t.Errorf("failed to remove category, %s", err)
	}

	names, descs := myRss.GetCategories()
	if len(names) != 1 {
		t.Errorf("incorrect number of categories, expected 1, got %d", len(names))
	}

	if names[0] != "Technology" {
		t.Errorf("incorrect category, expected Technology, got %s", names[0])
	}

	if descs[0] != "Discover a new use for your spare transistors!" {
		t.Errorf("incorrect description, expected Discover a new use for your spare transistors!, got %s", descs[0])
	}

	if err := myRss.RemoveCategory("Non-existent"); err == nil || err != ErrNotFound {
		t.Errorf("expected an error (ErrNotFound)")
	}
}

// TestRssFeedAdd if we get an error adding a feed doesn't work
func TestRssFeedAdd(t *testing.T) {
	myRss := getRss(t)
	if err := myRss.AddFeed("News", "New feed", "https://new.feed"); err != nil {
		t.Errorf("failed to add feed, %s", err)
	}

	names, urls, err := myRss.GetFeeds("News")
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(names) != 2 {
		t.Errorf("incorrect number of feeds, expected 2, got %d", len(names))
	}

	if names[1] != "New feed" {
		t.Errorf("incorrect feed, expected New feed, got %s", names[1])
	}

	if urls[1] != "https://new.feed" {
		t.Errorf("incorrect url, expected https://new.feed, got %s", urls[1])
	}

	if err = myRss.AddFeed("News", "New feed", "https://new.feed"); err == nil || err != ErrAlreadyExists {
		t.Errorf("expected an error (ErrAlreadyExists)")
	}

	if err = myRss.AddFeed("News", AllFeedsName, "https://new.feed"); err == nil || err != ErrReservedName {
		t.Errorf("expected an error, got nil")
	}

	if err = myRss.AddFeed("Non-existent", "New feed", "https://new.feed"); err == nil || err != ErrNotFound {
		t.Errorf("expected an error (ErrNotFound)")
	}

	// Check if we can add a new feed when there are more than 36 already
	for i := 0; i < 36; i++ {
		_ = myRss.AddFeed("News", strconv.Itoa(i), "https://new.feed")
	}

	if err = myRss.AddFeed("News", "New feed37", "https://new.feed"); err == nil {
		t.Errorf("expected an error, got nil")
	}
}

// TestRssFeedUpdate if we get an error updating a feed doesn't work
func TestRssFeedUpdate(t *testing.T) {
	myRss := getRss(t)
	if err := myRss.UpdateFeed("News", "Primordial soup", "New feed", "https://new.feed"); err != nil {
		t.Errorf("failed to update feed, %s", err)
	}

	names, urls, err := myRss.GetFeeds("News")
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(names) != 1 {
		t.Errorf("incorrect number of feeds, expected 1, got %d", len(names))
	}

	if names[0] != "New feed" {
		t.Errorf("incorrect feed, expected New feed, got %s", names[0])
	}

	if urls[0] != "https://new.feed" {
		t.Errorf("incorrect url, expected https://new.feed, got %s", urls[0])
	}

	// Check if we can update a feed to have the name AllFeedsName
	err = myRss.UpdateFeed("News", "New feed", AllFeedsName, "https://new.feed")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrReservedName {
		t.Errorf("expected ErrReservedName, got %s", err)
	}

	// Check if we can update a non-existent feed
	err = myRss.UpdateFeed("News", "Non-existent", "Other feed", "https://other.feed")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %s", err)
	}

	// Check if we can update a feed with the same name but different url
	err = myRss.UpdateFeed("News", "New feed", "New feed", "https://new.feed2")
	if err != nil {
		t.Errorf("failed to update feed, %s", err)
	}

	// Check if we can get the updated feed
	names, urls, err = myRss.GetFeeds("News")
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(names) != 1 {
		t.Errorf("incorrect number of feeds, expected 1, got %d", len(names))
	}

	if urls[0] != "https://new.feed2" {
		t.Errorf("incorrect url, expected https://new.feed2, got %s", urls[0])
	}
}

// TestRssFeedRemove if we get an error removing a feed doesn't work
func TestRssFeedRemove(t *testing.T) {
	myRss := getRss(t)
	if err := myRss.RemoveFeed("News", "Primordial soup"); err != nil {
		t.Errorf("failed to remove feed, %s", err)
	}

	names, _, err := myRss.GetFeeds("News")
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(names) != 0 {
		t.Errorf("incorrect number of feeds, expected 0, got %d", len(names))
	}

	if err = myRss.RemoveFeed("News", "Non-existent"); err == nil || err != ErrNotFound {
		t.Errorf("expected ErrNotFound got %s", err)
	}

	if err = myRss.RemoveFeed("Non-existent", "Primordial soup"); err == nil || err != ErrNotFound {
		t.Errorf("expected ErrNotFound got %s", err)
	}
}
