package feed

import (
	"errors"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/TypicalAM/goread/internal/colorscheme"
	"github.com/charmbracelet/lipgloss"
	"mvdan.cc/xurls/v2"
)

var ansiRe = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

// selector allows us to select links from a feed and open them in the browser
type selector struct {
	linkStyle lipgloss.Style
	article   string
	urls      []string
	indices   [][]int
	selection int
	active    bool
}

// newSelector creates a new selector
func newSelector(colors colorscheme.Colorscheme) *selector {
	return &selector{
		linkStyle: lipgloss.NewStyle().
			Background(colors.Color1).
			Underline(true),
	}
}

// newArticle finds the URLs for this article
func (s *selector) newArticle(content string) {
	rx := xurls.Relaxed()

	s.article = content
	s.selection = 0
	s.active = false

	rawIndices := rx.FindAllStringIndex(content, -1)

	// Fix the newline issues
	s.indices = make([][]int, 0)
	s.urls = make([]string, 0)

	// Link highlighting is stupid, so we have to do this
	for i := 0; i < len(rawIndices); i++ {
		str := s.article[rawIndices[i][0]:rawIndices[i][1]]
		if str[len(str)-1] == '-' {
			s.indices = append(s.indices, []int{rawIndices[i][0], rawIndices[i+1][1]})
			urlNoAnsi := ansiRe.ReplaceAllString(s.article[rawIndices[i][0]:rawIndices[i+1][1]], "")
			urlStripped := strings.ReplaceAll(strings.ReplaceAll(urlNoAnsi, " ", ""), "\n", "")
			s.urls = append(s.urls, urlStripped)
			i++
		} else {
			s.indices = append(s.indices, rawIndices[i])
			s.urls = append(s.urls, s.article[rawIndices[i][0]:rawIndices[i][1]])
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
	linkText := s.article[start:end]
	b.WriteString(s.article[:start])

	// This is tricky
	if strings.ContainsRune(s.article[start:end], '\n') {
		linkText = ansiRe.ReplaceAllString(linkText, "")
		newLine := strings.IndexRune(linkText, '\n')
		lastSpace := strings.LastIndex(linkText, " ")

		b.WriteString(s.linkStyle.Render(strings.TrimSpace(linkText[:lastSpace])))
		b.WriteString(linkText[newLine:lastSpace])
		b.WriteString(s.linkStyle.Render(linkText[lastSpace:]))
	} else {
		b.WriteString(s.linkStyle.Render(linkText))
	}

	b.WriteString(s.article[end:])
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
