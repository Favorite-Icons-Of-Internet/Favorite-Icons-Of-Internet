package gofavicon

import (
	"io/ioutil"
	"net/http"
	"time"
)

const (
	httpTimeout = 10 * time.Second
)

var (
	client = &http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
		Timeout:   httpTimeout,
	}
)

// Fetch URL using HTTP GET
func Fetch(url string) ([]byte, error) {
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
