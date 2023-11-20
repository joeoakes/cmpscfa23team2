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

// URLData holds information about each URL to be crawled.
type URLData struct {
	URL     string    // The URL to be crawled
	Created time.Time // Timestamp of URL creation or retrieval
	Links   []string
}

// crawlURL is responsible for crawling a single URL.
func crawlURL(urlData URLData, ch chan<- URLData, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure the WaitGroup counter is decremented on function exit
	c := colly.NewCollector(
		colly.UserAgent(getRandomUserAgent()), // Set a random user agent
	)
	// First, check if the URL is allowed by robots.txt rules
	allowed := isURLAllowedByRobotsTXT(urlData.URL)
	if !allowed {
		return // Skip crawling if not allowed
	}

	// Handler for errors during the crawl
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error occurred while crawling %s: %s\n", urlData.URL, err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		urlData.Links = append(urlData.Links, link)
	})

	// Handler for anchor tags found in HTML
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println("Found link:", link)
		// Here you can enqueue the link for further crawling or processing
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

// isURLAllowedByRobotsTXT checks if the given URL is allowed by the site's robots.txt.
func isURLAllowedByRobotsTXT(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Println("Error parsing URL:", err)
		return false
	}

	domain := parsedURL.Host
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

	return data.TestAgent(urlStr, "GoEngine")
}

// threadedCrawl starts crawling the provided URLs concurrently.
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
				colly.UserAgent(getRandomUserAgent()),
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

// InitializeCrawling sets up and starts the crawling process.
func InitializeCrawling() {
	log.Println("Fetching URLs to crawl...")
	urlDataList := getURLsToCrawl()
	log.Println("URLs to crawl:", urlDataList)

	threadedCrawl(urlDataList, 10)
}

// getURLsToCrawl retrieves a list of URLs to be crawled.
func getURLsToCrawl() []URLData {
	return []URLData{
		{URL: "https://www.kaggle.com/search?q=housing+prices"},
		{URL: "http://books.toscrape.com/"},
		{URL: "https://www.kaggle.com/search?q=stocks"},
		{URL: "https://www.kaggle.com/search?q=stock+market"},
		{URL: "https://www.kaggle.com/search?q=real+estate"},
	}
}

func main() {
	InitializeCrawling()
}
