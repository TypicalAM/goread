package rss

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gilliek/go-opml/opml"
	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v3"
)

// AllFeedsName is the name of the all feeds category
var AllFeedsName = "All Feeds"

// DownloadedFeedsName is the name of the downloaded feeds category
var DownloadedFeedsName = "Saved"

// DefaultCategoryName is the name of the default category
var DefaultCategoryName = "News"

// DefaultCategoryDescription is the description of the default category
var DefaultCategoryDescription = "News from around the world"

// ErrNotFound is returned when a feed or category is not found
var ErrNotFound = errors.New("not found")

// Default is the default rss structure
var Default = Rss{
	Categories: []Category{{
		Name:        AllFeedsName,
		Description: "All feeds",
		Subscriptions: []Feed{{
			Name:        "BBC",
			Description: "News from the BBC",
			URL:         "http://feeds.bbci.co.uk/news/rss.xml",
		}},
	}, {
		Name:        "News",
		Description: "News from around the world",
		Subscriptions: []Feed{{
			Name:        "Wired",
			Description: "News from the wired team",
			URL:         "https://www.wired.com/feed/rss",
		}},
	}, {
		Name:        "Tech",
		Description: "Tech news",
		Subscriptions: []Feed{{
			Name:        "Chris Titus Tech (virtualization)",
			Description: "Chris Titus Tech on virtualization",
			URL:         "https://christitus.com/categories/virtualization/index.xml",
		}},
	}},
}

// Rss will be used to structurize the rss feeds and categories
type Rss struct {
	filePath   string
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
func New(path string) (*Rss, error) {
	log.Println("Creating new rss structure")
	if path == "" {
		defaultPath, err := getDefaultPath()
		if err != nil {
			return nil, err
		}

		// Set the path
		path = defaultPath
	}

	rss := Default
	rss.filePath = path
	return &rss, nil
}

// Load will try to load the Rss structure from a file
func (rss *Rss) Load() error {
	log.Println("Loading rss from", rss.filePath)
	if _, err := os.Stat(rss.filePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	data, err := os.ReadFile(rss.filePath)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, rss); err != nil {
		return err
	}

	log.Printf("Rss loaded with %d categories\n", len(rss.Categories))
	return nil
}

// Save will write the Rss structure to a file
func (rss Rss) Save() error {
	yamlData, err := yaml.Marshal(rss)
	if err != nil {
		return err
	}

	if err = os.WriteFile(rss.filePath, yamlData, 0600); err != nil {
		if err = os.MkdirAll(filepath.Dir(rss.filePath), 0755); err != nil {
			return err
		}

		if err = os.WriteFile(rss.filePath, yamlData, 0600); err != nil {
			return err
		}
	}

	return nil
}

// GetCategories will return a list of all the names and descriptions of the categories
func (rss Rss) GetCategories() (names []string, descs []string) {
	names = make([]string, len(rss.Categories))
	descs = make([]string, len(rss.Categories))

	for i, cat := range rss.Categories {
		names[i] = cat.Name
		descs[i] = cat.Description
	}

	return names, descs
}

// GetFeeds will return a list of all the names and descriptions of the feeds
// in a category denoted by the name
func (rss Rss) GetFeeds(categoryName string) (names []string, urls []string, err error) {
	for _, cat := range rss.Categories {
		if cat.Name == categoryName {
			feeds := make([]string, len(cat.Subscriptions))
			urls = make([]string, len(cat.Subscriptions))

			for i, feed := range cat.Subscriptions {
				feeds[i] = feed.Name
				urls[i] = feed.URL
			}

			return feeds, urls, nil
		}
	}

	return nil, nil, ErrNotFound
}

// GetFeedURL will return the url of a feed denoted by the name
func (rss Rss) GetFeedURL(feedName string) (string, error) {
	if feedName == AllFeedsName || feedName == DownloadedFeedsName {
		return "", ErrReservedName
	}

	for _, cat := range rss.Categories {
		for _, feed := range cat.Subscriptions {
			if feed.Name == feedName {
				return feed.URL, nil
			}
		}
	}

	return "", ErrNotFound
}

// GetAllURLs will return a list of all the urls
func (rss Rss) GetAllURLs() []string {
	var urls []string

	for _, cat := range rss.Categories {
		for _, feed := range cat.Subscriptions {
			if feed.URL != AllFeedsName {
				urls = append(urls, feed.URL)
			}
		}
	}

	return urls
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

// LoadOPML will load the urls from an opml file.
func (rss *Rss) LoadOPML(path string) error {
	parsed, err := opml.NewOPMLFromFile(path)
	if err != nil {
		return err
	}

	for _, o := range parsed.Outlines() {
		catName := DefaultCategoryName
		catDesc := DefaultCategoryDescription

		if o.Type != "rss" {
			catName = o.Title
			catDesc = o.Text
		}

		if err = rss.AddCategory(catName, catDesc); err != nil && err != ErrAlreadyExists {
			return err
		}

		if len(o.Outlines) == 0 {
			if err = rss.AddFeed(DefaultCategoryName, o.Title, o.XMLURL); err != nil && err != ErrAlreadyExists {
				return err
			}

			continue
		}

		for _, so := range o.Outlines {
			log.Println("Adding feed:", so.Title)
			if err = rss.AddFeed(catName, so.Title, so.XMLURL); err != nil && err != ErrAlreadyExists {
				return err
			}
		}
	}

	return nil
}

// ExportOPML will export the urls to an opml file.
func (rss *Rss) ExportOPML(path string) error {
	result := opml.OPML{
		Version: "1.0",
		Head:    opml.Head{Title: "goread - Exported feeds"},
		Body:    opml.Body{},
	}

	for _, cat := range rss.Categories {
		result.Body.Outlines = append(result.Body.Outlines, opml.Outline{
			Title: cat.Name,
			Text:  cat.Description,
		})

		elem := &result.Body.Outlines[len(result.Body.Outlines) - 1]
		for _, feed := range cat.Subscriptions {
			elem.Outlines = append(elem.Outlines, opml.Outline{
				Type:   "rss",
				Text:   feed.Description,
				Title:  feed.Name,
				XMLURL: feed.URL,
			})
		}
	}

	data, err := result.XML()
	if err != nil {
		return err
	}

	if err = os.WriteFile(path, []byte(data), 0644); err != nil {
		return err
	}

	return nil
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

// getDefaultPath will return the default path for the urls file
func getDefaultPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "goread", "urls.yml"), nil
}
