package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"golang.org/x/exp/rand"
	"io/ioutil"
	"log"
	"net/http"
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

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
	"Mozilla/5.0 (iPad; CPU OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
	"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.58 Mobile Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.82 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:97.0) Gecko/20100101 Firefox/97.0",
	"Mozilla/5.0 (Windows NT 10.0; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
	"Opera/9.80 (Windows NT 6.0) Presto/2.12.388 Version/12.14",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/74.0",
}

func getRandomUserAgent() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	index := rand.Intn(len(userAgents))
	return userAgents[index]
}

func scrapeURLs(urls []string) []ItemData {

	allData := make([]ItemData, 0)

	client := &http.Client{}

	for _, url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Failed to create request for URL %s: %v", url, err)
			continue
		}

		req.Header.Set("User-Agent", getRandomUserAgent())

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to get URL %s: %v", url, err)
			continue
		}

		// Process the response as needed, for example, read the body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body for URL %s: %v", url, err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// For demonstration purposes, we're just printing the length of the body
		fmt.Printf("URL: %s, Response body length: %d\n", url, len(body))
		fmt.Println("Response body:", string(body))
	}

	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

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
