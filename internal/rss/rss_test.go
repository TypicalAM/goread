package rss

import (
	"reflect"
	"strconv"
	"testing"
)

// TestRssLoadNoFile if we get an error then the default categories are not generated
func TestRssLoadNoFile(t *testing.T) {
	// Create the rss object with an invalid/non-existent file
	myRss := New("../test/data/no-file")

	// Create the default data that we will be checking against
	defaultData := createBasicCategories()

	// If the data is not the same as the default data, then we have a problem
	if !reflect.DeepEqual(myRss.Categories, defaultData) {
		t.Errorf("default categories not created")
	}
}

// TestRssLoadFile if we get an error then the urls file is not loaded correctly
func TestRssLoadFile(t *testing.T) {
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check that we have the correct number of categories
	if len(myRss.Categories) != 2 {
		t.Errorf("incorrect number of categories, expected 2, got %d", len(myRss.Categories))
	}

	// Check if the first category has the correct feed
	if myRss.Categories[0].Subscriptions[0].Name != "Primordial soup" {
		t.Errorf("incorrect name, expected Primordial soup, got %s", myRss.Categories[0].Subscriptions[0].Name)
	}
}

// TestRssGetCategories if we get an error then the rss categories are not retrieved correctly
func TestRssGetCategories(t *testing.T) {
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check if the categories are returned correctly
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
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check if the feeds are returned correctly from the first category
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

	// Check if the feeds are returned correctly from a non-existent category
	_, _, err = myRss.GetFeeds("Non-existent")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrNotFound {
		t.Errorf("incorrect error, expected ErrNotFound, got %s", err)
	}
}

// TestRssGetFeedURL if we get an error then the rss feed url is not retrieved correctly
func TestRssGetFeedURL(t *testing.T) {
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check if the feed url is returned correctly from the first category and the first feed
	url, err := myRss.GetFeedURL(myRss.Categories[0].Subscriptions[0].Name)
	if err != nil {
		t.Errorf("failed to get feed url, %s", err)
	}

	if url != "https://primordialsoup.info/feed" {
		t.Errorf("incorrect url, expected https://primordialsoup.info/feed, got %s", url)
	}

	// Check if the feed url is returned correctly from a non-existent feed
	_, err = myRss.GetFeedURL("Non-existent")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrNotFound {
		t.Errorf("incorrect error, expected ErrNotFound, got %s", err)
	}
}

