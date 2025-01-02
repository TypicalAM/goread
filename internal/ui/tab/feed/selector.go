package feed

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"

	"github.com/TypicalAM/goread/internal/theme"
	"github.com/charmbracelet/lipgloss"
	"mvdan.cc/xurls/v2"
)

// selector allows us to select links from a feed and open them in the browser
type selector struct {
	linkStyle lipgloss.Style
	article   *string
	urls      []string
	indices   [][]int
	selection int
	active    bool
}

// newSelector creates a new selector
func newSelector(colors *theme.Colors) *selector {
	return &selector{
		linkStyle: lipgloss.NewStyle().
			Background(colors.Color1).
			Underline(true),
	}
}

// newArticle finds the URLs for this article
func (s *selector) newArticle(rawText, styledText *string) {
	s.article = styledText
	s.selection = 0
	s.active = false

	rx := xurls.Strict()
	urlsToIndex := rx.FindAllString(*rawText, -1)

	// map the url to their possible linebreak indices
	urlsMap := make(map[string][]int)
	for _, url := range urlsToIndex {
		if !strings.ContainsRune(url, '-') {
			urlsMap[url] = append(urlsMap[url], len(url)-1)
			continue
		}

		for i := 0; i < len(url); i++ {
			if url[i] == '-' {
				urlsMap[url] = append(urlsMap[url], i)
			}
		}

		urlsMap[url] = append(urlsMap[url], len(url)-1)
	}

	s.urls = make([]string, 0)
	s.indices = make([][]int, 0)

	for url, indices := range urlsMap {
		s.urls = append(s.urls, url)
		// Check if the entire url fits in one line
		start := strings.Index(*styledText, url[:indices[len(indices)-1]])
		if start != -1 {
			s.indices = append(s.indices, []int{start, start + len(url)})
			continue
		}

		// Let's check on the - character on which the url is broken down
		for i := len(indices) - 2; i >= 0; i-- {
			start = strings.Index(*styledText, url[:indices[i]])
			if start == -1 {
				continue
			}

			// The line is broken down on index, let's search where it ends on the next line
			end := 0
			for j := start + indices[i]; j < len(*styledText); j++ {
				if (*styledText)[j] == url[indices[i]+1] {
					end = j + len(url) - indices[i] - 1
					break
				}
			}

			s.indices = append(s.indices, []int{start, end})
			break
		}

		// Special case: the URL which appeared in the original raw text doesn't appear at all in the yassified markdown, for example image descriptions also being images
		if len(s.urls) > 0 {
			s.urls = s.urls[:len(s.urls)-1]
		}
	}
}

// cycle triggers the selection of the next link in the feed
func (s *selector) cycle() string {
	var b strings.Builder

	s.selection++
	if !s.active || s.selection == len(s.urls) {
		s.selection = 0
		s.active = true
	}

	start, end := s.indices[s.selection][0], s.indices[s.selection][1]
	b.WriteString((*s.article)[:start])
	linkText := (*s.article)[start:end]

	// This is tricky
	if strings.ContainsRune(linkText, '\n') {
		newLine := strings.IndexRune(linkText, '\n')
		lastSpace := strings.LastIndex(linkText, " ")

		b.WriteString(s.linkStyle.Render(strings.TrimSpace(linkText[:lastSpace])))
		b.WriteString(linkText[newLine:lastSpace])
		b.WriteString(s.linkStyle.Render(linkText[lastSpace:]))
	} else {
		b.WriteString(s.linkStyle.Render(linkText))
	}

	b.WriteString((*s.article)[end:])
	return b.String()
}

// open opens the URL in the browser
func (s *selector) open() error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", s.urls[s.selection]).Start() //nolint:gosec
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", s.urls[s.selection]).Start() //nolint:gosec
	case "darwin":
		return exec.Command("open", s.urls[s.selection]).Start() //nolint:gosec
	default:
		return errors.New("unsupported platform")
	}
}
