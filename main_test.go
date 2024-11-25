package goscraper

import (
	"fmt"
	"sync"
	"testing"
)

// Tests GenerateUrls without passing a startTime and endTime
func TestGenerateUrlsEmpty(t *testing.T) {
	startTime := ""
	endTime := ""
	urls, err := GenerateUrls(startTime, endTime)

	if len(urls) < 100 || err != nil {
		t.Fatalf("Expected a list of urls > 100, only recieved %v, error = %v", len(urls), err)
	}
}

// Tests TestGenerateUrls with a provided startTime and endTime
func TestGenerateUrlsProvided(t *testing.T) {
	startTime := "2022-04-10T19:08:52.997264Z"
	endTime := "2023-04-10T19:08:52.997264Z"
	urls, err := GenerateUrls(startTime, endTime)

	if len(urls) < 100 || err != nil {
		t.Fatalf("Expected a list of urls > 100, only recieved %v, error = %v", len(urls), err)
	}
}

func TestUrlProcessing(t *testing.T) {
	startTime := "2022-04-10T19:08:52.997264Z"
	endTime := "2023-04-10T19:08:52.997264Z"
	urls, _ := GenerateUrls(startTime, endTime)

	var wg sync.WaitGroup
	gi := ProcessUrls(urls, 10, true, &wg)

	// TODO: Better testing eventually but this shows that we are propogating packages and versions
	i := 0
	for path, pkgInfo := range gi.Packages {
		if i < 5 {
			i++
		} else {
			break
		}
		fmt.Println(path)
		fmt.Println(pkgInfo.Versions)
	}
}
