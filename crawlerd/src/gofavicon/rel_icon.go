package gofavicon

import (
	"net/url"
	"strings"
	"sync"
)

type RelIcon struct {
	// favicon URL
	IconURL string
	// value from rel attr
	RelType string
	// value from sizes attr
	Sizes string
	//
	url  *url.URL
	lock sync.Mutex
}

// check if RelIcon contains absolute url
func (rel RelIcon) URL() *url.URL {
	if rel.url == nil {
		rel.lock.Lock()
		u, _ := url.Parse(rel.IconURL)
		rel.url = u
		rel.lock.Unlock()
	}
	return rel.url
}

func (rel RelIcon) IsAbsURL() bool {
	return rel.URL().IsAbs()
}

func (rel RelIcon) IsEmbedded() bool {
	return len(rel.URL().Opaque) > 0
}

func (rel RelIcon) EmbeddedType() (string, bool) {
	i := strings.Index(rel.URL().Opaque, ";")
	if i == -1 {
		return "", false
	}
	return rel.URL().Opaque[0:i], true
}