// TestRssGetAllURLs if we get an error then all the urls are not retrieved correctly
func TestRssGetAllURLs(t *testing.T) {
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check if the urls are returned correctly
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
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check if we can add a category
	err := myRss.AddCategory("New", "New category")
	if err != nil {
		t.Errorf("failed to add category, %s", err)
	}

	// Check if we can get the new category
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

	// Check if we can add a new category with the same name
	err = myRss.AddCategory("New", "Some other new category")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrAlreadyExists {
		t.Errorf("incorrect error, expected ErrAlreadyExists, got %s", err)
	}

	// Check if we can add a new category if there are more than 36 already
	for i := 0; i < 36; i++ {
		_ = myRss.AddCategory(strconv.Itoa(i), "Some other new category")
	}

	err = myRss.AddCategory("37", "Some other new category")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

// TestRssCategoryUpdate if we get an error updating a category doesn't work
func TestRssCategoryUpdate(t *testing.T) {
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check if we can update a category
	err := myRss.UpdateCategory("News", "New", "New category")
	if err != nil {
		t.Errorf("failed to update category, %s", err)
	}

	// Check if we can get the updated category
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

	// Check if we can update a category to have the name AllFeedsName
	err = myRss.UpdateCategory("New", AllFeedsName, "New category")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrReservedName {
		t.Errorf("incorrect error, expected ErrReservedName, got %s", err)
	}

	// Check if we can update a non-existent category
	err = myRss.UpdateCategory("Non-existent", "A little bit of trolling", "Some other new category")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrNotFound {
		t.Errorf("incorrect error, expected ErrNotFound, got %s", err)
	}

	// Check if we can update a category with the same name
	err = myRss.UpdateCategory("New", "New", "Some other new category")
	if err != nil {
		t.Errorf("failed to update category, %s", err)
	}

	// Check if we can get the updated category
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
	err = myRss.UpdateCategory("New", "Technology", "Some other new category")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrAlreadyExists {
		t.Errorf("incorrect error, expected ErrAlreadyExists, got %s", err)
	}
}

// TestRssCategoryRemove if we get an error removing a category doesn't work
func TestRssCategoryRemove(t *testing.T) {
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check if we can remove a category
	err := myRss.RemoveCategory("News")
	if err != nil {
		t.Errorf("failed to remove category, %s", err)
	}

	// Check if we can get the categories
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

	// Check if we can remove a non-existent category
	err = myRss.RemoveCategory("Non-existent")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrNotFound {
		t.Errorf("incorrect error, expected ErrNotFound, got %s", err)
	}
}

// TestRssFeedAdd if we get an error adding a feed doesn't work
func TestRssFeedAdd(t *testing.T) {
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check if we can add a feed
	err := myRss.AddFeed("News", "New feed", "https://new.feed")
	if err != nil {
		t.Errorf("failed to add feed, %s", err)
	}

	// Check if we can get the new feed
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

	// Check if we can add a new feed with the same name
	err = myRss.AddFeed("News", "New feed", "https://new.feed")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrAlreadyExists {
		t.Errorf("incorrect error, expected ErrAlreadyExists, got %s", err)
	}

	// Check if we can add a new feed with the name AllFeedsName
	err = myRss.AddFeed("News", AllFeedsName, "https://new.feed")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrReservedName {
		t.Errorf("incorrect error, expected ErrReservedName, got %s", err)
	}

	// Check if we can add a new to a non-existent category
	err = myRss.AddFeed("Non-existent", "New feed", "https://new.feed")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrNotFound {
		t.Errorf("incorrect error, expected ErrNotFound, got %s", err)
	}

	// Check if we can add a new feed when there are more than 36 already
	for i := 0; i < 36; i++ {
		_ = myRss.AddFeed("News", strconv.Itoa(i), "https://new.feed")
	}

	err = myRss.AddFeed("News", "New feed37", "https://new.feed")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}
}

// TestRssFeedUpdate if we get an error updating a feed doesn't work
func TestRssFeedUpdate(t *testing.T) {
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check if we can update a feed
	err := myRss.UpdateFeed("News", "Primordial soup", "New feed", "https://new.feed")
	if err != nil {
		t.Errorf("failed to update feed, %s", err)
	}

	// Check if we can get the updated feed
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
		t.Errorf("incorrect error, expected ErrReservedName, got %s", err)
	}

	// Check if we can update a non-existent feed
	err = myRss.UpdateFeed("News", "Non-existent", "Other feed", "https://other.feed")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrNotFound {
		t.Errorf("incorrect error, expected ErrNotFound, got %s", err)
	}

	// Check if we can update a feed with the same name but diffferent url
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
	// Create the rss object with a valid file
	myRss := New("../test/data/urls.yml")

	// Check if we can remove a feed
	err := myRss.RemoveFeed("News", "Primordial soup")
	if err != nil {
		t.Errorf("failed to remove feed, %s", err)
	}

	// Check if we can get the feeds
	names, _, err := myRss.GetFeeds("News")
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(names) != 0 {
		t.Errorf("incorrect number of feeds, expected 0, got %d", len(names))
	}

	// Check if we can remove a non-existent feed
	err = myRss.RemoveFeed("News", "Non-existent")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrNotFound {
		t.Errorf("incorrect error, expected ErrNotFound, got %s", err)
	}

	// Check if we can remove a feed from a non-existent category
	err = myRss.RemoveFeed("Non-existent", "Primordial soup")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if err != ErrNotFound {
		t.Errorf("incorrect error, expected ErrNotFound, got %s", err)
	}
}
