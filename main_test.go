package goscraper

import (
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
