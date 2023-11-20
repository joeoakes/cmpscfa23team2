package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

// GetRandomUserAgent is accessible because it starts with a capital letter
func GetRandomUserAgent() string {
	userAgents := []string{
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
		"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_16) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/536.6",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:90.0) Gecko/20100101 Firefox/90.0",
		"Mozilla/5.0 (Android 11; Mobile; LG-M255; rv:90.0) Gecko/90.0 Firefox/90.0",
		"Mozilla/5.0 (iPad; CPU OS 14_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:91.0) Gecko/20100101 Firefox/91.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/90.0.4430.212 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Trident/7.0; Touch; rv:11.0) like Gecko",
		"Mozilla/5.0 (X11; CrOS x86_64 13729.56.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.95 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.2 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64; rv:92.0) Gecko/20100101 Firefox/92.0",
		"Mozilla/5.0 (Android 10; Tablet; rv:68.0) Gecko/68.0 Firefox/68.0",
		"Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.17",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
	}
	rand.Seed(int64(uint64(time.Now().UnixNano())))
	index := rand.Intn(len(userAgents))
	return userAgents[index]
}

// ScraperConfig holds the configuration for the scraper
type ScraperConfig struct {
	StartingURLs []string
}

type DomainConfig struct {
	Name                            string
	ItemSelector                    string
	TitleSelector                   string
	URLSelector                     string
	DescriptionSelector             string
	PriceSelector                   string
	FactorsSelector                 string
	DepreciationRatesSelector       string
	ModelsLeastDepreciationSelector string
	ModelsMostDepreciationSelector  string
}

var domainConfigurations = map[string]DomainConfig{
	"airfare": {
		Name:                "airfare",
		ItemSelector:        "div.article-content",      // Adjust this selector based on the structure of the page
		TitleSelector:       "h1",                       // Adjust this selector based on the structure of the page
		URLSelector:         "meta[property='og:url']",  // Adjust this selector based on the structure of the page
		DescriptionSelector: "meta[name='description']", // Adjust this selector based on the structure of the page
		PriceSelector:       "span.airfare-price",       // Adjust this selector based on the structure of the page
	},
	"books": {
		Name:                "books",
		ItemSelector:        "article.product_pod",
		TitleSelector:       "h3 a",
		URLSelector:         "h3 a",
		DescriptionSelector: "p.description", // Selector assumed, replace with the actual selector
		PriceSelector:       "div p.price_color",
	},
	"job-market": {
		Name:                "job-market",
		ItemSelector:        "div.job-posting",
		TitleSelector:       "h2.job-title",
		URLSelector:         "a.job-apply-link",
		DescriptionSelector: "div.job-description",
		PriceSelector:       "", // In job market, you might have salary rather than price
	},
	"nascar-predictem": {
		Name:                "nascar-predictem",
		ItemSelector:        "article",               // Selector for the main article
		TitleSelector:       "h1.entry-title",        // Selector for the title
		URLSelector:         "link[rel='canonical']", // Selector for the URL
		DescriptionSelector: "p",                     // Selector for the article paragraphs
		PriceSelector:       "",                      // Price selector might not be applicable here
	},
	"car-depreciation": {
		Name:                            "car-depreciation",
		ItemSelector:                    "div#content",
		TitleSelector:                   "h1.entry-title",
		URLSelector:                     "meta[property='og:url']",
		DescriptionSelector:             "div#content p:first-of-type",
		FactorsSelector:                 "div#content h2:contains('What causes a vehicle to depreciate?') + p, div#content h2:contains('What causes a vehicle to depreciate?') + ul > li",
		DepreciationRatesSelector:       "div#content h2:contains('How much does a car depreciate per year?') + p",
		ModelsLeastDepreciationSelector: "div#content h2:contains('Top 10 cars that depreciate the least') + p, div#content h2:contains('Top 10 cars that depreciate the least') + ul > li",
		ModelsMostDepreciationSelector:  "div#content h2:contains('Top 10 cars that depreciate the most') + p, div#content h2:contains('Top 10 cars that depreciate the most') + ul > li",
	},
}

