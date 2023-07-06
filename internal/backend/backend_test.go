package backend

import (
	"testing"
	"time"

	"github.com/TypicalAM/goread/internal/rss"
)

// getCache returns a new cache with the fake data
func getCache() (*Cache, error) {
	cache, err := newStore("../test/data/cache.json")
	if err != nil {
		return nil, err
	}

	err = cache.Load()
	if err != nil {
		return nil, err
	}

	return cache, nil
}

// TestCacheLoadNoFile if we get an error then there's no cache file
func TestCacheLoadNoFile(t *testing.T) {
	// Create a cache with no file
	cache, err := newStore("../test/data/no-file")
	if err != nil {
		t.Fatalf("couldn't get default path: %v", err)
	}

	err = cache.Load()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestCacheLoadCorrectly if we get an error then the cache file is bad
func TestCacheLoadCorrectly(t *testing.T) {
	// Create the cache object with a valid file
	cache, err := getCache()
	if err != nil {
		t.Fatalf("couldn't load the cache %v", err)
	}

	// Check if the cache is loaded correctly
	if len(cache.Content) != 1 {
		t.Fatal("expected 1 item in cache")
	}

	if _, ok := cache.Content["https://primordialsoup.info/feed"]; !ok {
		t.Fatal("expected https://primordialsoup.info/feed in cache")
	}
}

// TestCacheGetArticles if we get an error when there's a cache miss but the cache doesn't change
func TestCacheGetArticles(t *testing.T) {
	// Create the cache object with a valid file
	cache, err := getCache()
	if err != nil {
		t.Fatalf("couldn't load the cache %v", err)
	}

	// Check if the cache hit works
	_, err = cache.GetArticles("https://primordialsoup.info/feed")
	if err != nil {
		t.Fatalf("couldn't get article: %v", err)
	}

	if len(cache.Content) != 1 {
		t.Fatal("expected 1 item in cache")
	}

	// Check if the cache miss retrieves the item and puts it inside the cache
	_, err = cache.GetArticles("https://christitus.com/categories/virtualization/index.xml")
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
	// Create the cache object with a valid file
	cache, err := getCache()
	if err != nil {
		t.Fatalf("couldn't load the cache %v", err)
	}

	// Get the item from the cache
	oldItem, ok := cache.Content["https://primordialsoup.info/feed"]
	if !ok {
		t.Fatal("expected https://primordialsoup.info/feed in cache")
	}

	// Make the item expired and insert it back into the map
	oldItem.Expire = time.Now().Add(-2 * DefaultCacheDuration)
	cache.Content["https://primordialsoup.info/feed"] = oldItem

	_, err = cache.GetArticles("https://primordialsoup.info/feed")
	if err != nil {
		t.Fatalf("couldn't get article: %v", err)
	}

	// Check if item expiry is updated (cache miss)
	newItem, ok := cache.Content["https://primordialsoup.info/feed"]
	if !ok {
		t.Fatal("expected https://primordialsoup.info/feed in cache")
	}

	if newItem.Expire == oldItem.Expire {
		t.Fatal("expected the data to be refreshed and the expire to be updated")
	}
}

// getBackend creates a fake backend
func getBackend() (*Backend, error) {
	b, err := New("../test/data/urls.yml", "", false)
	if err != nil {
		return nil, err
	}

	return &b, err
}

// TestBackendLoad if we get an error loading doesn't work
func TestBackendLoad(t *testing.T) {
	// Create a backend with non-existent file
	_, err := New("../test/data/no-file", "", false)
	if err != nil {
		t.Fatal("expected no error, got", err)
	}

	// Create a backend with a valid file
	b, err := getBackend()
	if err != nil {
		t.Errorf("couldn't get the urls from the file")
	}

	// Check if the backend is loaded correctly
	if len(b.Rss.Categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(b.Rss.Categories))
	}
}

// TestBackendGetCategories if we get an error getting items doesn't work
func TestBackendGetCategories(t *testing.T) {
	// Create a backend with a valid file
	b, err := getBackend()
	if err != nil {
		t.Errorf("couldn't get the urls from the file")
	}

	// Try to fetch the categories
	result := b.FetchCategories()()
	if msg, ok := result.(FetchSuccessMessage); ok {
		if len(msg.Items) != 2 {
			t.Errorf("expected 2 items, got %d", len(msg.Items))
		}
	} else {
		t.Errorf("expected FetchSuccessMessage, got %T", msg)
	}
}

// TestBackendGetFeeds if we get an error getting feeds from a category doesn't work
func TestBackendGetFeeds(t *testing.T) {
	// Create a backend with a valid file
	b, err := getBackend()
	if err != nil {
		t.Errorf("couldn't get the urls from the file")
	}

	// Try to fetch the feeds
	result := b.FetchFeeds("News")()
	switch msg := result.(type) {
	case FetchSuccessMessage:
		if len(msg.Items) != 1 {
			t.Errorf("expected 1 item, got %d", len(msg.Items))
		}

	case FetchErrorMessage:
		t.Errorf("expected FetchSuccessMessage, got a FetchError message %v", msg.Err)

	default:
		t.Errorf("expected FetchSuccessMessage, got %T", msg)
	}

	// Try to fetch the feeds from a non-existent category
	result = b.FetchFeeds("No Category")()
	switch msg := result.(type) {
	case FetchSuccessMessage:
		t.Errorf("expected FetchErrorMessage, got a FetchSuccessMessage with %v items", len(msg.Items))

	case FetchErrorMessage:
		if msg.Err != rss.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", msg.Err)
		}

	default:
		t.Errorf("expected FetchErrorMessage, got %T", msg)
	}
}

// TestBackendGetArticles if we get an error getting items from a feed doesn't work
func TestBackendGetArticles(t *testing.T) {
	// Create a backend with a valid file
	b, err := getBackend()
	if err != nil {
		t.Errorf("couldn't get the urls from the file")
	}

	// Try to fetch the articles for a feed
	result := b.FetchArticles("Primordial soup")()
	switch msg := result.(type) {
	case FetchSuccessMessage:
		if len(msg.Items) != 9 {
			t.Errorf("expected 9 items, got %d", len(msg.Items))
		}

		if msg.Items[0].FilterValue() != "Blindingly Bright Side: the Problem of Positive Psychology in the Dialectic of Conscious and Unconscious" {
			t.Errorf("expected BBC, got %s", msg.Items[0].FilterValue())
		}

	case FetchErrorMessage:
		t.Errorf("expected FetchSuccessMessage, got a FetchErrorMessage with %v", msg.Err)

	default:
		t.Errorf("expected FetchSuccessMessage, got %T", msg)
	}

	// Try to fetch the articles for a non-existent feed
	result = b.FetchArticles("No Feed")()
	switch msg := result.(type) {
	case FetchSuccessMessage:
		t.Errorf("expected FetchErrorMessage, got a FetchSuccessMessage with %v items", len(msg.Items))

	case FetchErrorMessage:
		if msg.Err != rss.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", msg.Err)
		}

	default:
		t.Errorf("expected FetchErrorMessage, got %T", msg)
	}
}
