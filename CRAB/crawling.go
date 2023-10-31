package main

import (
	"fmt"
	"github.com/temoto/robotstxt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var visitedURLs = make(map[string]bool)
var visitedURLsMutex sync.Mutex

func isVisited(url string) bool {
	visitedURLsMutex.Lock()
	defer visitedURLsMutex.Unlock()
	return visitedURLs[url]
}

func markVisited(url string) {
	visitedURLsMutex.Lock()
	defer visitedURLsMutex.Unlock()
	visitedURLs[url] = true
}

func normalizeURL(baseURL, href string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return href
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return href
	}
	normalized := base.ResolveReference(uri)
	return normalized.String()
}

func extractLinksFromPage(doc *goquery.Document, err error) []string {
	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			links = append(links, href)
		}
	})
	return links
}

type ScrapeResult struct {
	URL   string
	Price string
	Error error
	Links []string
}

func scrape(url string, ch chan<- ScrapeResult, userAgent string, delay time.Duration) {
	if isVisited(url) {
		return // Skip already visited URLs
	}
	markVisited(url) // Mark the URL as visited

	time.Sleep(delay) // Rate limiting

	resp, err := http.Get(url)
	if err != nil {
		ch <- ScrapeResult{URL: url, Error: err}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ch <- ScrapeResult{URL: url, Error: fmt.Errorf("status code: %d", resp.StatusCode)}
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		ch <- ScrapeResult{URL: url, Error: err}
		return
	}

	price := doc.Find(".house-price-class").Text() // Assuming a class that holds the price
	links := extractLinksFromPage(doc, nil)
	ch <- ScrapeResult{URL: url, Price: price, Links: links}
}

func fetchRobotsAndExtractLinks(url string) (*robotstxt.Group, []string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		return nil, nil, err
	}

	links := extractLinksFromPage(goquery.NewDocumentFromReader(resp.Body))
	return data.FindGroup("*"), links, nil
}

//// InsertScrapedData inserts scraped data into the database using a stored procedure.
//func InsertScrapedData(url, data string) error {
//	_, err := dal.Exec("CALL InsertScrapedData(?, ?)", url, data)
//	return err
//}
//
//// InsertVisitedURL inserts visited URLs into the database using a stored procedure.
//func InsertVisitedURL(url string) error {
//	_, err := dal.Exec("CALL InsertVisitedURL(?)", url)
//	return err
//}

func main() {
	baseURL := "http://example-estate-site.com"
	robotsURL := baseURL + "/robots.txt"
	userAgent := "YourBotName"         // replace with your bot's name
	startingURL := baseURL + "/house1" // Replace with your starting URL
	urlsToCrawl := []string{startingURL}

	// Fetch robots.txt and extract links
	robotsGroup, _, err := fetchRobotsAndExtractLinks(robotsURL)
	if err != nil {
		log.Fatalf("Failed fetching robots.txt and extracting links: %v", err)
	}

	// Ensure our bot can access the desired URLs
	for _, url := range urlsToCrawl {
		if !robotsGroup.Test(url) {
			log.Fatalf("Access denied by robots.txt for URL: %s", url)
		}
	}

	ch := make(chan ScrapeResult)
	var wg sync.WaitGroup

	delay := 1 * time.Second // 1 request per second. Adjust as needed.
	for _, url := range urlsToCrawl {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			scrape(u, ch, userAgent, delay)
		}(url)
	}

	go func() {
		for result := range ch {
			if result.Error != nil {
				log.Printf("Error scraping %s: %v", result.URL, result.Error)
			} else {
				fmt.Printf("Scraped %s: %s\n", result.URL, result.Price)
				// Normalize and process links here.
				for _, link := range result.Links {
					normalizedLink := normalizeURL(result.URL, link)
					fmt.Printf("Found link: %s\n", normalizedLink)
					// You can add logic here to follow and crawl these links as needed.

					// Insert scraped data and visited URL into the database
					//err := InsertScrapedData(normalizedLink, result.Price)
					//if err != nil {
					//	log.Printf("Error inserting scraped data into the database: %v", err)
					//}
					//
					//err = InsertVisitedURL(normalizedLink)
					//if err != nil {
					//	log.Printf("Error inserting visited URL into the database: %v", err)
					//}
				}
			}
		}
	}()

	wg.Wait()
	close(ch)
}
