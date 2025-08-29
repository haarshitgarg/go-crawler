package internals

import (
	"fmt"
	"net/url"
	"strings"
)

func normaliseURL(u string) (string, error) {
	// Parse the URL and remove useless stuff
	pu, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("Could Not parse the given url: %s\nError: %s\n", u, err)
	}
	host := pu.Host
	path := pu.Path
	host = strings.TrimPrefix(host, "/")
	host = strings.TrimSuffix(host, "/")
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	normURL := fmt.Sprintf("%s/%s", host, path)

	return normURL, nil
}
