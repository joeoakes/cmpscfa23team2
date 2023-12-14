package crab

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// urlQueue is a channel used for queuing URLs to be processed by the scraper.
var urlQueue = make(chan string, 1000)

// visited is a map used for keeping track of URLs that have already been visited by the scraper.
var visited = make(map[string]bool)

// domainConfigurations maps domain names to their respective scraping configurations. Each domain
// has specific selectors and structures defined for scraping its relevant data.
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

// ReadCSV reads a CSV file from the given file path, parses the data, and returns a slice of PropertyData.
// It returns an error if it fails to read or parse the CSV file. This function is designed to handle CSV files
// with a specific format for real estate data.
func ReadCSV(filePath string) ([]PropertyData, error) {
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

// Scrape performs the web scraping process for a given domain. It takes a URL to start scraping from,
// a DomainConfig for scraping rules, and a WaitGroup for concurrency control. The function collects
// scraped data and saves it to a JSON file.
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
	maxRetries := 6
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
	err := InsertData(ItemData{
		Domain: domainConfig.Name,
		Data:   allData,
	}, filename)
	if err != nil {
		fmt.Printf("Error saving data to JSON file: %v\n", err)
	}
}

//end scrape ===========================================================================================================

// testScrape is a testing function for the scraper. It takes a domain name and triggers the Scrape
// function using predefined test URLs for the domain. This function helps in validating the scraping logic
// for different domains.
func TestScrape(domainName string) {
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

// The following functions (airdatatest, scrapeInflationData, scrapeGasInflationData, scrapeHousingData)
// are specific scraper implementations for different types of data like airfare, inflation, gasoline prices,
// and housing data. Each function fetches data from specific URLs and processes it according to predefined
// scraping rules and selectors, then writes the scraped data to JSON files.
func Airdatatest() {
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
func ScrapeInflationData() {
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
func ScrapeGasInflationData() {
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
func ScrapeHousingData() {
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
