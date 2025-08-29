package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormaliseURL(t *testing.T) {
	fmt.Println("Testing Normalise URL")
	url := "https://blog.boot.dev/path/"
	normUrl, err := normaliseURL(url)
	assert.NoError(t, err)
	assert.Equal(t, normUrl, "blog.boot.dev/path")

	url = "https://blog.boot.dev/path"
	normUrl, err = normaliseURL(url)
	assert.NoError(t, err)
	assert.Equal(t, normUrl, "blog.boot.dev/path")

	url = "http://blog.boot.dev/path/"
	normUrl, err = normaliseURL(url)
	assert.NoError(t, err)
	assert.Equal(t, normUrl, "blog.boot.dev/path")


	url = "http://blog.boot.dev/path"
	normUrl, err = normaliseURL(url)
	assert.NoError(t, err)
	assert.Equal(t, normUrl, "blog.boot.dev/path")
}
