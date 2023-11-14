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

// ItemData represents a generic item with metadata
type ItemData struct {
	Domain string      `json:"domain"`
	Data   GenericData `json:"data"`
}

type GenericData struct {
	Title          string            `json:"title"`
	URL            string            `json:"url"`
	Description    string            `json:"description"`
	Price          string            `json:"price"`
	Location       string            `json:"location,omitempty"`
	Features       []string          `json:"features,omitempty"`
	Reviews        []Review          `json:"reviews,omitempty"`
	Images         []string          `json:"images,omitempty"`
	AdditionalInfo map[string]string `json:"additional_info,omitempty"`
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

func marshallDataToJson(data []ItemData) ([]byte, error) {
	wrappedData := map[string][]ItemData{
		"items": data,
	}

	jsonData, err := json.MarshalIndent(wrappedData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling data to JSON: %w", err)
	}

	return jsonData, nil
}

func writeJsonToFile(filename string, jsonData []byte) error {
	err := ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}
	return nil
}

func scrapeURLs(urls []string) []ItemData {
	allData := make([]ItemData, 0)

	for _, pageURL := range urls {
		c := colly.NewCollector()
		extensions.RandomUserAgent(c)

		c.OnHTML("article.product_pod", func(e *colly.HTMLElement) {
			bookURL := e.ChildAttr("h3 a", "href")
			bookURL = e.Request.AbsoluteURL(bookURL)

			currentItem := ItemData{
				Domain: "books",
				Data: GenericData{
					Title:       e.ChildText("h3 a"),
					URL:         bookURL,
					Description: e.ChildText("p.description"),
					Price:       e.ChildText("div p.price_color"),
					Metadata: Metadata{
						Source:    pageURL,
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}

			allData = append(allData, currentItem)
		})

		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			err := c.Visit(pageURL)
			if err == nil {
				break
			}
			log.Printf("Error visiting %s: %s, retrying (%d/%d)\n", pageURL, err, i+1, maxRetries)
			if i < maxRetries-1 {
				time.Sleep(time.Second * 10)
			}
		}

		time.Sleep(time.Second * 5)
	}

	jsonData, err := marshallDataToJson(allData)
	if err != nil {
		log.Println(err)
		return allData
	}

	err = writeJsonToFile("scrapedData.json", jsonData)
	if err != nil {
		log.Println(err)
	}

	log.Println("Scraping completed and data has been saved to scrapedData.json")
	return allData
}

func getURLsToScrape() ([]string, error) {
	var urls []string

	jsonData, err := ioutil.ReadFile("crawledUrls.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &urls)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func main() {
	urlsToScrape, err := getURLsToScrape()
	if err != nil {
		log.Fatalf("Failed to get URLs to scrape: %v", err)
	}
	scrapedData := scrapeURLs(urlsToScrape)
	fmt.Println("Scraped Data:", scrapedData)
}
