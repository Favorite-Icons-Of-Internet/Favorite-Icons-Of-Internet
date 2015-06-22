package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"gofavicon"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	// path directory with downloaded icons
	outputDir string
	// path to executable icon processor
	iconProcessor string
)

type Req struct {
	ID                interface{} `json:"id"`
	Domain            string      `json:"domain"`
	PreviousHash      string      `json:"previous_hash"`
	PreviousFetchTime string      `json:"previous_fetch_time"`
}

type Res struct {
	ID                interface{} `json:"id"`
	Domain            string      `json:"domain"`
	PreviousHash      string      `json:"previous_hash"`
	PreviousFetchTime string      `json:"previous_fetch_time"`
	Changed           bool        `json:"changed"`
	IconFile          string      `json:"icon_file"`
	NewHash           string      `json:"new_hash"`
	NewFetchTime      string      `json:"new_fetch_time"`
}

// Extract favicon
func extract(r *Req) (*Res, error) {
	e := gofavicon.NewExtractor()

	d, err := resolveDomain(r.Domain, 10*time.Second)
	if err != nil {
		return nil, err
	}

	ico, err := e.Extract(d)
	if err != nil {
		return nil, err
	}

	var changed bool

	hash := fmt.Sprintf("%x", sha1.Sum(ico.Image))

	if !strings.EqualFold(hash, r.PreviousHash) {
		changed = true
	}

	var fpath string

	if changed {
		file, _ := ioutil.TempFile(outputDir, "")
		file.Write(ico.Image)
		fpath = file.Name() + ico.ImageExt
		if file.Name() != fpath {
			os.Rename(file.Name(), fpath)
		}
	}

	// run icon processor
	if len(iconProcessor) > 0 {
		abspath, err := filepath.Abs(fpath)
		if err != nil {
			return nil, err
		}

		b, err := execProcessor(abspath)
		if err != nil {
			return nil, fmt.Errorf("can't process file %s, error: %s, domain: %s", abspath, err, d)
		}

		fpath = strings.TrimSpace(string(b))
	}

	res := &Res{
		Domain:            d,
		ID:                r.ID,
		PreviousHash:      r.PreviousHash,
		PreviousFetchTime: r.PreviousFetchTime,
		NewFetchTime:      time.Now().Format(time.RFC3339),
		NewHash:           hash,
		Changed:           changed,
		IconFile:          fpath,
	}

	return res, nil
}

func processRequest(rch <-chan *Req, outCh chan<- *Res, ws *sync.WaitGroup, mon *selfMonitor) {
	for r := range rch {
		res, err := extract(r)

		if err != nil {
			log.Println(err)
			mon.AddFailed()
		} else {
			mon.AddProcessed()
			outCh <- res
		}
	}
	ws.Done()
}

func processResult(res <-chan *Res) {
	for r := range res {
		bytes, err := json.Marshal(&r)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("%s,\n", string(bytes))
		}
	}
}

func resolveDomain(d string, timeout time.Duration) (string, error) {
	hostExists := func(h string) bool {
		ch := lookupHost(h)
		select {
		case r := <-ch:
			return r
		case <-time.After(timeout):
			return false
		}
	}

	if hostExists(d) {
		return d, nil
	}

	w := fmt.Sprintf("www.%s", d)
	if hostExists(w) {
		return w, nil
	}

	return "", fmt.Errorf("domain %s not resolved", d)
}

func lookupHost(h string) chan bool {
	ch := make(chan bool)
	go func() {
		_, err := net.LookupHost(h)
		ch <- err == nil
		close(ch)
	}()
	return ch
}

func checkDirectory(p string) {
	err := os.MkdirAll(p, 0744)
	if err != nil {
		panic(err)
	}
}

func execProcessor(filename string) ([]byte, error) {
	args := strings.Split(iconProcessor, " ")
	args = append(args, filename)
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func init() {
	flag.StringVar(&outputDir, "output", "icons", "path to downloaded icons")
	flag.StringVar(&iconProcessor, "processor", "", "path to executable applied to each icon")
}

func main() {
	flag.Parse()

	checkDirectory(outputDir)

	var reqCh = make(chan *Req, 10)
	var outCh = make(chan *Res)

	var ws sync.WaitGroup

	monitor := NewSelfMonitor()
	monitor.Start()

	go processResult(outCh)

	ws.Add(10)
	for i := 0; i < 10; i++ {
		go processRequest(reqCh, outCh, &ws, monitor)
	}

	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		data := []byte(strings.TrimSuffix(reader.Text(), ","))

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
