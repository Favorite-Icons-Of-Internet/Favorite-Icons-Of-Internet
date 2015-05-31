package gofavicon

import (
	"regexp"
)

var (
	linkRe = regexp.MustCompile("(?i)<link([^>])+>")
	relRe = regexp.MustCompile(`(?i)rel="?(icon|shortcut icon)"?`)
	typeRe = regexp.MustCompile(`(?i)type="?([^"\s\>]+)"?`)
	hrefRe = regexp.MustCompile(`(?i)href="?([^"\s\>]+)"?`)
	sizesRe = regexp.MustCompile(`(?i)sizes="?([^"\s\>]+)"?`)
)

type regexpParser struct {}

func (parser regexpParser) Parse(html []byte) (*RelIcon, bool) {
	find := func(r *regexp.Regexp, body []byte) string {
		matches := r.FindSubmatch(body)
		if len(matches) >= 1 {
			return string(matches[1])
		}
		return ""
	}

	for _, match := range linkRe.FindAll(html, -1) {
		rel := find(relRe, match)
		typ := find(typeRe, match)
		href := find(hrefRe, match)
		sizes := find(sizesRe, match)

		if len(rel) > 0 && len(href) > 0 {

			ico := &RelIcon{
				IconURL: href,
				RelType: typ,
				Sizes: sizes,
			}

			return ico, true
		}
	}

	return nil, false
}

func NewReParser() *regexpParser {
	return &regexpParser{}
}