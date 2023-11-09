package main

import (
	"cmpscfa23team2/DAL"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/temoto/robotstxt"
	"log"
	"net/http"
	"sync"
	"time"
)

type URLData struct {
	URL     string
	Tags    map[string]interface{}
	Domain  string
	Created time.Time
}

func crawlURL(urlData URLData, ch chan<- URLData, wg *sync.WaitGroup) {
	defer wg.Done()

	// Check robots.txt before crawling
	allowed := isURLAllowedByRobotsTXT(urlData.URL)
	if !allowed {
		return
	}

	// Use colly library to fetch the content.
	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		// Handle errors while fetching content (e.g., HTTP errors)
		fmt.Printf("Error occurred while crawling %s: %s\n", urlData.URL, err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Process the link as needed.
		// For this example, we'll just print it.
		fmt.Println("Found link:", link)
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			// Process the response, extract data, and add to URLData.
			// For this example, we'll just mark it as visited and pass it.
			urlData.Tags = make(map[string]interface{}) // Set tags to an empty map, you can fill this.
			ch <- urlData
		} else {
			// Handle non-200 status codes
			fmt.Printf("Non-200 status code while crawling %s: %d\n", urlData.URL, r.StatusCode)
		}
	})

	c.Visit(urlData.URL)
}

func isURLAllowedByRobotsTXT(url string) bool {
	// Fetch and parse robots.txt for the domain using go-robotstxt library
	robotsURL := url + "/robots.txt"
	resp, err := http.Get(robotsURL)
	if err != nil {
		// Handle the error as needed. If there's no robots.txt, it's assumed that all paths are allowed.
		return true
	}

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		// Handle the error as needed. If parsing fails, it's assumed that all paths are allowed.
		return true
	}

	// Check if the URL is allowed
	return data.TestAgent(url, "YourBotName") // Replace "YourBotName" with your actual bot name.
}

func threadedCrawl(urls []URLData, concurrentCrawlers int) {
	var wg sync.WaitGroup
	ch := make(chan URLData, len(urls))

	for i := 0; i < concurrentCrawlers && i < len(urls); i++ {
		wg.Add(1)
		go crawlURL(urls[i], ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
}

func main() {
	InitializeCrawling()
}

func InitializeCrawling() {
	// Fetch the list of URLs to crawl from the DB.
	urlsToCrawl, err := GetURLsOnly()
	if err != nil {
		log.Fatal("Error fetching URLs: ", err)
	}

	var urlDataList []URLData
	for _, url := range urlsToCrawl {
		tags, domain, err := GetURLTagsAndDomain(url)
		if err != nil {
			fmt.Printf("Error fetching data for URL %s: %v\n", url, err)
			continue
		}
		urlDataList = append(urlDataList, URLData{
			URL:    url,
			Tags:   tags,
			Domain: domain,
		})
	}

	threadedCrawl(urlDataList, 10) // Example: Use 10 concurrent crawlers.
}

// Placeholder for fetching URLs from your data source.
func GetURLsOnly() ([]string, error) {
	// Replace this with your implementation to fetch URLs from your data source.
	// For testing purposes, return some example URLs.
	return []string{"https://example.com", "https://example.org"}, nil
}

// Placeholder for fetching tags and domain for a given URL.
func GetURLTagsAndDomain(url string) (map[string]interface{}, string, error) {
	// Replace this with your implementation to fetch tags and domain for a given URL.
	// For testing purposes, return some example data.
	tags := map[string]interface{}{"tag1": "value1", "tag2": "value2"}
	return tags, "example.com", nil
}
