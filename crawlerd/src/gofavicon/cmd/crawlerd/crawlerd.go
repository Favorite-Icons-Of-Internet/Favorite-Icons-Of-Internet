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
	ID                interface{} `json:"id"`
	Domain            string `json:"domain"`
	PreviousHash      string `json:"previous_hash"`
	PreviousFetchTime string `json:"previous_fetch_time"`
}

type Res struct {
	ID                interface{} `json:"id"`
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

	ico, err := e.Extract(r.Domain)
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
			log.Println(err)
			continue
		}
		outCh <- res
	}
	ws.Done()
}

func processResult(res <-chan *Res) {
	for r := range res {
		bytes, err := json.Marshal(&r)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(string(bytes))
		}
	}
}

func main() {
	var reqCh = make(chan *Req, 10)
	var outCh = make(chan *Res)

	var ws sync.WaitGroup

	go processResult(outCh)

	for i := 0; i < 10; i++ {
		ws.Add(1)
		go processRequest(reqCh, outCh, &ws)
	}

	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		data := []byte(reader.Text())

		var req *Req
		err := json.Unmarshal(data, &req)
		if err != nil {
			log.Println(err)
			continue
		}

		reqCh <- req
	}

	close(reqCh)
	ws.Wait()

	close(outCh)
}