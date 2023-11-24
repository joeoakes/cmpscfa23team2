package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/temoto/robotstxt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// begin structs ========================================================================================================
// Data Structs
type URLData struct {
	URL     string    // The URL to be crawled
	Created time.Time // Timestamp of URL creation or retrieval
	Links   []string  // URLs found on this page
}

type MonthData struct {
	Month string `json:"month"`
	Rate  string `json:"rate"`
}

type AirfareData struct {
	Domain string `json:"domain"`
	URL    string `json:"url"`
	Data   struct {
		Title          string   `json:"title"`
		Year           string   `json:"year"`
		Location       string   `json:"location"`
		Features       []string `json:"features"`
		AdditionalInfo struct {
			Country    string      `json:"country"`
			MonthsData []MonthData `json:"months_data"`
		} `json:"additional_info"`
		Metadata struct {
			Source    string `json:"source"`
			Timestamp string `json:"timestamp"`
		} `json:"metadata"`
	} `json:"data"`
}

type YearData struct {
	Year string `json:"year"`
	Jan  string `json:"jan"`
	Feb  string `json:"feb"`
	Mar  string `json:"mar"`
	Apr  string `json:"apr"`
	May  string `json:"may"`
	Jun  string `json:"jun"`
	July string `json:"july"`
	Aug  string `json:"aug"`
	Sept string `json:"sept"`
	Oct  string `json:"oct"`
	Nov  string `json:"nov"`
	Dec  string `json:"dec"`
	Avg  string `json:"avg"`
}

type GasolineData struct {
	Year                     string `json:"year"`
	AverageGasolinePrices    string `json:"average_gasoline_prices"`
	AverageAnnualCPIForGas   string `json:"average_annual_cpi_for_gas"`
	GasPricesAdjustedForInfl string `json:"gas_prices_adjusted_for_inflation"`
}

