package main

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// TopLevelStruct represents the top-level structure of the JSON file
type TopLevelStruct struct {
	Items []ItemData `json:"items"`
}

type SiteMap map[string][]string

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
	rand.Seed(int64(uint64(time.Now().UnixNano())))
	index := rand.Intn(len(userAgents))
	return userAgents[index]
}

// Function to read site map from a JSON file
func readSiteMap(filePath string) (SiteMap, error) {
	var siteMap SiteMap
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fileContents, &siteMap); err != nil {
		return nil, err
	}
	return siteMap, nil
}

// Function to save data to a JSON file
func saveToJson(filePath string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, jsonData, 0644)
}

func getDomainFromURL(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		// Handle the error according to your needs
		log.Printf("Error parsing URL '%s': %v", urlStr, err)
		return ""
	}

	// Split the host by '.' and extract the domain part
	// This simple split by '.' assumes a basic domain structure like 'example.com'
	parts := strings.Split(parsedURL.Host, ".")
	if len(parts) > 1 {
		return parts[len(parts)-2] + "." + parts[len(parts)-1]
	}
	return parsedURL.Host
}

func scrapeURLs(urls []string) []ItemData {
	allData := make([]ItemData, 0)

	for _, pageURL := range urls {
		c := colly.NewCollector(
			colly.UserAgent(getRandomUserAgent()),
		)
		// Random user agent for each request
		extensions.RandomUserAgent(c)

		domain := getDomainFromURL(pageURL)

		c.OnHTML("article.product_pod", func(e *colly.HTMLElement) {
			bookURL := e.ChildAttr("h3 a", "href")
			bookURL = e.Request.AbsoluteURL(bookURL)

			currentItem := ItemData{
				Domain: domain, // Use the extracted domain
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
			if err := c.Visit(pageURL); err == nil {
				break
			}
			log.Printf("Error visiting %s: %s, retrying (%d/%d)\n", pageURL, i+1, maxRetries)
			if i < maxRetries-1 {
				time.Sleep(time.Second * 10)
			}
		}

		time.Sleep(time.Second * 5)
	}

	log.Println("Scraping completed and data has been saved to scrapedData.json")
	return allData
}

func getURLsToScrape() ([]string, error) {
	var siteMap SiteMap
	jsonData, err := ioutil.ReadFile("processedSiteMap.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonData, &siteMap)
	if err != nil {
		return nil, err
	}

	// Extract URLs from the siteMap
	var urls []string
	for url := range siteMap {
		urls = append(urls, url)
	}

	return urls, nil
}

func main() {
	// Load site map from file
	siteMap, err := readSiteMap("siteMap.json")
	if err != nil {
		log.Fatalf("Error reading site map: %v", err)
	}

	// Save the unmarshaled data to a new file for later access
	err = saveToJson("processedSiteMap.json", siteMap)
	if err != nil {
		log.Fatalf("Error saving processed site map: %v", err)
	}

	// Get URLs to scrape from the processed site map
	urlsToScrape, err := getURLsToScrape()
	if err != nil {
		log.Fatalf("Failed to get URLs to scrape: %v", err)
	}

	// Scrape URLs
	scrapedData := scrapeURLs(urlsToScrape)
	log.Printf("Scraped Data: %+v\n", scrapedData)

	// Save scraped data
	if err := saveToJson("scrapedData.json", scrapedData); err != nil {
		log.Fatalf("Error saving scraped data: %v", err)
	}

	log.Println("Scraping completed and data has been saved to scrapedData.json")
}
