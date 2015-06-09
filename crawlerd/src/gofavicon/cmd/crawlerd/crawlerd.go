package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"gofavicon"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Req struct {
	ID                string `json:"id"`
	Domain            string `json:"domain"`
	PreviousHash      string `json:"previous_hash"`
	PreviousFetchTime string `json:"previous_fetch_time"`
}

type Res struct {
	ID                string `json:"id"`
	Domain            string `json:"domain"`
	PreviousHash      string `json:"previous_hash"`
	PreviousFetchTime string `json:"previous_fetch_time"`
	Changed           bool   `json:"changed"`
	IconFile          string `json:"icon_file"`
	NewHash           string `json:"new_hash"`
	NewFetchTime      string `json:"new_fetch_time"`
}

// Extract favicon
func extract(r *Req) (*Res, error) {
	e := gofavicon.NewExtractor()
	url := r.Domain
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https") {
		url = fmt.Sprintf("http://%s", url)
	}

	ico, err := e.Extract(url)
	if err != nil {
		return nil, err
	}

	var changed bool

	hash := fmt.Sprintf("%x", sha1.Sum(ico.Image))

	if !strings.EqualFold(hash, r.PreviousHash) {
		changed = true
	}

	var filepath string
	if changed {
		file, _ := ioutil.TempFile(os.TempDir(), "")
		file.Write(ico.Image)
		filepath = file.Name()
	}

	res := &Res{
		Domain: r.Domain,
		ID: r.ID,
		PreviousHash: r.PreviousHash,
		PreviousFetchTime: r.PreviousFetchTime,
		NewFetchTime: time.Now().Format(time.RFC3339),
		NewHash: hash,
		Changed: changed,
		IconFile: filepath,
	}

	return res, nil
}

func processRequest(rch <-chan *Req, outCh chan<- *Res, ws *sync.WaitGroup) {
	for r := range rch {
		res, err := extract(r)
		if err != nil {

		}
		outCh <- res
	}
	ws.Done()
}

func processResult(res <-chan *Res) {
	for r := range res {
		bytes, _ := json.Marshal(&r)
		fmt.Println(string(bytes))
	}
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	reader.Split(TextSplit)

	var reqCh = make(chan *Req, 10)
	var outCh = make(chan *Res)

	var ws sync.WaitGroup

	go processResult(outCh)

	for i := 0; i < 10; i++ {
		ws.Add(1)
		go processRequest(reqCh, outCh, &ws)
	}

	for reader.Scan() {
		data := reader.Bytes()

		var requests []*Req
		err := json.Unmarshal(data, &requests)
		if err != nil {
			log.Print(err)
			continue
		}

		for _, r := range requests {
			reqCh <- r
		}
	}

	close(reqCh)
	ws.Wait()

	close(outCh)
}

// custom split function for reader
// read all from stdin until atEOF reached.
func TextSplit(data []byte, atEOF bool) (advanced int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
