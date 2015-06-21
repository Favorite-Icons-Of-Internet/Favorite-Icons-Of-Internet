package gofavicon

import (
	"encoding/base64"
	"errors"
	"net/url"
	"path"
	"strings"
)

var ErrUnknownFormat = errors.New("Unknown dataurl format")

type RelIcon struct {
	// favicon URL
	IconURL *url.URL
	// value from rel attr
	RelType string
	// value from sizes attr
	Sizes string
}

// Initialize new RelIcon from string s
func NewRelIcon(s string) (*RelIcon, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	i := &RelIcon{IconURL: u}
	return i, nil
}

// Check if rel contains absolute url to favicon
func (rel RelIcon) IsAbsURL() bool {
	return rel.IconURL.IsAbs()
}

// Check if rel contains data (i.e. data:image/png;base64...)
func (rel RelIcon) IsEmbedded() bool {
	return rel.IconURL.Scheme == "data" && len(rel.IconURL.Opaque) > 0
}

// Extract MimeType and content bytes from data url
func (rel RelIcon) Embedded() (string, []byte, error) {
	i := strings.Index(rel.IconURL.Opaque, ";")
	j := strings.Index(rel.IconURL.Opaque, ",")
	if i == -1 || j == -1 {
		return "", nil, ErrUnknownFormat
	}

	mimeType := rel.IconURL.Opaque[0:i]
	body := rel.IconURL.Opaque[j+1:]

	bytes, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return "", nil, err
	}

	return mimeType, bytes, nil
}

// Extract file extension from url or empty string if extension is not set
func (rel RelIcon) Ext() string {
	return path.Ext(rel.IconURL.Path)
}
