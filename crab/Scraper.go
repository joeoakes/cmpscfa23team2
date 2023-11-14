package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"io/ioutil"
	"log"
	"time"
)

// TopLevelStruct represents the top-level structure of the JSON file
type TopLevelStruct struct {
	Items []ItemData `json:"items"`
}

// ItemData represents a generic item with metadata
type ItemData struct {
	Domain string      `json:"domain"`
	Data   GenericData `json:"data"`
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

type Metadata struct {
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

func main() {
	// Container for all scraped data
	allData := make([]ItemData, 0)

	startingURL := "http://books.toscrape.com/catalogue/category/books/fiction_10/index.html"
	fmt.Println("Constructed URL:", startingURL)

	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	// Collect the data for each book on the first page of the fiction section
	c.OnHTML("article.product_pod", func(e *colly.HTMLElement) {
		bookURL := e.ChildAttr("h3 a", "href")
		// Assuming we have a function to resolve the relative bookURL to absolute
		bookURL = e.Request.AbsoluteURL(bookURL)

		currentItem := ItemData{
			Domain: "books",
			Data: GenericData{
				Title:       e.ChildText("h3 a"),
				URL:         bookURL,
				Description: e.ChildText("p.description"), // Selector assumed, replace with the actual selector
				Price:       e.ChildText("div p.price_color"),
				Metadata: Metadata{
					Source:    startingURL,
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

	// Read the JSON file
	fileContents, err := ioutil.ReadFile("scrapedData.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}

	// Unmarshal JSON data into struct
	var data TopLevelStruct
	err = json.Unmarshal(fileContents, &data)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON data: %s", err)
	}

	// Now you can access the data from the JSON file
	// For example, printing the titles of each item
	for _, item := range data.Items {
		fmt.Println("Title:", item.Data.Title)
	}

}
