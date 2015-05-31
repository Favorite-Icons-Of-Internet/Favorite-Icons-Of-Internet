package gofavicon

type FaviconParser interface {
	Parse(html []byte) (*RelIcon, bool)
}