type PropertyData struct {
	Status    string `json:"status"`
	Bedrooms  string `json:"bedrooms"`
	Bathrooms string `json:"bathrooms"`
	AcreLot   string `json:"acre_lot"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	HouseSize string `json:"house_size"`
	SoldDate  string `json:"prev_sold_date"`
	Price     string `json:"price"`
}

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

//end structs ==========================================================================================================

//begin user agents ====================================================================================================

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

//end user agents ======================================================================================================

// begin crawler vars =========================================================================================================
var (
	urlQueue = make(chan string, 1000)
	visited  = make(map[string]bool)
	//dataChan = make(chan URLData, 1000)//uncomment this use case appears
)

//end crawler vars =====================================================================================================

//begin domain configurations ==========================================================================================

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

//end domain configurations ============================================================================================

//currently unused =====================================================================================================
//
//// NewScraperConfig creates a new ScraperConfig with default values
//func NewScraperConfig(startingURLs []string) ScraperConfig {
//	return ScraperConfig{
//		StartingURLs: startingURLs,
//	}
//}
//======================================================================================================================

//begin crawler ========================================================================================================

// begin intialize crawling =============================================================================================
//
//	InitializeCrawling sets up and starts the crawling process.
func InitializeCrawling() {
	log.Println("Fetching URLs to crawl...")
	urlDataList := getURLsToCrawl()
	log.Println("URLs to crawl:", urlDataList)
	threadedCrawl(urlDataList, 10)
}

//end intialize crawling ===============================================================================================

// begin crawl url ======================================================================================================
// getURLsToCrawl retrieves a list of URLs to be crawled.
func getURLsToCrawl() []URLData {
	return []URLData{
		{URL: "https://www.kaggle.com/search?q=housing+prices"},
		{URL: "http://books.toscrape.com/"},
		{URL: "https://www.kaggle.com/search?q=stocks"},
		{URL: "https://www.kaggle.com/search?q=stock+market"},
		{URL: "https://www.kaggle.com/search?q=real+estate"},
	}
}

//end crawl url ========================================================================================================

// begin insert data ====================================================================================================
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

//end insert data ======================================================================================================

// begin crawlurl ======================================================================================================
// crawlURL is responsible for crawling a single URL.
func crawlURL(urlData URLData, ch chan<- URLData, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure the WaitGroup counter is decremented on function exit
	c := colly.NewCollector(
		colly.UserAgent(GetRandomUserAgent()), // Set a random user agent
	)
	// First, check if the URL is allowed by robots.txt rules
	allowed := isURLAllowedByRobotsTXT(urlData.URL)
	if !allowed {
		return // Skip crawling if not allowed
	}

	// Handler for errors during the crawl
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error occurred while crawling %s: %s\n", urlData.URL, err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if !visited[link] && isURLAllowedByRobotsTXT(link) {
			urlData.Links = append(urlData.Links, link)
			urlQueue <- link
		}
	})

	// Handler for successful HTTP responses
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			// Successful crawl, process the response here
			ch <- urlData // Send the URLData to the channel
			fmt.Printf("Crawled URL: %s\n", urlData.URL)
		} else {
			// Handle cases where the status code is not 200
			fmt.Printf("Non-200 status code while crawling %s: %d\n", urlData.URL, r.StatusCode)
		}
	})

	// Start the crawl
	c.Visit(urlData.URL)

	ch <- urlData
}

//end crawlurl ========================================================================================================

// begin create sitemap =================================================================================================
func createSiteMap(urls []URLData) error {
	siteMap := make(map[string][]string)
	for _, u := range urls {
		siteMap[u.URL] = u.Links
	}

	jsonData, err := json.Marshal(siteMap)
	err = ioutil.WriteFile("siteMap.json", jsonData, 0644)
	if err != nil {
		log.Printf("Error writing sitemap to file: %v\n", err)
		return err
	}

	log.Println("Sitemap created successfully.")
	return nil
}

//end create sitemap ===================================================================================================

// begin robot.txt ======================================================================================================
// isURLAllowedByRobotsTXT checks if the given URL is allowed by the site's robots.txt.
func isURLAllowedByRobotsTXT(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Println("Error parsing URL:", err)
		return false
	}

	if parsedURL.Host == "" {
		log.Println("Invalid URL, no host found:", urlStr)
		return false
	}

	robotsURL := "http://" + parsedURL.Host + "/robots.txt"

	resp, err := http.Get(robotsURL)
	if err != nil {
		log.Println("Error fetching robots.txt:", err)
		return true
	}

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		log.Println("Error parsing robots.txt:", err)
		return true
	}

	return data.TestAgent(urlStr, "GoEngine")
}

//end robot.txt ========================================================================================================

// begin threaded crawl =================================================================================================
// threadedCrawl starts crawling the provided URLs concurrently.
func threadedCrawl(urls []URLData, concurrentCrawlers int) {
	var wg sync.WaitGroup
	ch := make(chan URLData, len(urls))

	rateLimitRule := &colly.LimitRule{
		DomainGlob:  "*",             // Apply to all domains
		Delay:       5 * time.Second, // Wait 5 seconds between requests
		RandomDelay: 5 * time.Second, // Add up to 5 seconds of random delay
	}

	log.Println("Starting crawling...")
	for _, urlData := range urls {
		wg.Add(1)

		go func(u URLData) {
			c := colly.NewCollector(
				colly.UserAgent(GetRandomUserAgent()),
			)
			c.Limit(rateLimitRule) // Set the rate limit rule

			crawlURL(u, ch, &wg)
		}(urlData)

		log.Println("Crawling URL:", urlData.URL)
		if len(urls) >= concurrentCrawlers {
			break
		}
	}

	log.Println("Waiting for crawlers to finish...")
	go func() {
		wg.Wait()
		close(ch)
		log.Println("All goroutines finished, channel closed.")
	}()

	var crawledURLs []URLData
	for urlData := range ch {
		crawledURLs = append(crawledURLs, urlData)
	}
	if err := createSiteMap(crawledURLs); err != nil {
		log.Println("Error creating sitemap:", err)
	}
}

//end threaded crawl ===================================================================================================

//end crawler ==========================================================================================================

//begin scraper ========================================================================================================

// begin scrape =========================================================================================================
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

//end scrape ===========================================================================================================

// begin test scrape ====================================================================================================
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

//end test scrape ======================================================================================================

// new scrapers =========================================================================================================
// begin airfare scraper =================================================================================================
func airdatatest() {
	scrapeurl := "https://www.usinflationcalculator.com/inflation/airfare-inflation/"
	res, err := http.Get(scrapeurl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	const switchYear = "2023" // Replace with the actual year
	const switchMonth = "Dec" // Replace with the actual month

	var months = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	var isSecondTable = false
	var file *os.File

	// Open the first file
	file, err = os.OpenFile("airfare_data_inflation.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open first JSON file: %s", err)
	}

	doc.Find("table tbody tr").Each(func(rowIndex int, rowHtml *goquery.Selection) {
		if rowIndex == 0 {
			return
		}

		var airfareData AirfareData
		airfareData.Domain = "airfare"
		airfareData.URL = scrapeurl
		airfareData.Data.Title = "Airfare Inflation Data"
		airfareData.Data.Location = "United States"
		airfareData.Data.Features = []string{"Month", "Inflation Rate"}
		airfareData.Data.AdditionalInfo.Country = "USA"
		airfareData.Data.Metadata.Source = scrapeurl
		airfareData.Data.Metadata.Timestamp = time.Now().Format(time.RFC3339)
		airfareData.Data.AdditionalInfo.MonthsData = make([]MonthData, 0)

		rowHtml.Find("td").Each(func(cellIndex int, cellHtml *goquery.Selection) {
			cellText := cellHtml.Text()
			if cellIndex == 0 {
				airfareData.Data.Year = cellText
			} else if cellIndex >= 1 && cellIndex <= 12 {
				monthData := MonthData{
					Month: months[cellIndex-1], // Correct usage of months array
					Rate:  cellText,
				}
				airfareData.Data.AdditionalInfo.MonthsData = append(airfareData.Data.AdditionalInfo.MonthsData, monthData)
			}
		})

		// Switching logic for files
		if airfareData.Data.Year == switchYear && !isSecondTable {
			for _, monthData := range airfareData.Data.AdditionalInfo.MonthsData {
				if monthData.Month == switchMonth {
					file.Close()
					file, err = os.OpenFile("airfare_data_price.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
					if err != nil {
						log.Fatalf("Failed to open second JSON file: %s", err)
					}
					isSecondTable = true
					break
				}
			}
		}

		jsonData, err := json.MarshalIndent(airfareData, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		if _, err := file.Write(jsonData); err != nil {
			log.Fatalf("Failed to write JSON data to file: %s", err)
		}
		if rowIndex < doc.Find("table tbody tr").Length()-1 {
			file.WriteString(",\n")
		} else {
			file.WriteString("\n")
		}
	})

	file.Close()
	log.Println("Airfare data written to respective files")
}

//end airfare scraper ==================================================================================================

// begin inflation scraper ==============================================================================================
func scrapeInflationData() {
	scrapeurl := "https://www.usinflationcalculator.com/inflation/current-inflation-rates/"
	res, err := http.Get(scrapeurl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data []YearData
	doc.Find("table tbody tr").Each(func(rowIndex int, rowHtml *goquery.Selection) {
		if rowIndex == 0 { // Skip the header row
			return
		}

		var yearData YearData
		rowHtml.Find("td").Each(func(cellIndex int, cellHtml *goquery.Selection) {
			text := cellHtml.Text()
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
				yearData.Avg = text
			}
		})
		data = append(data, yearData)
	})

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("inflation_data.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON data to file: %s", err)
	}

	fmt.Println("Inflation data written to inflation_data.json")
}

//end inflation scraper ================================================================================================

// begin gasoline scraper =================================================================================================
func scrapeGasInflationData() {
	scrapeurl := "https://www.usinflationcalculator.com/gasoline-prices-adjusted-for-inflation/"
	res, err := http.Get(scrapeurl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data []GasolineData
	doc.Find("table tbody tr").Each(func(rowIndex int, rowHtml *goquery.Selection) {
		if rowIndex == 0 { // Skip the header row
			return
		}

		var gasData GasolineData
		rowHtml.Find("td").Each(func(cellIndex int, cellHtml *goquery.Selection) {
			text := cellHtml.Text()
			switch cellIndex {
			case 0:
				gasData.Year = text
			case 1:
				gasData.AverageGasolinePrices = text
			case 2:
				gasData.AverageAnnualCPIForGas = text
			case 3:
				gasData.GasPricesAdjustedForInfl = text
			}
		})
		data = append(data, gasData)
	})

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("gasoline_data.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON data to file: %s", err)
	}

	fmt.Println("Gasoline data written to gasoline_data.json")
}

//end gasoline scraper =================================================================================================

// begin housing scraper =================================================================================================
func scrapeHousingData() {
	scrapeurl := "https://www.kaggle.com/datasets/ahmedshahriarsakib/usa-real-estate-dataset"
	res, err := http.Get(scrapeurl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var properties []PropertyData
	doc.Find(".sc-fLdTid.sc-eZkIzG.iXbLwD.cefCfQ").Each(func(i int, s *goquery.Selection) {
		var data PropertyData
		s.Find("div").Each(func(index int, item *goquery.Selection) {
			switch index {
			case 0:
				data.Status = item.Text()
			case 1:
				data.Bedrooms = item.Text()
			case 2:
				data.Bathrooms = item.Text()
			case 3:
				data.AcreLot = item.Text()
			case 4:
				data.City = item.Text()
			case 5:
				data.State = item.Text()
			case 6:
				data.ZipCode = item.Text()
			case 7:
				data.HouseSize = item.Text()
			case 8:
				data.SoldDate = item.Text()
			case 9:
				data.Price = item.Text()
			}
		})
		properties = append(properties, data)
	})

	jsonData, err := json.MarshalIndent(properties, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("property_data.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON data to file: %s", err)
	}

	fmt.Println("Property data written to property_data.json")
}

//end housing scraper ===================================================================================================

//end scrapers ==========================================================================================================

// begin readcsv =========================================================================================================
func readCSV(filePath string) ([]PropertyData, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	log.Println("Successfully opened CSV file.")
	// Create a CSV reader from the file
	reader := csv.NewReader(file)
	reader.Comma = ',' // Set the delimiter to comma
	reader.TrimLeadingSpace = true
	log.Println("Reading CSV file...")
	// Read all the records at once
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return []PropertyData{}, nil
	}
	log.Println("Successfully read", len(records), "records from CSV file.")
	// Process records after the header row
	properties := make([]PropertyData, 0, len(records)-1)
	for i, record := range records {
		if i == 0 {
			continue // Skip the header row
		}

		if len(record) != 10 {
			return nil, fmt.Errorf("unexpected number of fields at row %d: got %d, want 10", i+1, len(record))
		}

		// Create a PropertyData struct for each record
		property := PropertyData{
			Status:    record[0],
			Bedrooms:  record[1],
			Bathrooms: record[2],
			AcreLot:   record[3],
			City:      record[4],
			State:     record[5],
			ZipCode:   record[6],
			HouseSize: record[7],
			SoldDate:  record[8],
			Price:     record[9],
		}
		properties = append(properties, property)
	}
	log.Println("Successfully read", len(properties), "properties from CSV file.")
	return properties, nil
}

//end readcsv ===========================================================================================================

// begin main ===========================================================================================================
func main() {

	//begin crawler

	airdatatest()
	//scrapeInflationData()
	//scrapeGasInflationData()
	//scrapeHousingData()
	//end crawler

	//begin scraper
	//fmt.Println("Available domains:")
	//for domainName := range domainConfigurations {
	//	fmt.Printf("- %s\n", domainName)
	//}
	//
	//// Ask the user to choose a domain
	//var domainName string
	//fmt.Print("Enter the domain you want to scrape: ")
	//fmt.Scanln(&domainName)
	//
	//// Check if the chosen domain is valid
	//_, exists := domainConfigurations[domainName]
	//if !exists {
	//	fmt.Printf("Invalid domain name provided: %s\n", domainName)
	//	return
	//}
	//
	//// Perform the scraping for the chosen domain
	//testScrape(domainName)
	//
	////csvread
	//filePath := "crab/csv"
	//properties, err := readCSV(filePath)
	//if err != nil {
	//	fmt.Printf("Error reading CSV file: %s\n", err)
	//	return
	//}
	//
	//// Print the PropertyData for demonstration purposes
	//for _, property := range properties {
	//	fmt.Printf("%+v\n", property)
	//}
}

//end main =============================================================================================================

//bfscrawler test =====================================================================================================
//func bfsCrawlTest(startURLs []string, concurrentCrawlers int) {
//	var wg sync.WaitGroup
//
//	// Initialize the queue with start URLs and mark as visited
//	for _, url := range startURLs {
//		urlQueue <- url
//		visited[url] = true
//	}
//
//	// Start crawling using BFS
//	for i := 0; i < concurrentCrawlers; i++ {
//		wg.Add(1)
//		go func() {
//			for url := range urlQueue {
//				if visited[url] {
//					continue
//				}
//				visited[url] = true
//				crawlURL(URLData{URL: url}, dataChan, &wg)
//			}
//		}()
//	}
//
//	// Close the dataChan when all crawls are complete
//	go func() {
//		wg.Wait()
//		close(dataChan)
//	}()
//
//	// Process the data collected
//	crawledData := make([]URLData, 0)
//	for data := range dataChan {
//		crawledData = append(crawledData, data)
//		// Additional processing can be done here
//	}
//
//	// Convert crawled data to JSON and save to file
//	jsonData, err := json.MarshalIndent(crawledData, "", "  ")
//	if err != nil {
//		log.Fatal("Error marshalling data:", err)
//	}
//	err = ioutil.WriteFile("crawled_data.json", jsonData, 0644)
//	if err != nil {
//		log.Fatal("Error writing to file:", err)
//	}
//
//	// Confirm data has been saved
//	fmt.Println("Crawled data written to crawled_data.json")
//}
//end bfscrawler test =================================================================================================

//add to main if reusing bfs
//	//startURLs := []string{
//	//	"https://www.kaggle.com/search?q=housing+prices",
//	//	"http://books.toscrape.com/",
//	//	"https://www.kaggle.com/search?q=stocks",
//	//	"https://www.kaggle.com/search?q=stock+market",
//	//	"https://www.kaggle.com/search?q=real+estate",
//	//	"https://www.madronavl.com/launchable/public-data-sources-real-estate",
//	//	"https://www.census.gov/programs-surveys/housing.html", "https://data.world/datasets/real-estate",
//	//}
//	//bfsCrawlTest(startURLs, 10)
