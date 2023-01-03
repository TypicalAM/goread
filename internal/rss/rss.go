package rss

import (
	"errors"
	"os"
	"sort"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v3"
)

// Rss will be used to structurize the rss feeds and categories
// it will usually be read from a file
type Rss struct {
	filePath   string     `yaml:"file_path"`
	Categories []Category `yaml:"categories"`
}

// Category will be used to structurize the rss feeds
type Category struct {
	Name          string
	Description   string `yaml:"desc"`
	Subscriptions []Feed `yaml:"subscriptions"`
}

// Feed is a single rss feed
type Feed struct {
	Name        string `yaml:"name"`
	Description string `yaml:"desc"`
	URL         string `yaml:"url"`
}

// ErrNotFound is returned when a feed or category is not found
var ErrNotFound = errors.New("not found")

// New will create a new Rss structure, it fills it with basic items for now
func New(filePath string) Rss {
	rss := Rss{filePath: "rss.yml"}
	err := rss.loadFromFile()
	if err == nil {
		return rss
	}

	rss.Categories = append(rss.Categories, Category{
		Name:        "News",
		Description: "News from around the world",
	})

	rss.Categories = append(rss.Categories, Category{
		Name:        "Tech",
		Description: "Tech news",
	})

	rss.Categories[0].Subscriptions = append(rss.Categories[0].Subscriptions, Feed{
		Name:        "BBC",
		Description: "News from the BBC",
		URL:         "http://feeds.bbci.co.uk/news/rss.xml",
	})

	rss.Categories[1].Subscriptions = append(rss.Categories[1].Subscriptions, Feed{
		Name:        "Hacker News",
		Description: "News from Hacker News",
		URL:         "https://news.ycombinator.com/rss",
	})

	rss.Categories[1].Subscriptions = append(rss.Categories[1].Subscriptions, Feed{
		Name:        "Golang subreddit",
		Description: "News from the Golang subreddit",
		URL:         "https://www.reddit.com/r/golang/.rss",
	})

	return rss
}

// loadFromFile will load the Rss structure from a file
func (rss *Rss) loadFromFile() error {
	// Try to open the file
	fileContent, err := os.ReadFile(rss.filePath)
	if err != nil {
		return err
	}

	// Try to decode the file
	err = yaml.Unmarshal(fileContent, rss)
	if err != nil {
		return err
	}

	// Successfully loaded the file
	return nil
}

// WriteToFile will write the Rss structure to a file
func (rss Rss) WriteToFile() error {
	// Try to marshall the data
	yamlData, err := yaml.Marshal(rss)
	if err != nil {
		return err
	}

	// Try to open the file, if it doesn't exist, create it
	file, err := os.OpenFile(rss.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the data to the file
	_, err = file.Write(yamlData)
	if err != nil {
		return err
	}

	// Successfully wrote the file
	return nil
}

// GetCategories will return a alphabetically sorted list of all categories
func (rss Rss) GetCategories() []string {
	// Create a list of categories
	categories := make([]string, len(rss.Categories))
	for i, cat := range rss.Categories {
		categories[i] = cat.Name
	}

	// Sort the list
	sort.Strings(categories)

	// Return the list
	return categories
}

// GetFeeds will return a alphabetically sorted list of the feeds
// in a category denoted by the name
func (rss Rss) GetFeeds(categoryName string) ([]string, error) {
	// Find the category
	for _, cat := range rss.Categories {
		if cat.Name == categoryName {
			// Create a list of feeds
			feeds := make([]string, len(cat.Subscriptions))
			for i, feed := range cat.Subscriptions {
				feeds[i] = feed.Name
			}

			// Sort the list
			sort.Strings(feeds)

			// Return the list
			return feeds, nil
		}
	}

	// Category not found
	return nil, ErrNotFound
}

// GetFeedURL will return the url of a feed denoted by the name
func (rss Rss) GetFeedURL(feedName string) (string, error) {
	// Iterate over all categories
	for _, cat := range rss.Categories {
		// Iterate over all feeds
		for _, feed := range cat.Subscriptions {
			if feed.Name == feedName {
				return feed.URL, nil
			}
		}
	}

	// Feed not found
	return "", ErrNotFound
}

// Glamourize will return a string that can be used to display the rss feeds
func Glamourize(item gofeed.Item) string {
	var mdown string

	// Add the title
	mdown += "# " + item.Title + "\n "

	// If there are no authors, then we don't want to add a newline
	if item.Authors != nil {
		mdown += item.Authors[0].Name + "\n"
	}

	// Show when the article was published if available
	if item.PublishedParsed != nil {
		mdown += "\n"
		mdown += "Published: " + item.PublishedParsed.Format("2006-01-02 15:04:05")
	}

	mdown += "\n\n"

	// Convert the description html to markdown
	converter := md.NewConverter("", true, nil)
	descMarkdown, err := converter.ConvertString(item.Description)
	if err != nil {
		panic(err)
	}

	mdown += descMarkdown

	return mdown
}
