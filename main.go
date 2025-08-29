package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	//"github.com/haarshitgarg/go-crawler.git/internals"
)

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

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Printf("Got error: %s", err)
		os.Exit(1)
	}
	req.Header.Set("User-Agent", "gocrawler/1.0 (contact: hgarg9460@gmail.com)")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Got error: %s", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	

	cType := res.Header.Get("Content-Type")
	if !strings.Contains(cType, "text/html") {
		fmt.Printf("The Content is not html. It is of type: %s\n", res.Header.Get("Content-Type"))
		os.Exit(1)
	}

	htmlBody := res.Body
	buf, err := io.ReadAll(htmlBody)
	if err != nil {
		fmt.Printf("Got error: %s", err)
		os.Exit(1)
	}
	fmt.Println(string(buf))
	//urls, err := internals.GetURLsFromHTML(htmlBody, baseURL)
	//if err != nil {
	//	fmt.Printf("Got error: %s", err)
	//}
	//fmt.Println(urls)

	os.Exit(0)

}
