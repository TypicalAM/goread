package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/TypicalAM/goread/internal/backend"
	"github.com/TypicalAM/goread/internal/rss"
)

// getCache returns a new cache with the fake data
func getCache() (*Cache, error) {
	cache, err := newStore()
	if err != nil {
		return nil, err
	}

	cache.filePath = "../../test/data/cache.json"
	err = cache.Load()
	if err != nil {
		return nil, err
	}

	return &cache, nil
}

// TestCacheLoadNoFile if we get an error when there's no cache file
func TestCacheLoadNoFile(t *testing.T) {
	cache, err := newStore()
	if err != nil {
		t.Fatalf("couldn't get default path: %v", err)
	}

	cache.filePath = "../test/data/no-file"
	err = cache.Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestCacheLoadCorrectly if we get an error then the cache file is bad
func TestCacheLoadCorrectly(t *testing.T) {
	cache, err := getCache()
	if err != nil {
		t.Fatalf("couldn't load the cache %v", err)
	}

	if len(cache.Content) != 1 {
		t.Fatal("expected 1 item in cache")
	}

	if _, ok := cache.Content["https://primordialsoup.info/feed"]; !ok {
		t.Fatal("expected https://primordialsoup.info/feed in cache")
	}
}

// TestCacheGetArticle if we get an error when there's a cache miss but the cache doesn't change
func TestCacheGetArticle(t *testing.T) {
	cache, err := getCache()
	if err != nil {
		t.Fatalf("couldn't load the cache %v", err)
	}

	_, err = cache.GetArticle("https://primordialsoup.info/feed")
	if err != nil {
		t.Fatalf("couldn't get article: %v", err)
	}

	if len(cache.Content) != 1 {
		t.Fatal("expected 1 item in cache")
	}

	_, err = cache.GetArticle("https://christitus.com/categories/virtualization/index.xml")
	if err != nil {
		t.Fatalf("couldn't get article: %v", err)
	}

	if len(cache.Content) != 2 {
		t.Fatal("expected 2 items in cache")
	}

	if _, ok := cache.Content["https://christitus.com/categories/virtualization/index.xml"]; !ok {
		t.Fatal("expected https://christitus.com/categories/virtualization/index.xml in cache")
	}
}

// TestCacheGetArticleExpired if we get an error then the store doesn't delete expired cache when getting data
func TestCacheGetArticleExpired(t *testing.T) {
	cache, err := getCache()
	if err != nil {
		t.Fatalf("couldn't load the cache %v", err)
	}

	oldItem, ok := cache.Content["https://primordialsoup.info/feed"]
	if !ok {
		t.Fatal("expected https://primordialsoup.info/feed in cache")
	}

	oldItem.Expire = time.Now().Add(-2 * defaultCacheDuration)
	cache.Content["https://primordialsoup.info/feed"] = oldItem

	_, err = cache.GetArticle("https://primordialsoup.info/feed")
	if err != nil {
		t.Fatalf("couldn't get article: %v", err)
	}

	newItem, ok := cache.Content["https://primordialsoup.info/feed"]
	if !ok {
		t.Fatal("expected https://primordialsoup.info/feed in cache")
	}

	fmt.Println(oldItem.Expire.Nanosecond(), newItem.Expire.Nanosecond())
	if newItem.Expire == oldItem.Expire {
		t.Fatal("expected the data to be refreshed and the expire to be updated")
	}
}

// getBackend creates a fake backend
func getBackend() (*Backend, error) {
	b, err := New("../../test/data/urls.yml")
	if err != nil {
		return nil, err
	}

	return &b, err
}

// TestBackendLoad if we get an error loading doesn't work
func TestBackendLoad(t *testing.T) {
	_, err := New("../test/data/no-file")
	if err != nil {
		t.Fatal("expected no error, got", err)
	}

	b, err := getBackend()
	if err != nil {
		t.Errorf("couldn't get the urls from the file")
	}

	if len(b.rss.Categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(b.rss.Categories))
	}
}

// TestBackendGetCategories if we get an error getting items doesn't work
func TestBackendGetCategories(t *testing.T) {
	b, err := getBackend()
	if err != nil {
		t.Errorf("couldn't get the urls from the file")
	}

	result := b.FetchCategories()()
	if msg, ok := result.(backend.FetchSuccessMessage); ok {
		if len(msg.Items) != 2 {
			t.Errorf("expected 2 items, got %d", len(msg.Items))
		}
	} else {
		t.Errorf("expected FetchSuccessMessage, got %T", msg)
	}
}

// TestBackendGetFeeds if we get an error getting feeds from a category doesn't work
func TestBackendGetFeeds(t *testing.T) {
	b, err := getBackend()
	if err != nil {
		t.Errorf("couldn't get the urls from the file")
	}

	result := b.FetchFeeds("News")()
	switch msg := result.(type) {
	case backend.FetchSuccessMessage:
		if len(msg.Items) != 1 {
			t.Errorf("expected 1 item, got %d", len(msg.Items))
		}

	case backend.FetchErrorMessage:
		t.Errorf("expected FetchSuccessMessage, got a FetchError message %v", msg.Err)

	default:
		t.Errorf("expected FetchSuccessMessage, got %T", msg)
	}

	result = b.FetchFeeds("No Category")()
	switch msg := result.(type) {
	case backend.FetchSuccessMessage:
		t.Errorf("expected FetchErrorMessage, got a FetchSuccessMessage with %v items", len(msg.Items))

	case backend.FetchErrorMessage:
		if msg.Err != rss.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", msg.Err)
		}

	default:
		t.Errorf("expected FetchErrorMessage, got %T", msg)
	}
}

// TestBackendGetArticles if we get an error getting items from a feed doesn't work
func TestBackendGetArticles(t *testing.T) {
	b, err := getBackend()
	if err != nil {
		t.Errorf("couldn't get the urls from the file")
	}

	result := b.FetchArticles("Primordial soup")()
	switch msg := result.(type) {
	case backend.FetchSuccessMessage:
		if len(msg.Items) != 9 {
			t.Errorf("expected 9 items, got %d", len(msg.Items))
		}

		if msg.Items[0].FilterValue() != "Blindingly Bright Side: the Problem of Positive Psychology in the Dialectic of Conscious and Unconscious" {
			t.Errorf("expected BBC, got %s", msg.Items[0].FilterValue())
		}

	case backend.FetchErrorMessage:
		t.Errorf("expected FetchSuccessMessage, got a FetchErrorMessage with %v", msg.Err)

	default:
		t.Errorf("expected FetchSuccessMessage, got %T", msg)
	}

	result = b.FetchArticles("No Feed")()
	switch msg := result.(type) {
	case backend.FetchSuccessMessage:
		t.Errorf("expected FetchErrorMessage, got a FetchSuccessMessage with %v items", len(msg.Items))

	case backend.FetchErrorMessage:
		if msg.Err != rss.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", msg.Err)
		}

	default:
		t.Errorf("expected FetchErrorMessage, got %T", msg)
	}
}
