package gofavicon_test

import (
	"gofavicon"
	"net/url"
	"strings"
	"testing"
)

func TestDefaultFaviconURL(t *testing.T) {
	urls := make(map[string]string)
	urls["http://host.com"] = "http://host.com/favicon.ico"
	urls["https://host.com"] = "https://host.com/favicon.ico"
	urls["https://host.com:12345"] = "https://host.com:12345/favicon.ico"
	urls["https://host.com:12345/y"] = "https://host.com:12345/favicon.ico"

	for website, expected := range urls {
		u, _ := url.Parse(website)
		ico := gofavicon.DefaultFaviconURL(u)
		if !strings.EqualFold(ico, expected) {
			t.Errorf("Unexpected default URL for favicon. Expected: %s, got: %s", expected, ico)
		}
	}

}

func TestExtractor(t *testing.T) {
	extractor := gofavicon.NewExtractor()
	ico, err := extractor.Extract("http://amazon.com")
	if err != nil {
		t.Error(err)
	}
	_ = ico
//	if len(ico.Image) == 0 {
//		t.Error("Image len is zero")
//	}
}