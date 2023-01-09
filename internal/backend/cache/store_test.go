package cache

import (
	"fmt"
	"testing"
	"time"
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

// TestLoadNoFile if we get an error when there's no cache file
func TestLoadNoFile(t *testing.T) {
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

// TestLoadCorrectly if we get an error then the cache file is bad
func TestLoadCorrectly(t *testing.T) {
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

// TestGetArticle if we get an error when there's a cache miss but the cache doesn't change
func TestGetArticle(t *testing.T) {
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

// TestGetArticleExpired if we get an error then the store doesn't delete expired cache when getting data
func TestGetArticleExpired(t *testing.T) {
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
