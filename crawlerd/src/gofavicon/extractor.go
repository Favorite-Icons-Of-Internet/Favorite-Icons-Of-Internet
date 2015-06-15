package gofavicon

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Extractor struct {
	Parser FaviconParser
	Client *http.Client
}

type Content struct {
	Body     []byte
	Location *url.URL
}

type fetchError struct {
	URL string
	Reason error
}

func (err fetchError) Error() string {
	return fmt.Sprintf("Fetch failed:\n\tURL: %s\n\tError: %s\n", err.URL, err.Reason)
}


var DefaultHttpClient = &http.Client{
	Transport: &http.Transport{DisableKeepAlives: true},
	Timeout:   10 * time.Second,
}

var DefaultParser = NewReParser()

func NewExtractor() *Extractor {
	return &Extractor{
		Client: DefaultHttpClient,
		Parser: DefaultParser,
	}
}

func (e Extractor) Extract(u string) (*Favicon, error) {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https") {
		u = fmt.Sprintf("http://%s", u)
	}

	ico, err := e.ExtractFromURL(u)
	if err != nil {
		return nil, err
	}

	if ico != nil {
		return ico, nil
	}

	ico, err = e.ExtractDefault(u)
	if err != nil {
		return nil, err
	}

	return ico, err
}

func (e Extractor) ExtractFromURL(u string) (*Favicon, error) {
	page, err := e.Fetch(u)
	if err != nil {
		return nil, fetchError{u, err}
	}

	rel, ok := e.Parser.Parse(page.Body)

	if !ok {
		return nil, nil
	}

	if rel.IsEmbedded() {
		_, bytes, err := rel.Embedded()

		if err != nil {
			return nil, fetchError{rel.IconURL.String(), err}
		}

		i := &Favicon{
			ImageURL: rel.IconURL.String(),
			Image: bytes,
		}

		return i, nil
	}

	var faviconURL string

	if rel.IsAbsURL() {
		faviconURL = rel.IconURL.String()
	} else {
		p := *page.Location
		p.Path = rel.IconURL.Path
		faviconURL = p.String()
	}

	icon, err := e.Fetch(faviconURL)
	if err != nil {
		return nil, fetchError{u, err}
	}

	i := &Favicon{
		ImageURL: faviconURL,
		Image:    icon.Body,
	}

	return i, nil
}

func (e Extractor) ExtractDefault(s string) (*Favicon, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, fetchError{s, err}
	}

	faviconURL := DefaultFaviconURL(u)
	icon, err := e.Fetch(faviconURL)
	if err != nil {
		return nil, fetchError{s, err}
	}

	i := &Favicon{
		ImageURL: faviconURL,
		Image:    icon.Body,
	}

	return i, nil
}

// Generates URL for default location of favicon.ico
func DefaultFaviconURL(u *url.URL) string {
	r := url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   "/favicon.ico",
	}
	return r.String()
}

// Fetch URL using HTTP GET
func (e Extractor) Fetch(resource string) (*Content, error) {
	res, err := e.Client.Get(resource)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	loc, err := res.Location()
	if err != nil {
		loc, _ = url.Parse(resource)
	}

	page := &Content{body, loc}

	return page, nil
}
