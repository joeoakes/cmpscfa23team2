package CRAB

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// ScrapeResult stores the result of a scrape
type ScrapeResult struct {
	URL   string
	Data  string // You can replace this with a more detailed structure
	Error error
}

// scrape performs web scraping on a given url
func scrape(url string, ch chan<- ScrapeResult) {
	res, err := http.Get(url)
	if err != nil {
		ch <- ScrapeResult{URL: url, Error: err}
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		ch <- ScrapeResult{URL: url, Error: fmt.Errorf("status code: %d", res.StatusCode)}
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		ch <- ScrapeResult{URL: url, Error: err}
		return
	}

	// This is a basic example. You'd typically extract more detailed information here.
	content, _ := doc.Find("body").Html()
	ch <- ScrapeResult{URL: url, Data: content}
}

func main() {
	// Sample list of URLs to scrape. You'd replace this with your actual list.
	urls := []string{
		"http://example.com",
		"http://example.org",
		"http://example.net",
	}

	// Channel to collect scrape results
	ch := make(chan ScrapeResult)

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			scrape(u, ch)
		}(url)
	}

	// This goroutine collects results and prints them.
	go func() {
		for result := range ch {
			if result.Error != nil {
				log.Printf("Error scraping %s: %v", result.URL, result.Error)
			} else {
				// Just printing the beginning of the content for brevity.
				// Replace this with your data processing logic.
				fmt.Printf("Scraped %s: %s...\n", result.URL, result.Data[:100])
			}
		}
	}()

	// Wait for all scraping goroutines to complete
	wg.Wait()
	close(ch)
}
