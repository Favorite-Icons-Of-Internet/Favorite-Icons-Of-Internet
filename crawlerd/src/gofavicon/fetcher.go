package gofavicon

import (
	"net/http"
	"io/ioutil"
)

var (
	transport = &http.Transport{DisableKeepAlives: true}
	client = &http.Client{Transport: transport}
)

// fetch URL
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