package remote

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type RemoteFile struct {
	Content  []byte
	MimeType string
}

var (
	errUnknownFormat = errors.New("Unknown dataurl format")

	transport = &http.Transport{DisableKeepAlives: true}
	client    = &http.Client{Transport: transport}
)

func Get(u string) (*RemoteFile, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme == "data" {
		return extract(parsed)
	} else {
		return download(parsed)
	}
}

// download file
func download(u *url.URL) (*RemoteFile, error) {
	res, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	r := &RemoteFile{bytes, res.Header.Get("Content-Type")}

	return r, nil
}

// extract file from data url
func extract(u *url.URL) (*RemoteFile, error) {
	i := strings.Index(u.Opaque, ";")
	j := strings.Index(u.Opaque, ",")
	if i == -1 || j == -1 {
		return nil, errUnknownFormat
	}

	mimeType := u.Opaque[0:i]
	body := u.Opaque[j+1:]

	bytes, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, err
	}

	file := &RemoteFile{Content: bytes, MimeType: mimeType}

	return file, nil
}
