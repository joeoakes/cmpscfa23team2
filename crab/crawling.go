package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/temoto/robotstxt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
	// Import the "dal" package to access its functions.
	"cmpscfa23team2/dal"
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
		fmt.Println("Found link:", link)
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			// Process the response, extract data, and add to URLData.
			urlData.Tags = make(map[string]interface{}) // Set tags to an empty map, you can fill this.
			ch <- urlData

			// Print the crawled URL data.
			fmt.Printf("Crawled URL: %s\n", urlData.URL)
			fmt.Printf("Tags: %+v\n", urlData.Tags)
			fmt.Printf("Domain: %s\n", urlData.Domain)
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

	for _, urlData := range urls {
		wg.Add(1) // Increment the WaitGroup counter

		// Call the crawlURL function as a goroutine
		go func(u URLData) {
			crawlURL(u, ch, &wg) // Call the crawlURL function as a goroutine
		}(urlData) // Pass the URLData to the goroutine

		// Check if the number of concurrent crawlers has reached the limit
		if len(urls) >= concurrentCrawlers {
			break // Limit the number of concurrent goroutines
		}
	}
	// Close the channel when all crawlers are done
	go func() {
		wg.Wait()
		close(ch)
	}()
}

func main() {
	urlsToCrawl := getURLsToCrawl() // Function to fetch URLs dynamically
	InitializeCrawling(urlsToCrawl) // Function to initialize crawling
}

func InitializeCrawling(urls []URLData) {
	// Fetch the list of URLs to crawl from the DB.
	urlsToCrawl, err := dal.GetURLsOnly() // Use the GetURLsOnly function from the "dal" package.
	// Handle the error as needed.
	if err != nil {
		log.Fatal("Error fetching URLs: ", err)
	}
	// Create a list of URLData objects.
	var urlDataList []URLData

	// Loop through the URLs and fetch the tags and domain for each URL.
	for _, url := range urlsToCrawl {
		tags, domain, err := dal.GetURLTagsAndDomain(url) // Use the GetURLTagsAndDomain function from the "dal" package.

		if err != nil {
			fmt.Printf("Error fetching data for URL %s: %v\n", url, err)
			continue
		}
		// Append the URLData object to the list.
		urlDataList = append(urlDataList, URLData{
			URL:    url,
			Tags:   tags,
			Domain: domain,
		})

		// Insert the URL into the database with associated tags.
		_, err = dal.InsertURL(url, domain, tags) // Use the InsertURL function from the "dal" package.
		if err != nil {
			fmt.Printf("Error inserting URL %s: %v\n", url, err)
		}
	}
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
	}
}
