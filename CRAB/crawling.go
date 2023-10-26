/*

First:
go get github.com/PuerkitoBio/goquery
go get github.com/temoto/robotstxt

Explanation:
	We first fetch and parse the robots.txt file to make sure we're allowed to scrape our target URLs.
	Rate limiting is implemented with a simple sleep. Depending on your needs and the site's policies, you might want to refine this.
	The scrape function is tailored to extract house prices. The specific HTML element and,
 	class (.house-price-class) would need to match the actual structure of the website you're targeting.
	Always remember that web scraping may be subject to legal and ethical restrictions. 
 	Always ensure you have permission to scrape a website, and always respect its robots.txt and terms of use.
*/

package CRAB

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/temoto/robotstxt-go"
)

type ScrapeResult struct {
	URL    string
	Price  string
	Error  error
}

func fetchRobotsTxt(url string) (*robotstxt.Group, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		return nil, err
	}

	return data.FindGroup("*"), nil
}

func scrape(url string, ch chan<- ScrapeResult, userAgent string, delay time.Duration) {
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
	ch <- ScrapeResult{URL: url, Price: price}
}

func main() {
	baseURL := "http://example-estate-site.com"
	robotsURL := baseURL + "/robots.txt"
	userAgent := "YourBotName" // replace with your bot's name
	urls := []string{ // Some sample house URLs
		baseURL + "/house1",
		baseURL + "/house2",
		// ... add as many URLs as needed
	}

	// Fetch robots.txt
	robotsGroup, err := fetchRobotsTxt(robotsURL)
	if err != nil {
		log.Fatalf("Failed fetching robots.txt: %v", err)
	}

	// Ensure our bot can access the desired URLs
	for _, url := range urls {
		if !robotsGroup.Test(url) {
			log.Fatalf("Access denied by robots.txt for URL: %s", url)
		}
	}

	ch := make(chan ScrapeResult)
	var wg sync.WaitGroup

	delay := 1 * time.Second // 1 request per second. Adjust as needed.
	for _, url := range urls {
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
			}
		}
	}()

	wg.Wait()
	close(ch)
}
