package gofavicon

import (
	"errors"
	"gofavicon/remote"
	"net/url"
)

// Favicon for site
type Favicon struct {
	// url of website
	SiteURL string
	// image url
	ImageURL string
	// type of image
	MimeType string
	// Icon width in pixels
	Width int
	// Icon height in pixels
	Height int
	// Icon
	Image []byte
}

type Resource string

var (
	errNoRel  = errors.New("favicon is not defined in link tag")
	errNoFile = errors.New("cant download file")
)

// extract favicon by URL
func Extract(res string) (*Favicon, error) {
	body, err := Fetch(res)
	if err != nil {
		return nil, err
	}

	parser := NewReParser()
	relIcon, ok := parser.Parse(body)
	if !ok {
		relIcon = NewRelIcon(faviconURL(res))
		//return nil, errNoRel
	}

	var icoUrl string

	if !relIcon.IsAbsURL() {
		u := Resource(res).Parsed()
		u.Path = relIcon.IconURL
		icoUrl = u.String()
	} else {
		icoUrl = relIcon.IconURL
	}

	file, err := remote.Get(icoUrl)
	if err != nil {
		return nil, err
	}

	i := &Favicon{
		SiteURL:  res,
		ImageURL: relIcon.IconURL,
		MimeType: file.MimeType,
		Width:    0,
		Height:   0,
		Image:    file.Content,
	}

	return i, nil
}

func (r Resource) Parsed() *url.URL {
	u, _ := url.Parse(string(r))
	return u
}

func faviconURL(resource string) string {
	u, _ := url.Parse(resource)
	r := url.URL{}
	r.Scheme = u.Scheme
	r.Host = u.Host
	r.Path = "/favicon.ico"
	return r.String()
}

func NewRelIcon(iconUrl string) *RelIcon {
	return &RelIcon{IconURL: iconUrl}
}
