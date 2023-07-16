package rss

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v3"
)

// AllFeedsName is the name of the all feeds category
var AllFeedsName = "All Feeds"

// DownloadedFeedsName is the name of the downloaded feeds category
var DownloadedFeedsName = "Saved"

// ErrNotFound is returned when a feed or category is not found
var ErrNotFound = errors.New("not found")

// Rss will be used to structurize the rss feeds and categories
// it will usually be read from a file
type Rss struct {
	FilePath   string     `yaml:"-"`
	Categories []Category `yaml:"categories"`
}

// Category will be used to structurize the rss feeds
type Category struct {
	Name          string `yaml:"name"`
	Description   string `yaml:"desc"`
	Subscriptions []Feed `yaml:"subscriptions"`
}

// Feed is a single rss feed
type Feed struct {
	Name        string `yaml:"name"`
	Description string `yaml:"desc"`
	URL         string `yaml:"url"`
}

// New will create a new Rss structure
func New(urlFilePath string) Rss {
	// Create the rss object
	rss := Rss{FilePath: urlFilePath}

	// Check if we can load it from file
	err := rss.loadFromFile()
	if err == nil {
		return rss
	}

	// Append some default categories
	rss.Categories = append(rss.Categories, createBasicCategories()...)
	return rss
}

// loadFromFile will load the Rss structure from a file
func (rss *Rss) loadFromFile() error {
	// Check if the path is valid
	if rss.FilePath == "" {
		// Get the default path
		path, err := getDefaultPath()
		if err != nil {
			return err
		}

		// Set the path
		rss.FilePath = path
	}

	// Try to open the file
	fileContent, err := os.ReadFile(rss.FilePath)
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

// Save will write the Rss structure to a file
func (rss *Rss) Save() error {
	// Try to marshall the data
	yamlData, err := yaml.Marshal(rss)
	if err != nil {
		return err
	}

	// Try to write the data to the file
	if err = os.WriteFile(rss.FilePath, yamlData, 0600); err != nil {
		// Try to create the directory
		err = os.MkdirAll(filepath.Dir(rss.FilePath), 0755)
		if err != nil {
			return err
		}

		// Try to write to the file again
		err = os.WriteFile(rss.FilePath, yamlData, 0600)
		if err != nil {
			return err
		}
	}

	// Successfully wrote the file
	return nil
}

// GetCategories will return a list of all the names and descriptions of the categories
func (rss Rss) GetCategories() (names []string, descs []string) {
	// Create a list of categories
	names = make([]string, len(rss.Categories))
	descs = make([]string, len(rss.Categories))

	for i, cat := range rss.Categories {
		names[i] = cat.Name
		descs[i] = cat.Description
	}

	// Return the list
	return names, descs
}

// GetFeeds will return a list of all the names and descriptions of the feeds
// in a category denoted by the name
func (rss Rss) GetFeeds(categoryName string) (names []string, urls []string, err error) {
	// Find the category
	for _, cat := range rss.Categories {
		if cat.Name == categoryName {
			// Create a list of feeds
			feeds := make([]string, len(cat.Subscriptions))
			urls = make([]string, len(cat.Subscriptions))
			for i, feed := range cat.Subscriptions {
				feeds[i] = feed.Name
				urls[i] = feed.URL
			}

			// Return the list
			return feeds, urls, nil
		}
	}

	// Category not found
	return nil, nil, ErrNotFound
}

// GetFeedURL will return the url of a feed denoted by the name
func (rss Rss) GetFeedURL(feedName string) (string, error) {
	// Check if the feed is reserved
	if feedName == AllFeedsName || feedName == DownloadedFeedsName {
		return "", ErrReservedName
	}

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

// GetAllURLs will return a list of all the urls
func (rss Rss) GetAllURLs() []string {
	// Create a list of urls
	var urls []string

	// Iterate over all categories
	for _, cat := range rss.Categories {
		// Iterate over all feeds
		for _, feed := range cat.Subscriptions {
			if feed.URL != AllFeedsName {
				urls = append(urls, feed.URL)
			}
		}
	}

	// Return the list
	return urls
}

// getDefaultPath will return the default path for the urls file
func getDefaultPath() (string, error) {
	// Get the default config path
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// Create the config path
	return filepath.Join(configDir, "goread", "urls.yml"), nil
}

// YassifyItem will return a yassified string which is used in the viewport
// to view a single item
func YassifyItem(item *gofeed.Item) string {
	var mdown string

	// Add the title
	mdown += "# " + item.Title + "\n "

	// If there are no authors, then don't add the author
	if item.Authors != nil {
		mdown += item.Authors[0].Name + "\n"
	}

	// Show when the article was published if available
	if item.PublishedParsed != nil {
		mdown += "\n"
		mdown += "Published: " + item.PublishedParsed.Format("2006-01-02 15:04:05")
	}

	// Convert the html to markdown
	mdown += "\n\n"
	htmlMarkdown, err := HTMLToMarkdown(item.Description)
	if err != nil {
		// If there is an error, then just print the html
		mdown += item.Description
	} else {
		mdown += htmlMarkdown
	}

	// Add the links if there are any
	if len(item.Links) > 0 {
		mdown += "\n## Links\n"
		for _, link := range item.Links {
			mdown += "- " + link + "\n"
		}
	}

	// Add padding
	mdown += "\n\n"

	// Return the markdown
	return mdown
}

// HTMLToMarkdown converts html to markdown using the html-to-markdown library
func HTMLToMarkdown(content string) (string, error) {
	// Create a new converter
	converter := md.NewConverter("", true, nil)

	// Convert the html to markdown
	markdown, err := converter.ConvertString(content)
	if err != nil {
		return "", err
	}

	// Return the markdown
	return markdown, nil
}

// HTMLToText converts html to text using the goquery library
func HTMLToText(content string) (string, error) {
	// Create a new document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return "", err
	}

	// Return the text
	return doc.Text(), nil
}

// createBasicCategories will create some basic categories
func createBasicCategories() []Category {
	var categories []Category

	categories = append(categories, Category{
		Name:        AllFeedsName,
		Description: "All feeds",
	})

	categories = append(categories, Category{
		Name:        "News",
		Description: "News from around the world",
	})

	categories = append(categories, Category{
		Name:        "Tech",
		Description: "Tech news",
	})

	categories[0].Subscriptions = append(categories[0].Subscriptions, Feed{
		Name:        "BBC",
		Description: "News from the BBC",
		URL:         "http://feeds.bbci.co.uk/news/rss.xml",
	})

	categories[1].Subscriptions = append(categories[1].Subscriptions, Feed{
		Name:        "Wired",
		Description: "News from the wired team",
		URL:         "https://www.wired.com/feed/rss",
	})

	categories[1].Subscriptions = append(categories[1].Subscriptions, Feed{
		Name:        "Chris Titus Tech (virtualization)",
		Description: "Chris Titus Tech on virtualization",
		URL:         "https://christitus.com/categories/virtualization/index.xml",
	})

	return categories
}
