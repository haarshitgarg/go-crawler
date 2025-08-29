package internals

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)


func GetURLsFromHTML(r io.Reader, baseURL string) ([]string, error) {
	htmlNode, err := html.Parse(r)
	if err != nil {
		return make([]string, 0), err
	}

	urls := traverseURL(htmlNode, baseURL)

	return urls, nil
}

func traverseURL(node *html.Node, baseURL string) []string {
	if node == nil {
		return make([]string, 0)
	}
	var urls []string = make([]string, 0)
	if node.Data == "a" && node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				aUrl := attr.Val
				aUrl, err := validURL(aUrl, baseURL) 
				if err == nil {
					urls = append(urls, aUrl)
				} else {
					fmt.Printf("Error found: %s\n", err)
				}
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		urls = append(urls, traverseURL(c, baseURL)...)
	}

	return urls
}

func validURL(s string, basePath string) (string, error) {
	basePath = strings.TrimSuffix(basePath, "/")
	urlObj, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	if urlObj.Scheme != "http" && urlObj.Scheme != "https" {
		if strings.HasPrefix(s, "/") {
			return fmt.Sprintf("%s%s", basePath, s), nil
		}
		return "", fmt.Errorf("Not a valid URL anchor")
	}
	return  s, nil
}
