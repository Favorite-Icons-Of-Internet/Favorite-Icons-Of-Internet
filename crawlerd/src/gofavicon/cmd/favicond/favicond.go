package main

import (
	"bufio"
	"crypto/sha1"
	"expvar"
	"fmt"
	"gofavicon"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	processLog               *log.Logger
	queued, extracted, saved int32
)

func goroutines() interface{} {
	return runtime.NumGoroutine()
}

func StartSelfMonitor() {
	expvar.Publish("Goroutines", expvar.Func(goroutines))
	expvar.Publish("Queued", expvar.Func(func() interface{} {
		return queued
	}))
	expvar.Publish("Extracted", expvar.Func(func() interface{} {
		return extracted
	}))
	expvar.Publish("Saved", expvar.Func(func() interface{} {
		return saved
	}))

	go http.ListenAndServe(":12345", nil)
}

func extractListener(inCh <-chan string, out chan<- *gofavicon.Favicon, wg *sync.WaitGroup) {
	for url := range inCh {
		ico, err := gofavicon.Extract(url)
		if err != nil {
			log.Printf("%s - %s", url, err)
			continue
		}
		atomic.AddInt32(&extracted, 1)
		out <- ico
	}
	wg.Done()
}

func imageSaver(ch <-chan *gofavicon.Favicon, wg *sync.WaitGroup) {
	var tempDirs []string

	for i := 0; i < 20; i++ {
		tmpdir, _ := ioutil.TempDir("test", "")
		tempDirs = append(tempDirs, tmpdir)
	}

	for ico := range ch {
		h := sha1.New()
		io.WriteString(h, ico.SiteURL)
		dir := tempDirs[rand.Intn(len(tempDirs))]
		filename := fmt.Sprintf("%s/%x.ico", dir, h.Sum(nil))
		ioutil.WriteFile(filename, ico.Image, 0777)
		atomic.AddInt32(&saved, 1)
		fmt.Printf("PROCESSED\t%s\t%s\n", ico.SiteURL, filename)
	}

	wg.Done()
}

func main() {
	var (
		ch  chan string             = make(chan string, 1000)
		fav chan *gofavicon.Favicon = make(chan *gofavicon.Favicon, 100)
		wg  *sync.WaitGroup         = &sync.WaitGroup{}
	)

	StartSelfMonitor()

	var nrListeners int = 10
	wg.Add(nrListeners)
	for i := 0; i < nrListeners; i++ {
		go extractListener(ch, fav, wg)
	}

	wg.Add(1)
	go imageSaver(fav, wg)

	reader := bufio.NewScanner(os.Stdin)

	for reader.Scan() {
		url := reader.Text()

		if !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
			url = fmt.Sprintf("http://%s", url)
		}

		atomic.AddInt32(&queued, 1)
		ch <- url
	}

	close(ch)

	wg.Wait()
}
