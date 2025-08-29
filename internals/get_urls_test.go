package internals

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLLinkExtraction(t *testing.T) {
	buf, err := os.ReadFile("../htmls/test_html1.html")
	require.NoError(t, err)
	r := bytes.NewReader(buf)
	urls, err := getURLsFromHTML(r, "https://blog.boot.dev")
	require.NoError(t, err)
	assert.Equal(t, []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"}, urls)

	buf, err = os.ReadFile("../htmls/test_html2.html")
	require.NoError(t, err)
	r = bytes.NewReader(buf)
	urls, err = getURLsFromHTML(r, "https://blog.boot.dev")
	require.NoError(t, err)
	assert.Equal(t, []string{
		"https://blog.boot.dev/home",
		"https://blog.boot.dev/about",
		"https://blog.boot.dev/contact",
	}, urls)

	buf, err = os.ReadFile("../htmls/test_html3.html")
	require.NoError(t, err)
	r = bytes.NewReader(buf)
	urls, err = getURLsFromHTML(r, "https://blog.boot.dev")
	require.NoError(t, err)
	assert.Equal(t, []string{
		"https://example.com",
		"https://blog.boot.dev/products",
		"https://google.com/search",
	}, urls)

	buf, err = os.ReadFile("../htmls/test_html4.html")
	require.NoError(t, err)
	r = bytes.NewReader(buf)
	urls, err = getURLsFromHTML(r, "https://blog.boot.dev")
	require.NoError(t, err)
	assert.Equal(t, []string{
		"https://blog.boot.dev/profile?id=123",
		"https://blog.boot.dev/settings#privacy",
		"https://shop.com/cart?item=book&qty=2",
	}, urls)

	buf, err = os.ReadFile("../htmls/test_html5.html")
	require.NoError(t, err)
	r = bytes.NewReader(buf)
	urls, err = getURLsFromHTML(r, "https://blog.boot.dev")
	require.NoError(t, err)
	assert.Equal(t, []string{
		"https://blog.boot.dev/blog",
		"https://blog.boot.dev/faq",
		"https://news.com/latest",
	}, urls)

	buf, err = os.ReadFile("../htmls/test_html6.html")
	require.NoError(t, err)
	r = bytes.NewReader(buf)
	urls, err = getURLsFromHTML(r, "https://blog.boot.dev")
	require.NoError(t, err)
	assert.Equal(t, []string{}, urls)

}
