package rss

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/gilliek/go-opml/opml"
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
	if len(myRss.Categories) != 2 {
		t.Errorf("incorrect number of categories, expected 2, got %d", len(myRss.Categories))
	}

	if myRss.Categories[0].Name != "News" {
		t.Errorf("incorrect category, expected News, got %s", myRss.Categories[0].Name)
	}

	if myRss.Categories[1].Description != "Discover a new use for your spare transistors!" {
		t.Errorf("incorrect description, expected Discover a new use for your spare transistors!, got %s", myRss.Categories[1].Description)
	}
}

// TestRssGetFeeds if we get an error then the rss feeds are not retrieved correctly
func TestRssGetFeeds(t *testing.T) {
	myRss := getRss(t)

	feeds, err := myRss.GetFeeds(myRss.Categories[0].Name)
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(feeds) != 1 {
		t.Errorf("incorrect number of urls, expected 1, got %d", len(feeds))
	}

	if feeds[0].Name != "Primordial soup" {
		t.Errorf("incorrect feed, expected Primordial soup, got %s", feeds[0].Name)
	}

	if _, err = myRss.GetFeeds("Non-existent"); err == nil || err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %s", err)
	}
}

// TestRssGetFeedURL if we get an error then the rss feed url is not retrieved correctly
func TestRssGetFeedURL(t *testing.T) {
	myRss := getRss(t)
	feed, err := myRss.GetFeed(myRss.Categories[0].Subscriptions[0].Name)
	if err != nil {
		t.Errorf("failed to get feed url, %s", err)
	}

	if feed.URL != "https://primordialsoup.info/feed" {
		t.Errorf("incorrect url, expected https://primordialsoup.info/feed, got %s", feed.URL)
	}

	if _, err = myRss.GetFeed("Non-existent"); err == nil || err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %s", err)
	}
}

// TestRssGetAllURLs if we get an error then all the urls are not retrieved correctly
func TestRssGetAllURLs(t *testing.T) {
	myRss := getRss(t)
	urls := myRss.GetAllFeeds()

	if len(urls) != 3 {
		t.Errorf("incorrect number of urls, expected 2, got %d", len(urls))
	}

	if urls[0].URL != "https://primordialsoup.info/feed" {
		t.Errorf("incorrect url, expected https://primordialsoup.info/feed, got %s", urls[0])
	}
}

