package main

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/haarshitgarg/go-crawler.git/internals"
)

type config struct {
	pages map[string]int
	baseURL url.URL 
	mu *sync.Mutex
	concurrencyControl chan struct{}
	wg *sync.WaitGroup
	maxPage int
}

func (c *config) crawlPage(rawCurrURL string) {
	c.concurrencyControl <- struct{}{}
	defer func () {
		<-c.concurrencyControl
		c.wg.Done()
	}()
	if c.maxPageReached() {
		return
	}

	nCurrURL, err := internals.NormaliseURL(rawCurrURL)
	if err != nil {
		fmt.Printf("Encountered a URL that is not correct. Check it. URL: %s", rawCurrURL)
		fmt.Printf(". Ending with error: %s\n", err)
		return
	}
	currURL, err := url.Parse(rawCurrURL)
	if err != nil {
		fmt.Printf("Encountered a URL that cannot be parsed. Check it. URL: %s", rawCurrURL)
		fmt.Printf(". Ending with error: %s\n", err)
		return
	}

	if c.addPageVisited(nCurrURL) {
		return
	}

	if c.baseURL.Hostname() != currURL.Hostname() {
		return
	}

	htmlBody, err := getHtml(rawCurrURL)
	if err != nil {
		// TODO: Decide what to do here
		fmt.Printf("Error received: %s\n", err)
		return
	}

	urls, err := internals.GetURLsFromHTML(htmlBody, c.baseURL.String())
	htmlBody.Close()
	for _, url := range urls {
		c.wg.Add(1)
		go c.crawlPage(url)
	}
}

func (c *config) addPageVisited(normURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.pages[normURL]
	if ok {
		c.pages[normURL] = val + 1
		return true
	}
	c.pages[normURL] = 1
	return false
}

func (c *config) maxPageReached() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.pages) >= c.maxPage {
		return true
	}
	return  false
}