// NewScraperConfig creates a new ScraperConfig with default values
func NewScraperConfig(startingURLs []string) ScraperConfig {
	return ScraperConfig{
		StartingURLs: startingURLs,
	}
}

func insertData(data ItemData, filename string) error {
	// Save data to JSON file
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Scrape performs the scraping based on the provided configuration
func Scrape(startingURL string, domainConfig DomainConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	c := colly.NewCollector(
		colly.UserAgent(GetRandomUserAgent()),
	)

	// Container for scraped data
	var allData []GenericData

	// Define scraping logic based on the domain
	switch domainConfig.Name {
	case "car-depreciation":
		c.OnHTML(domainConfig.ItemSelector, func(e *colly.HTMLElement) {
			var factors, modelsLeast, modelsMost []string
			var depreciationRates = make(map[string]string)

			// Extracting factors
			e.ForEach(domainConfig.FactorsSelector, func(_ int, el *colly.HTMLElement) {
				factors = append(factors, el.Text)
			})

			// Extracting models with least and most depreciation
			e.ForEach(domainConfig.ModelsLeastDepreciationSelector, func(_ int, el *colly.HTMLElement) {
				modelsLeast = append(modelsLeast, el.Text)
			})
			e.ForEach(domainConfig.ModelsMostDepreciationSelector, func(_ int, el *colly.HTMLElement) {
				modelsMost = append(modelsMost, el.Text)
			})

			// Handling depreciation rates
			e.ForEach(domainConfig.DepreciationRatesSelector, func(_ int, el *colly.HTMLElement) {
				// Here you might need to parse the text to extract meaningful data
				// This is an example, adjust it as per the actual text format
				lines := strings.Split(el.Text, "\n")
				for _, line := range lines {
					parts := strings.SplitN(line, ": ", 2)
					if len(parts) == 2 {
						depreciationRates[parts[0]] = parts[1]
					}
				}
			})

			// Constructing the item
			currentItem := GenericData{
				Title:                   e.ChildText(domainConfig.TitleSelector),
				URL:                     e.Request.AbsoluteURL(e.ChildAttr(domainConfig.URLSelector, "href")),
				Description:             e.ChildText(domainConfig.DescriptionSelector),
				Price:                   "", // No specific price data for this domain
				Factors:                 factors,
				DepreciationRates:       depreciationRates,
				ModelsLeastDepreciation: modelsLeast,
				ModelsMostDepreciation:  modelsMost,
				Metadata: Metadata{
					Source:    e.Request.URL.String(),
					Timestamp: time.Now().Format(time.RFC3339),
				},
			}

			allData = append(allData, currentItem)
		})
	case "airfare", "books", "job-market", "nascar-predictem":
		// General scraping logic for other domains
		c.OnHTML(domainConfig.ItemSelector, func(e *colly.HTMLElement) {
			itemURL := e.Request.AbsoluteURL(e.ChildAttr(domainConfig.URLSelector, "href"))

			currentItem := GenericData{
				Title:       e.ChildText(domainConfig.TitleSelector),
				URL:         itemURL,
				Description: e.ChildText(domainConfig.DescriptionSelector),
				Price:       e.ChildText(domainConfig.PriceSelector),
				Metadata: Metadata{
					Source:    e.Request.URL.String(),
					Timestamp: time.Now().Format(time.RFC3339),
				},
			}
			allData = append(allData, currentItem)
		})
	}

	// Visit the URL with retry logic
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := c.Visit(startingURL)
		if err == nil {
			break
		}
		fmt.Printf("Error visiting %s: %s, retrying (%d/%d)\n", startingURL, err, i+1, maxRetries)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 10)
		}
	}

	// Save data to JSON file
	filename := fmt.Sprintf("%s_data.json", domainConfig.Name)
	err := insertData(ItemData{
		Domain: domainConfig.Name,
		Data:   allData,
	}, filename)
	if err != nil {
		fmt.Printf("Error saving data to JSON file: %v\n", err)
	}
}
func airfair() {
	content, err := ioutil.ReadFile("data.html")
	if err != nil {
		log.Fatal(err)
	}
	html := string(content)

	var data []YearData

	doc, err := goquery.NewDocumentFromReader(strings.NewReader())
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("table tbody tr").Each(func(rowIndex int, rowHtml *goquery.Selection) {
		if rowIndex == 0 { // Skip the header row
			return
		}

		var yearData YearData
		rowHtml.Find("td").Each(func(cellIndex int, cellHtml *goquery.Selection) {
			switch cellIndex {
			case 0:
				yearData.Year = cellHtml.Text()
			case 1:
				yearData.Jan = cellHtml.Text()
			case 2:
				yearData.Feb = cellHtml.Text()
			case 3:
				yearData.Mar = cellHtml.Text()
			case 4:
				yearData.Apr = cellHtml.Text()
			case 5:
				yearData.May = cellHtml.Text()
			case 6:
				yearData.Jun = cellHtml.Text()
			case 7:
				yearData.July = cellHtml.Text()
			case 8:
				yearData.Aug = cellHtml.Text()
			case 9:
				yearData.Sept = cellHtml.Text()
			case 10:
				yearData.Oct = cellHtml.Text()
			case 11:
				yearData.Nov = cellHtml.Text()
			case 12:
				yearData.Dec = cellHtml.Text()
			case 13:
				yearData.Avg = cellHtml.Text()

			}
		})

		data = append(data, yearData)
	})

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
}
func testScrape(domainName string) {
	domainConfig, exists := domainConfigurations[domainName]
	if !exists {
		fmt.Printf("Invalid domain name provided: %s\n", domainName)
		return
	}

	// Test URLs for the specified domain
	testURLs := map[string][]string{
		"airfare":          {"https://www.usinflationcalculator.com/inflation/airfare-inflation/"},
		"books":            {"http://books.toscrape.com/catalogue/category/books/fiction_10/index.html"},
		"job-market":       {"https://www.example.com/job-market"},
		"nascar-predictem": {"https://www.predictem.com/nascar/xfinity-500-race-preview-picks/"},
		"car-depreciation": {"https://www.thinkinsure.ca/insurance-help-centre/car-deprecation.html"},

	}

	startingURLs := testURLs[domainName]
	var wg sync.WaitGroup

	// Launch a goroutine for each URL
	for _, url := range startingURLs {
		wg.Add(1)
		go Scrape(url, domainConfig, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Printf("Scraping for domain %s completed and data has been saved to JSON files\n", domainName)
}

type Metadata struct {
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

type GenericData struct {
	Title                   string            `json:"title"`
	URL                     string            `json:"url"`
	Description             string            `json:"description"`
	Price                   string            `json:"price"`
	Factors                 []string          `json:"factors"`
	DepreciationRates       map[string]string `json:"depreciation_rates"`
	ModelsLeastDepreciation []string          `json:"models_least_depreciation"`
	ModelsMostDepreciation  []string          `json:"models_most_depreciation"`
	Metadata                Metadata          `json:"metadata"`
}

type ItemData struct {
	Domain string        `json:"domain"`
	Data   []GenericData `json:"data"`
}

func main() {
	//GetRandomUserAgent()
	// Display available domain options to the user
	fmt.Println("Available domains:")
	for domainName := range domainConfigurations {
		fmt.Printf("- %s\n", domainName)
	}

	// Ask the user to choose a domain
	var domainName string
	fmt.Print("Enter the domain you want to scrape: ")
	fmt.Scanln(&domainName)

	// Check if the chosen domain is valid
	_, exists := domainConfigurations[domainName]
	if !exists {
		fmt.Printf("Invalid domain name provided: %s\n", domainName)
		return
	}

	// Perform the scraping for the chosen domain
	testScrape(domainName)

}
