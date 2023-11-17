package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"io/ioutil"
	"sync"
	"time"
)

// ScraperConfig holds the configuration for the scraper
type ScraperConfig struct {
	StartingURLs []string
}

// NewScraperConfig creates a new ScraperConfig with default values
func NewScraperConfig(startingURLs []string) ScraperConfig {
	return ScraperConfig{
		StartingURLs: startingURLs,
	}
}

// Scrape performs the scraping based on the provided configuration
func Scrape(startingURL string, domain string, wg *sync.WaitGroup, resultChan chan []ItemData) {
	defer wg.Done()

	// Container for scraped data
	var allData []ItemData

	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	// Collect the data for each book on the first page of each URL
	c.OnHTML("article.product_pod", func(e *colly.HTMLElement) {
		bookURL := e.ChildAttr("h3 a", "href")
		// Assuming we have a function to resolve the relative bookURL to absolute
		bookURL = e.Request.AbsoluteURL(bookURL)

		currentItem := ItemData{
			Domain: domain,
			Data: GenericData{
				Title:       e.ChildText("h3 a"),
				URL:         bookURL,
				Description: e.ChildText("p.description"), // Selector assumed, replace with the actual selector
				Price:       e.ChildText("div p.price_color"),
				Metadata: Metadata{
					Source:    e.Request.URL.String(),
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}

		allData = append(allData, currentItem)
	})

	// Visit the URL with retry logic
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := c.Visit(startingURL)
		if err == nil {
			break // No error, break the retry loop
		}
		fmt.Printf("Error visiting %s: %s, retrying (%d/%d)\n", startingURL, err, i+1, maxRetries)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 10) // Wait before retrying
		}
	}

	// Sleep to prevent rate-limiting issues
	time.Sleep(time.Second * 5)

	resultChan <- allData
}

