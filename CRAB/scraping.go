package main

import (
	"encoding/json"
	// Import the flag package for command-line argument parsing
	"flag"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"os"
	"time"
)

const (
	// Define the base URL for gas price data (sample site url)
	BaseURL = "https://gasprices.aaa.com/"
	// Set the user agent string
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edg/100.0.0.0 Safari/537.36"
)

// Define a struct to hold gas price information
type GasPriceInfo struct {
	Location string
	Regular  string
	MidGrade string
	Premium  string
	Diesel   string
}

// setupCollector initializes and configures the Colly collector
func setupCollector() *colly.Collector {
	log.Println("Setting up collector...")
	c := colly.NewCollector(
		// Set the user agent for the collector
		colly.UserAgent(UserAgent),
	)
	// Randomize the user agent
	extensions.RandomUserAgent(c)

	// Attach a robot.txt middleware to the collector
	// Set a reasonable timeout for requests
	c.SetRequestTimeout(20 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		// Set the Accept-Language header
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		log.Printf("Visiting %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Printf("Error while scraping: %s\n", e.Error())
	})

	log.Println("Collector setup complete.")
	return c
}

// main function for testing (with the current url):
func main() {
	log.Println("Main function started...")

	// Define a command-line flag for specifying queries or regions
	queryFlag := flag.String("query", "chalfont-pa-gas-prices/", "Specify the query or region")
	// Parse the command-line flags
	flag.Parse()

	// Use the provided query value or the default if not specified
	query := *queryFlag

	gasPriceInfo := GasPriceInfo{}
	gasPriceInfos := make([]GasPriceInfo, 0, 1)

	// Initialize the collector
	c := setupCollector()

	// Set up callback for scraping location
	c.OnHTML(".location", func(h *colly.HTMLElement) {
		gasPriceInfo.Location = h.Text
		log.Printf("Scraped location: %s\n", gasPriceInfo.Location)
	})
	// Set up callback for scraping regular gas price
	c.OnHTML(".regular", func(h *colly.HTMLElement) {
		gasPriceInfo.Regular = h.Text
		log.Printf("Scraped regular price: %s\n", gasPriceInfo.Regular)
	})

	// Assume similar callbacks for MidGrade, Premium, Diesel

	// Callback when a page is scraped
	c.OnScraped(func(r *colly.Response) {
		gasPriceInfos = append(gasPriceInfos, gasPriceInfo)
		// Reset gasPriceInfo for the next round
		gasPriceInfo = GasPriceInfo{}
		log.Println("Scraping complete for this page.")
	})

	targetURL := BaseURL + query
	log.Printf("Processing URL: %s\n", targetURL)

	// Respect robots.txt and visit the URL
	err := c.Visit(targetURL)
	if err != nil {
		log.Printf("Error visiting site: %s", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	// Encode and print gas price data
	enc.Encode(gasPriceInfos)
	log.Println("Main function complete.")
}

// For running:
// go run scraping.go -query "abington-pa-gas-prices/"
