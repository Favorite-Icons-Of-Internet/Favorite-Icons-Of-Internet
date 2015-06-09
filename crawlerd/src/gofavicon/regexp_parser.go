package gofavicon

import (
	"regexp"
)

var (
	// <link...> tag
	linkRe = regexp.MustCompile("(?i)<link([^>])+>")
	// rel attribute inside <link>
	relRe = regexp.MustCompile(`(?i)rel="?(icon|shortcut icon)"?`)
	// type attribute inside <link>
	typeRe = regexp.MustCompile(`(?i)type="?([^"\s\>]+)"?`)
	// href attribute inside <link>
	hrefRe = regexp.MustCompile(`(?i)href="?([^"\s\>]+)"?`)
	// sizes attribute inside <link>
	sizesRe = regexp.MustCompile(`(?i)sizes="?([^"\s\>]+)"?`)
)

type regexpParser struct{}

func findFirst(r *regexp.Regexp, body []byte) string {
	matches := r.FindSubmatch(body)
	if len(matches) >= 1 {
		return string(matches[1])
	}
	return ""
}

func NewReParser() *regexpParser {
	return &regexpParser{}
}

func (parser regexpParser) Parse(html []byte) (*RelIcon, bool) {
	for _, match := range linkRe.FindAll(html, -1) {
		rel := findFirst(relRe, match)
		typ := findFirst(typeRe, match)
		href := findFirst(hrefRe, match)
		sizes := findFirst(sizesRe, match)

		if len(rel) > 0 && len(href) > 0 {
			ico, err := NewRelIcon(href)
			if err != nil {
				return nil, false
			}
			ico.RelType = typ
			ico.Sizes = sizes
			return ico, true
		}
	}

	return nil, false
}