// TestRssCategoryAdd if we get an error adding a category doesn't work
func TestRssCategoryAdd(t *testing.T) {
	myRss := getRss(t)
	if err := myRss.AddCategory("New", "New category"); err != nil {
		t.Errorf("failed to add category, %s", err)
	}

	if len(myRss.Categories) != 3 {
		t.Errorf("incorrect number of categories, expected 3, got %d", len(myRss.Categories))
	}

	if myRss.Categories[2].Name != "New" {
		t.Errorf("incorrect category, expected New, got %s", myRss.Categories[2].Name)
	}

	if myRss.Categories[2].Description != "New category" {
		t.Errorf("incorrect description, expected New category, got %s", myRss.Categories[2].Description)
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

	if len(myRss.Categories) != 2 {
		t.Errorf("incorrect number of categories, expected 2, got %d", len(myRss.Categories))
	}

	if myRss.Categories[0].Name != "New" {
		t.Errorf("incorrect category, expected New, got %s", myRss.Categories[0].Name)
	}

	if myRss.Categories[0].Description != "New category" {
		t.Errorf("incorrect description, expected New category, got %s", myRss.Categories[0].Description)
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

	if len(myRss.Categories) != 2 {
		t.Errorf("incorrect number of categories, expected 2, got %d", len(myRss.Categories))
	}

	if myRss.Categories[0].Name != "New" {
		t.Errorf("incorrect category, expected New, got %s", myRss.Categories[0].Name)
	}

	if myRss.Categories[0].Description != "Some other new category" {
		t.Errorf("incorrect description, expected Some other new category, got %s", myRss.Categories[0].Description)
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

	if len(myRss.Categories) != 1 {
		t.Errorf("incorrect number of categories, expected 1, got %d", len(myRss.Categories))
	}

	if myRss.Categories[0].Name != "Technology" {
		t.Errorf("incorrect category, expected Technology, got %s", myRss.Categories[0].Name)
	}

	if myRss.Categories[0].Description != "Discover a new use for your spare transistors!" {
		t.Errorf("incorrect description, expected Discover a new use for your spare transistors!, got %s", myRss.Categories[0].Description)
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

	feeds, err := myRss.GetFeeds("News")
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(feeds) != 2 {
		t.Errorf("incorrect number of feeds, expected 2, got %d", len(feeds))
	}

	if feeds[1].Name != "New feed" {
		t.Errorf("incorrect feed, expected New feed, got %s", feeds[1].Name)
	}

	if feeds[1].URL != "https://new.feed" {
		t.Errorf("incorrect url, expected https://new.feed, got %s", feeds[1].URL)
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

	feeds, err := myRss.GetFeeds("News")
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(feeds) != 1 {
		t.Errorf("incorrect number of feeds, expected 1, got %d", len(feeds))
	}

	if feeds[0].Name != "New feed" {
		t.Errorf("incorrect feed, expected New feed, got %s", feeds[0].Name)
	}

	if feeds[0].URL != "https://new.feed" {
		t.Errorf("incorrect url, expected https://new.feed, got %s", feeds[0].URL)
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
	feeds, err = myRss.GetFeeds("News")
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(feeds) != 1 {
		t.Errorf("incorrect number of feeds, expected 1, got %d", len(feeds))
	}

	if feeds[0].URL != "https://new.feed2" {
		t.Errorf("incorrect url, expected https://new.feed2, got %s", feeds[0].URL)
	}
}

// TestRssFeedRemove if we get an error removing a feed doesn't work
func TestRssFeedRemove(t *testing.T) {
	myRss := getRss(t)
	if err := myRss.RemoveFeed("News", "Primordial soup"); err != nil {
		t.Errorf("failed to remove feed, %s", err)
	}

	feeds, err := myRss.GetFeeds("News")
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(feeds) != 0 {
		t.Errorf("incorrect number of feeds, expected 0, got %d", len(feeds))
	}

	if err = myRss.RemoveFeed("News", "Non-existent"); err == nil || err != ErrNotFound {
		t.Errorf("expected ErrNotFound got %s", err)
	}

	if err = myRss.RemoveFeed("Non-existent", "Primordial soup"); err == nil || err != ErrNotFound {
		t.Errorf("expected ErrNotFound got %s", err)
	}
}

// TestOPMLImport if we get an error importing an OPML file doesn't work
func TestRssOPMLImport(t *testing.T) {
	myRss := &Rss{}
	if err := myRss.LoadOPML("../../test/data/opml_flat.xml"); err != nil {
		t.Errorf("failed to import OPML, %s", err)
	}

	feeds, err := myRss.GetFeeds(DefaultCategoryName)
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(feeds) != 3 {
		t.Errorf("incorrect number of feeds, expected 3, got %d", len(feeds))
	}

	myRss = getRss(t)
	if err := myRss.LoadOPML("../../test/data/opml_nested.xml"); err != nil {
		t.Errorf("failed to import OPML, %s", err)
	}

	for _, cat := range myRss.Categories {
		log.Println(cat.Name)
	}

	feeds, err = myRss.GetFeeds("test")
	if err != nil {
		t.Errorf("failed to get feeds, %s", err)
	}

	if len(feeds) != 10 {
		t.Errorf("incorrect number of feeds, expected 3, got %d", len(feeds))
	}
}

// TestOPMLExport if we get an error exporting an OPML file doesn't work
func TestOPMLExport(t *testing.T) {
	rss := getRss(t)
	if err := rss.ExportOPML("test.xml"); err != nil {
		t.Errorf("failed to export OPML, %s", err)
	}

	parsed, err := opml.NewOPMLFromFile("test.xml")
	if err != nil {
		t.Errorf("failed to parse the exported xml into a struct, %s", err)
	}

	if len(parsed.Body.Outlines) != 2 {
		t.Errorf("")
	}

	if len(parsed.Body.Outlines[0].Outlines) != 1 {
		t.Errorf("incorrect number of feeds, expected 1, got %d", len(parsed.Body.Outlines[0].Outlines))
	}

	if err := os.Remove("test.xml"); err != nil {
		t.Errorf("cannot remove the fake file, %s", err)
	}
}
