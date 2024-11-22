// Package goscraper implements simple utility methods to gather and query the
// go-index. (https://index.golang.org/)
//
// This package is not meant to support any real-time applications and rather
// intended to be a pre-cursor to any use you would have for a list of all go packages
// with version information.
package goscraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// A snapshot of a packages version and its version timestamp (RFC3339Nano)
type GoPackageVersion struct {
	Version   string
	Timestamp string
}

// Container struct for a single go package from the index
type GoPackage struct {
	Path     string
	Versions []GoPackageVersion
}

// Raw struct to store JSON data per entry of the index.golang.org API
type GoIndexEntry struct {
	Path      string `json:"Path"`
	Version   string `json:"Version"`
	Timestamp string `json:"Timestamp"`
}

// GenerateUrls returns a list of all URLS from index.golang.org to be scraped
// since a startTime and up to an endTime (both RFC3339Nano format)
//
// StartTime can be before 2019-04-10T19:08:52.997264Z but will fetch no additional data.
//
// Pass an empty string for either to default to the max range
func GenerateUrls(startTime string, endTime string) ([]string, error) {
	indexStartTime := "2019-04-10T19:08:52.997264Z"
	if startTime == "" {
		startTime = indexStartTime
	}
	parsedStartTime, err := time.Parse(time.RFC3339Nano, startTime)
	if err != nil {
		return make([]string, 0), err
	}

	if endTime == "" {
		endTime = time.Now().Format(time.RFC3339Nano)
	}
	parsedEndTime, err := time.Parse(time.RFC3339Nano, endTime)
	if err != nil {
		return make([]string, 0), err
	}

	// NOTE: The step could be an arg but may cause missing indexs if too great or slow performance if too small
	step := time.Duration(12) * time.Hour
	baseUrl := "https://index.golang.org/index?since="
	urls := make([]string, 0)

	for parsedStartTime.Unix() < parsedEndTime.Unix() {
		urls = append(urls, baseUrl+parsedStartTime.Format(time.RFC3339Nano))
		parsedStartTime = parsedStartTime.Add(step)
	}

	return urls, nil
}

// ProcessUrls returns a list of GoPackage's from the provided urls
//
// You can provide the size of semaphore (maxWorkers). For average machines recommended is 10-20
//
// Optional logging flag as well.
func ProcessUrls(urls []string, maxWorkers int, logProgress bool, wg *sync.WaitGroup) {
	sem := make(chan int, maxWorkers)
	urlCount := len(urls)
	indexEntries := make(chan *GoIndexEntry, (urlCount * 2000))

	allPackages := make(map[string]*GoPackage, 0)

	for i, url := range urls {
		if logProgress {
			fmt.Printf("\r Processing url: %d / %d", i, urlCount)
		}

		sem <- 1
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(url)
			// TODO: Collection/handling of bad urls
			if err != nil {
				<-sem
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				<-sem
				return
			}

			lines := strings.Split(string(body), "\n")
			for _, e := range lines {
				gie := &GoIndexEntry{}
				json.Unmarshal([]byte(e), gie)
				if len(gie.Path) < 5 {
					<-sem
					return
				}
				indexEntries <- gie
			}
			<-sem
		}()
	}

	wg.Wait()
	close(indexEntries)

	for e := range indexEntries {
		sem <- 1
		wg.Add(1)
		go func() {
			defer wg.Done()
			// TODO: Store/sort entries
		}()
	}
}
