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
	"os"
	"sync"
	"time"
)

// InitializeCrawling starts the web crawling process. It first fetches URLs to crawl from a predefined list,
// and then initiates a threaded crawl process with a specified number of concurrent crawlers.
func InitializeCrawling() {
	log.Println("Fetching URLs to crawl...")
	urlDataList := getURLsToCrawl()
	log.Println("URLs to crawl:", urlDataList)
	threadedCrawl(urlDataList, 10)
}

// / getURLsToCrawl returns a slice of URLData representing a list of URLs to be crawled.
// This function is used internally within the InitializeCrawling function.
func getURLsToCrawl() []URLData {
	return []URLData{
		{URL: "https://www.kaggle.com/search?q=housing+prices"},
		{URL: "http://books.toscrape.com/"},
		{URL: "https://www.kaggle.com/search?q=stocks"},
		{URL: "https://www.kaggle.com/search?q=stock+market"},
		{URL: "https://www.kaggle.com/search?q=real+estate"},
	}
}

// InsertData takes structured data (ItemData) and a filename, marshals the data into JSON format,
// and writes it to the specified file. It returns an error if any occurs during the marshaling or file operations.
func InsertData(data ItemData, filename string) error {
	// Save data to JSON file
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// crawlURL is the core function responsible for crawling a single URL. It takes URLData, a channel to send
// crawled data, and a WaitGroup to handle concurrency. It uses the Colly library for crawling and processes
// each URL based on the received HTML content.
func crawlURL(urlData URLData, ch chan<- URLData, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure the WaitGroup counter is decremented on function exit
	c := colly.NewCollector(
		colly.UserAgent(GetRandomUserAgent()), // Set a random user agent
		colly.AllowURLRevisit(),               // Allow URL revisit
	)

	// Handler for errors during the crawl
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error occurred while crawling %s: %s\n", urlData.URL, err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		urlData.Links = append(urlData.Links, link)
		urlQueue <- link
	})

	// Handler for successful HTTP responses
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			// Successful crawl, process the response here
			ch <- urlData // Send the URLData to the channel
			fmt.Printf("Crawled URL: %s\n", urlData.URL)
		} else {
			// Handle cases where the status code is not 200
			fmt.Printf("Non-200 status code while crawling %s: %d\n", urlData.URL, r.StatusCode)
		}
	})

	// Start the crawl
	c.Visit(urlData.URL)

	ch <- urlData
}

// createSiteMap generates a sitemap from the given slice of URLData. Each URLData contains links found
// at a specific URL. The function marshals this data into JSON format and writes it to a file named "siteMap.json".
// It returns an error if the marshaling or file operations fail.
func createSiteMap(urls []URLData) error {
	siteMap := make(map[string][]string)
	for _, u := range urls {
		siteMap[u.URL] = u.Links
	}

	jsonData, err := json.Marshal(siteMap)
	err = ioutil.WriteFile("siteMap.json", jsonData, 0644)
	if err != nil {
		log.Printf("Error writing sitemap to file: %v\n", err)
		return err
	}

	log.Println("Sitemap created successfully.")
	return nil
}

// isURLAllowedByRobotsTXT checks if the given URL is allowed by the site's robots.txt file.
// It parses the URL to extract the domain, fetches the robots.txt file from the domain, and tests
// if the URL is allowed. It returns true if allowed, false otherwise.
func isURLAllowedByRobotsTXT(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Println("Error parsing URL:", err)
		return false
	}

	if parsedURL.Host == "" {
		log.Println("Invalid URL, no host found:", urlStr)
		return false
	}

	robotsURL := "http://" + parsedURL.Host + "/robots.txt"

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

	return data.TestAgent(urlStr, "GoEngine")
}

//end robot.txt ========================================================================================================

// threadedCrawl manages the concurrent crawling of multiple URLs. It takes a slice of URLData and
// an integer specifying the number of concurrent crawlers. The function sets up each crawler with rate limiting
// and starts the crawling process. The resulting crawled data is used to create a sitemap.
func threadedCrawl(urls []URLData, concurrentCrawlers int) {
	var wg sync.WaitGroup
	ch := make(chan URLData, len(urls))

	rateLimitRule := &colly.LimitRule{
		DomainGlob:  "*",             // Apply to all domains
		Delay:       5 * time.Second, // Wait 5 seconds between requests
		RandomDelay: 5 * time.Second, // Add up to 5 seconds of random delay
	}

	log.Println("Starting crawling...")
	for _, urlData := range urls {
		wg.Add(1)

		go func(u URLData) {
			c := colly.NewCollector(
				colly.UserAgent(GetRandomUserAgent()),
			)
			c.Limit(rateLimitRule) // Set the rate limit rule

			crawlURL(u, ch, &wg)
		}(urlData)

		log.Println("Crawling URL:", urlData.URL)
		if len(urls) >= concurrentCrawlers {
			break
		}
	}

	log.Println("Waiting for crawlers to finish...")
	go func() {
		wg.Wait()
		close(ch)
		log.Println("All goroutines finished, channel closed.")
	}()

	var crawledURLs []URLData
	for urlData := range ch {
		crawledURLs = append(crawledURLs, urlData)
	}
	if err := createSiteMap(crawledURLs); err != nil {
		log.Println("Error creating sitemap:", err)
	}
}
