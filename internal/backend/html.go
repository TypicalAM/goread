package backend

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Since the rss feed "content" is HTML, we need to parse it and get the text
// from it. This is a helper function to do that.
func parseHTML(content string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return ""
	}

	return doc.Text()
}
