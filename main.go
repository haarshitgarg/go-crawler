package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/haarshitgarg/go-crawler.git/internals"
)

func crawlPage(rawBaseURL string, rawCurrURL string, pages map[string]int) {
	fmt.Printf("Crawling page: %s\n", rawCurrURL)
	nCurrURL, err := internals.NormaliseURL(rawCurrURL)
	if err != nil {
		fmt.Printf("Encountered a URL that is not correct. Check it. URL: %s", rawCurrURL)
		fmt.Printf(". Ending with error: %s\n", err)
		return
	}
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Encountered a URL that cannot be parsed. Check it. URL: %s", rawCurrURL)
		fmt.Printf(". Ending with error: %s\n", err)
		return
	}
	currURL, err := url.Parse(rawCurrURL)
	if err != nil {
		fmt.Printf("Encountered a URL that cannot be parsed. Check it. URL: %s", rawCurrURL)
		fmt.Printf(". Ending with error: %s\n", err)
		return
	}

	val, ok := pages[nCurrURL]
	if ok {
		fmt.Printf("Found the page %s already. No need to check the links of this one\n", nCurrURL)
		pages[nCurrURL] = val + 1
		return
	} else {
		fmt.Printf("First time visiting the page %s\n", nCurrURL)
		pages[nCurrURL] = 1
	}

	if baseURL.Hostname() != currURL.Hostname() {
		return
	}

	htmlBody, err := getHtml(rawCurrURL)
	if err != nil {
		// TODO: Decide what to do here
		fmt.Printf("Error received: %s\n", err)
		return
	}

	fmt.Printf("Getting more URLs for the page %s\n", nCurrURL)
	urls, err := internals.GetURLsFromHTML(htmlBody, rawBaseURL)
	htmlBody.Close()
	for _, url := range urls {
		crawlPage(rawBaseURL, url, pages)
	}
}

func getHtml(baseURL string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gocrawler/1.0 (contact: hgarg9460@gmail.com)")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	cType := res.Header.Get("Content-Type")
	if !strings.Contains(cType, "text/html") {
		return nil, fmt.Errorf("The content type is not text/html. It is %s", cType) 
	}

	return res.Body, nil
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	baseURL := args[0]
	fmt.Printf("starting crawl of: %s\n", baseURL)

	pages := make(map[string]int)
	crawlPage(baseURL, baseURL, pages)
	fmt.Println("======================")
	fmt.Println(pages)
	os.Exit(0)

}
