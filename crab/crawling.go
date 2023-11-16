package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/temoto/robotstxt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type URLData struct {
	URL     string
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
		fmt.Println("Found link:", link)
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			// Process the response, extract data, and add to URLData.
			ch <- urlData

			// Print the crawled URL data.
			fmt.Printf("Crawled URL: %s\n", urlData.URL)
		} else {
			// Handle non-200 status codes
			fmt.Printf("Non-200 status code while crawling %s: %d\n", urlData.URL, r.StatusCode)
		}
	})

	c.Visit(urlData.URL)
}

func isURLAllowedByRobotsTXT(urlStr string) bool {
	// Parse domain from URL
	parsedURL, err := url.Parse(urlStr)

	if err != nil {
		log.Println("Error parsing URL:", err)
		return false
	}
	domain := parsedURL.Host

	// Correct URL for robots.txt
	robotsURL := "http://" + domain + "/robots.txt"

	resp, err := http.Get(robotsURL)
	if err != nil {
		log.Println("Error fetching robots.txt:", err)
		return true
	}

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		log.Println("Error parsing robots.txt:", err)
		return true
	}

	// Check if the URL is allowed
	return data.TestAgent(urlStr, "GoEngine")
}

func threadedCrawl(urls []URLData, concurrentCrawlers int) {
	var wg sync.WaitGroup               // WaitGroup to wait for all crawlers to finish
	ch := make(chan URLData, len(urls)) // Channel to receive crawled URL data
	log.Println("Starting crawling...")
	for _, urlData := range urls {
		wg.Add(1) // Increment the WaitGroup counter

		// Call the crawlURL function as a goroutine
		go func(u URLData) {
			crawlURL(u, ch, &wg) // Call the crawlURL function as a goroutine
		}(urlData) // Pass the URLData to the goroutine
		log.Println("Crawling URL:", urlData.URL)
		// Check if the number of concurrent crawlers has reached the limit
		if len(urls) >= concurrentCrawlers {
			break // Limit the number of concurrent goroutines
		}
	}
	log.Println("Waiting for crawlers to finish...")
	// Close the channel when all crawlers are done
	go func() {
		wg.Wait()
		close(ch)
	}()
	log.Println("Crawling finished!")
	writeCrawledURLsToFile(urls)
}

func main() {
	InitializeCrawling() // Function to initialize crawling

}

func InitializeCrawling() {
	// Fetch the list of URLs to crawl from the DB.
	// Get the list of URLs to crawl.
	log.Println("Fetching URLs to crawl...")
	urlDataList := getURLsToCrawl()
	log.Println("URLs to crawl:", urlDataList)
	// Call the threadedCrawl function to crawl the URLs.
	threadedCrawl(urlDataList, 10) // Use 10 concurrent crawlers.
}

// Function to fetch URLs dynamically.
func getURLsToCrawl() []URLData {
	// For testing, we are using 'http://books.toscrape.com/'
	return []URLData{
		{
			URL: "http://books.toscrape.com/",
		},
		{
			URL: "http://www.trulia.com",
		},
		{
			URL: "http://www.zillow.com",
		},
		{
			URL: "http://www.realtor.com",
		},
		{
			URL: "http://www.redfin.com",
		},
		{
			URL: "http://www.craigslist.com",
		},
		{
			URL: "http://www.propertyguru.com",
		},
	}
}

func writeCrawledURLsToFile(urls []URLData) error {
	urlStrings := make([]string, len(urls))
	for i, u := range urls {
		urlStrings[i] = u.URL
	}

	jsonData, err := json.Marshal(urlStrings)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("crawledUrls.json", jsonData, 0644)
}
