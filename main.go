package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/haarshitgarg/go-crawler.git/internals"
)

func crawlPage(rawBaseURL string, rawCurrURL string, pages map[string]int) {
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

func printReport(pages map[string]int, baseURL string) {
	type kv struct {
		key string
		value int 
	}
	var kvSlice []kv 
	for k, v := range pages {
		kvSlice = append(kvSlice, kv{key: k, value: v})
	}

	sort.Slice(kvSlice, func(i, j int) bool {
		return kvSlice[i].value < kvSlice[j].value
	})

	fmt.Println("=============================")
  	fmt.Printf("REPORT for %s", baseURL)
	fmt.Println("=============================")
	for _, kv := range kvSlice {
		fmt.Printf("Found %d internal links to %s\n", kv.value, kv.key)
	}

}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	rawBaseURL := args[0]
	maxPage := args[1]
	maConc := args[2]
	nPage, err := strconv.Atoi(maxPage)
	if err != nil {
		fmt.Printf("The max page variable needs to be integer\n")
		os.Exit(1)
	}
	nConc, err := strconv.Atoi(maConc)
	if err != nil {
		fmt.Printf("The max conc variable needs to be integer\n")
		os.Exit(1)
	}


	baseURL, _ := url.Parse(rawBaseURL)
	pages := make(map[string]int)
	wg := sync.WaitGroup{}

	configuration := config{
		pages: pages,
		baseURL: *baseURL,
		mu: &sync.Mutex{},
		concurrencyControl: make(chan struct{}, nConc),
		wg: &wg,
		maxPage: nPage,
	}

	configuration.wg.Add(1)
	go configuration.crawlPage(rawBaseURL)
	configuration.wg.Wait()

	//crawlPage(rawBaseURL, rawBaseURL, pages)

	printReport(pages, rawBaseURL)
	os.Exit(0)
}
