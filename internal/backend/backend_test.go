package backend

import (
	"testing"

	"github.com/TypicalAM/goread/internal/backend/rss"
)

// getBackend creates a fake backend
func getBackend() (*Backend, error) {
	b, err := New("../test/data/urls.yml", "", false)
	if err != nil {
		return nil, err
	}

	return b, err
}

// TestBackendLoad if we get an error loading doesn't work
func TestBackendLoad(t *testing.T) {
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
	result := b.FetchCategories("")()
	if msg, ok := result.(FetchSuccessMsg); ok {
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
	case FetchSuccessMsg:
		if len(msg.Items) != 1 {
			t.Errorf("expected 1 item, got %d", len(msg.Items))
		}

	case FetchErrorMsg:
		t.Errorf("expected FetchSuccessMessage, got a FetchError message %v", msg.Err)

	default:
		t.Errorf("expected FetchSuccessMessage, got %T", msg)
	}

	// Try to fetch the feeds from a non-existent category
	result = b.FetchFeeds("No Category")()
	switch msg := result.(type) {
	case FetchSuccessMsg:
		t.Errorf("expected FetchErrorMessage, got a FetchSuccessMessage with %v items", len(msg.Items))

	case FetchErrorMsg:
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
	result := b.FetchArticles("Primordial soup", false)()
	switch msg := result.(type) {
	case FetchArticleSuccessMsg:
		if len(msg.Items) != 9 {
			t.Errorf("expected 9 items, got %d", len(msg.Items))
		}

		if msg.Items[0].FilterValue() != "Blindingly Bright Side: the Problem of Positive Psychology in the Dialectic of Conscious and Unconscious" {
			t.Errorf("expected BBC, got %s", msg.Items[0].FilterValue())
		}

	case FetchErrorMsg:
		t.Errorf("expected FetchSuccessMessage, got a FetchErrorMessage with %v", msg.Err)

	default:
		t.Errorf("expected FetchSuccessMessage, got %T", msg)
	}

	// Try to fetch the articles for a non-existent feed
	result = b.FetchArticles("No Feed", false)()
	switch msg := result.(type) {
	case FetchSuccessMsg:
		t.Errorf("expected FetchErrorMessage, got a FetchSuccessMessage with %v items", len(msg.Items))

	case FetchErrorMsg:
		if msg.Err != rss.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", msg.Err)
		}

	default:
		t.Errorf("expected FetchErrorMessage, got %T", msg)
	}
}