type Metadata struct {
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

type GenericData struct {
	Title          string            `json:"title"`
	URL            string            `json:"url"`
	Description    string            `json:"description"` // This could be the book synopsis if available
	Price          string            `json:"price"`
	Location       string            `json:"location,omitempty"`        // Omitted if not applicable
	Features       []string          `json:"features,omitempty"`        // Omitted if not applicable
	Reviews        []Review          `json:"reviews,omitempty"`         // Omitted if not applicable
	Images         []string          `json:"images,omitempty"`          // Omitted if not applicable
	AdditionalInfo map[string]string `json:"additional_info,omitempty"` // Flexible for any additional data
	Metadata       Metadata          `json:"metadata"`
}

type Review struct {
	User    string `json:"user"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type ItemData struct {
	Domain string      `json:"domain"`
	Data   GenericData `json:"data"`
}

func main() {
	startingURLs := []string{
		"http://books.toscrape.com/catalogue/category/books/fiction_10/index.html",
		"https://books.toscrape.com/catalogue/category/books/philosophy_7/index.html",
	}

	var wg sync.WaitGroup
	resultChan := make(chan []ItemData, len(startingURLs))

	// Launch a goroutine for each URL
	for _, url := range startingURLs {
		wg.Add(1)
		go Scrape(url, "dynamic_domain", &wg, resultChan) // Pass the dynamic domain here
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from channels
	var allData []ItemData
	for result := range resultChan {
		allData = append(allData, result...)
	}

	// Wrap the data
	wrappedData := map[string][]ItemData{
		"items": allData,
	}

	// Marshal the wrapped data into JSON
	jsonData, err := json.MarshalIndent(wrappedData, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling data to JSON:", err)
		return
	}

	// Write the JSON data to a file
	err = ioutil.WriteFile("scrapedData.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
	}

	fmt.Println("Scraping completed and data has been saved to scrapedData.json")
}

//package main
//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/gocolly/colly"
//	"github.com/gocolly/colly/extensions"
//	"io/ioutil"
//	"sync"
//	"time"
//)
//
//// ScraperConfig holds the configuration for the scraper
//type ScraperConfig struct {
//	StartingURLs []string
//}
//
//// NewScraperConfig creates a new ScraperConfig with default values
//func NewScraperConfig(startingURLs []string) ScraperConfig {
//	return ScraperConfig{
//		StartingURLs: startingURLs,
//	}
//}
//
//// Scrape performs the scraping based on the provided configuration
//func Scrape(startingURL string, domain string, wg *sync.WaitGroup, resultChan chan []ItemData) {
//	defer wg.Done()
//
//	// Container for scraped data
//	var allData []ItemData
//
//	c := colly.NewCollector()
//	extensions.RandomUserAgent(c)
//
//	// Collect the data for each book on the first page of each URL
//	c.OnHTML("article.product_pod", func(e *colly.HTMLElement) {
//		bookURL := e.ChildAttr("h3 a", "href")
//		// Assuming we have a function to resolve the relative bookURL to absolute
//		bookURL = e.Request.AbsoluteURL(bookURL)
//
//		currentItem := ItemData{
//			Domain: domain,
//			Data: GenericData{
//				Title:       e.ChildText("h3 a"),
//				URL:         bookURL,
//				Description: e.ChildText("p.description"), // Selector assumed, replace with the actual selector
//				Price:       e.ChildText("div p.price_color"),
//				Metadata: Metadata{
//					Source:    e.Request.URL.String(),
//					Timestamp: time.Now().Format(time.RFC3339),
//				},
//			},
//		}
//
//		allData = append(allData, currentItem)
//	})
//
//	// Visit the URL with retry logic
//	maxRetries := 3
//	for i := 0; i < maxRetries; i++ {
//		err := c.Visit(startingURL)
//		if err == nil {
//			break // No error, break the retry loop
//		}
//		fmt.Printf("Error visiting %s: %s, retrying (%d/%d)\n", startingURL, err, i+1, maxRetries)
//		if i < maxRetries-1 {
//			time.Sleep(time.Second * 10) // Wait before retrying
//		}
//	}
//
//	// Sleep to prevent rate-limiting issues
//	time.Sleep(time.Second * 5)
//
//	resultChan <- allData
//}
//
//type Metadata struct {
//	Source    string `json:"source"`
//	Timestamp string `json:"timestamp"`
//}
//
//type GenericData struct {
//	Title          string            `json:"title"`
//	URL            string            `json:"url"`
//	Description    string            `json:"description"` // This could be the book synopsis if available
//	Price          string            `json:"price"`
//	Location       string            `json:"location,omitempty"`        // Omitted if not applicable
//	Features       []string          `json:"features,omitempty"`        // Omitted if not applicable
//	Reviews        []Review          `json:"reviews,omitempty"`         // Omitted if not applicable
//	Images         []string          `json:"images,omitempty"`          // Omitted if not applicable
//	AdditionalInfo map[string]string `json:"additional_info,omitempty"` // Flexible for any additional data
//	Metadata       Metadata          `json:"metadata"`
//}
//
//type Review struct {
//	User    string `json:"user"`
//	Rating  int    `json:"rating"`
//	Comment string `json:"comment"`
//}
//
//type ItemData struct {
//	Domain string      `json:"domain"`
//	Data   GenericData `json:"data"`
//}
//
//func main() {
//	startingURLs := []string{
//		"http://books.toscrape.com/catalogue/category/books/fiction_10/index.html",
//		"https://books.toscrape.com/catalogue/category/books/philosophy_7/index.html",
//	}
//
//	var wg sync.WaitGroup
//	resultChan := make(chan []ItemData, len(startingURLs))
//
//	// Launch a goroutine for each URL
//	for _, url := range startingURLs {
//		wg.Add(1)
//		go Scrape(url, "dynamic_domain", &wg, resultChan) // Pass the dynamic domain here
//	}
//
//	// Wait for all goroutines to finish
//	go func() {
//		wg.Wait()
//		close(resultChan)
//	}()
//
//	// Collect results from channels
//	var allData []ItemData
//	for result := range resultChan {
//		allData = append(allData, result...)
//	}
//
//	// Wrap the data
//	wrappedData := map[string][]ItemData{
//		"items": allData,
//	}
//
//	// Marshal the wrapped data into JSON
//	jsonData, err := json.MarshalIndent(wrappedData, "", "  ")
//	if err != nil {
//		fmt.Println("Error marshalling data to JSON:", err)
//		return
//	}
//
//	// Write the JSON data to a file
//	err = ioutil.WriteFile("scrapedData.json", jsonData, 0644)
//	if err != nil {
//		fmt.Println("Error writing JSON to file:", err)
//	}
//
//	fmt.Println("Scraping completed and data has been saved to scrapedData.json")
//}
