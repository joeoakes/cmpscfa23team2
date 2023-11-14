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
		// Handle error
		return false
	}
	domain := parsedURL.Host

	// Correct URL for robots.txt
	robotsURL := "http://" + domain + "/robots.txt"

	resp, err := http.Get(robotsURL)
	if err != nil {
		// Handle the error as needed.
		return true
	}

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		// Handle the error as needed.
		return true
	}

	// Check if the URL is allowed
	return data.TestAgent(urlStr, "YourBotName") // Replace "YourBotName" with your actual bot name.
}

func threadedCrawl(urls []URLData, concurrentCrawlers int) {
	var wg sync.WaitGroup
	ch := make(chan URLData, len(urls))

	for _, urlData := range urls {
		wg.Add(1)
		go func(u URLData) {
			crawlURL(u, ch, &wg)
		}(urlData)
		if len(urls) >= concurrentCrawlers {
			break // Limit the number of concurrent goroutines
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
}

func main() {
	urlsToCrawl := getURLsToCrawl() // Function to fetch URLs dynamically
	InitializeCrawling(urlsToCrawl)
}

func InitializeCrawling(urls []URLData) {
	// Fetch the list of URLs to crawl from the DB.
	urlsToCrawl, err := dal.GetURLsOnly() // Use the GetURLsOnly function from the "dal" package.
	if err != nil {
		log.Fatal("Error fetching URLs: ", err)
	}

	var urlDataList []URLData
	for _, url := range urlsToCrawl {
		tags, domain, err := dal.GetURLTagsAndDomain(url) // Use the GetURLTagsAndDomain function from the "dal" package.
		if err != nil {
			fmt.Printf("Error fetching data for URL %s: %v\n", url, err)
			continue
		}
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

	threadedCrawl(urlDataList, 10) // Use 10 concurrent crawlers.
}

func getURLsToCrawl() []URLData {
	// For testing, we are using 'http://books.toscrape.com/'
	return []URLData{
		{
			URL: "http://books.toscrape.com/",
		},
	}
}
