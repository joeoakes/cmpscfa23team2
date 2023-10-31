//package main
//
//import (
//	"encoding/json"
//	_ "errors"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/gocolly/colly"
//	"os"
//	"strings"
//)
//
///*
//Constants for magic values
//And just using this as a placeholder we will be adjusting down the line
//*/
//const (
//	BaseURL  = "https://www.trackingdifferences.com/ETF/ISIN/"
//	Children = 3
//)
//
//// EtfInfo struct holds information about an ETF.
//type EtfInfo struct {
//	Title              string
//	Replication        string
//	Earnings           string
//	TotalExpenseRatio  string
//	TrackingDifference string
//	FundSize           string
//}
//
//// Initialize and configure the Colly collector.
//func setupCollector() *colly.Collector {
//	return colly.NewCollector(colly.AllowedDomains("www.trackingdifferences.com", "trackingdifferences.com"))
//}
//
//func main() {
//	isins := []string{"IE00B1XNHC34", "IE00B4L5Y983", "LU1838002480"}
//	etfInfo := EtfInfo{}
//	etfInfos := make([]EtfInfo, 0, len(isins))
//
//	c := setupCollector()
//
//	// Set headers before making a request
//	c.OnRequest(func(r *colly.Request) {
//		r.Headers.Set("Accept-Language", "en-US;q=0.9")
//		fmt.Printf("Visiting %s\n", r.URL)
//	})
//
//	// Handle any errors during scraping
//	c.OnError(func(r *colly.Response, e error) {
//		fmt.Printf("Error while scraping: %s\n", e.Error())
//	})
//
//	// Scrape ETF title
//	c.OnHTML("h1.page-title", func(h *colly.HTMLElement) {
//		etfInfo.Title = h.Text
//	})
//
//	// Scrape various ETF attributes
//	c.OnHTML("div.descfloat p.desc", func(h *colly.HTMLElement) {
//		selection := h.DOM
//
//		// Get all child nodes of the selection
//		childNodes := selection.Children().Nodes
//
//		// Check the number of child nodes to make sure we're looking at the right elements
//		if len(childNodes) == 3 {
//			description := cleanDesc(selection.Find("span.desctitle").Text())
//			value := selection.FindNodes(childNodes[2]).Text()
//
//			// Populate etfInfo based on scraped information
//			switch description {
//			case "Replication":
//				etfInfo.Replication = value
//			case "TER":
//				etfInfo.TotalExpenseRatio = value
//			case "TD":
//				etfInfo.TrackingDifference = value
//			case "Earnings":
//				etfInfo.Earnings = value
//			case "Fund size":
//				etfInfo.FundSize = value
//			}
//		}
//	})
//
//	// After scraping is complete, append the populated etfInfo to the slice
//	c.OnScraped(func(r *colly.Response) {
//		etfInfos = append(etfInfos, etfInfo)
//		etfInfo = EtfInfo{} // Reset etfInfo for the next round
//	})
//
//	// Loop through ISINs and perform scraping operations
//	for _, isin := range isins {
//		if err := c.Visit(scrapeUrl(isin)); err != nil {
//			fmt.Printf("Error visiting site: %s", err)
//			continue
//		}
//	}
//
//	// Encode and print the scraped data
//	enc := json.NewEncoder(os.Stdout)
//	enc.SetIndent("", " ")
//	enc.Encode(etfInfos)
//}
//
//// Cleans up description strings
//func cleanDesc(s string) string {
//	return strings.TrimSpace(s)
//}
//
//// Constructs the URL for scraping based on the ISIN
//func scrapeUrl(isin string) string {
//	return BaseURL + isin
//}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"sync"
)

const (
	BaseURL         = "https://www.trackingdifferences.com/ETF/ISIN/"
	AdditionalURL   = "https://www.cntraveler.com/story/best-places-to-go-in-2023"
	Workers         = 3
	TotalURLs       = 2 // Total number of URLs to scrape
	URLsPerCategory = 1 // Number of URLs per category (you can adjust this based on your needs)
)

type EtfInfo struct {
	Title              string
	Replication        string
	Earnings           string
	TotalExpenseRatio  string
	TrackingDifference string
	FundSize           string
}

type TravelInfo struct {
	Title              string
	DescriptionOfPlace string
}

func setupCollector() *colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains("www.trackingdifferences.com", "trackingdifferences.com", "www.cntraveler.com"),
	)
}

func scrapeEtfInfo(url string, wg *sync.WaitGroup, resultChan chan interface{}) {
	defer wg.Done()

	if url == BaseURL {
		etfInfo := EtfInfo{} // Initialize with empty fields
		c := setupCollector()

		// Set headers before making a request
		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Accept-Language", "en-US;q=0.9")
			fmt.Printf("Visiting %s\n", r.URL)
		})

		// Handle any errors during scraping
		c.OnError(func(r *colly.Response, e error) {
			fmt.Printf("Error while scraping: %s\n", e.Error())
		})

		// Implement scraping logic for ETF info here
		// For example:
		// c.OnHTML("your ETF info CSS selector", func(e *colly.HTMLElement) {
		//   etfInfo.Title = e.ChildText("title CSS selector")
		//   etfInfo.Replication = e.ChildText("replication CSS selector")
		//   // ... populate other fields similarly
		// })

		// After scraping is complete, send the result back through the channel
		c.OnScraped(func(r *colly.Response) {
			resultChan <- etfInfo
		})

		// Visit the ETF URL
		err := c.Visit(url)
		if err != nil {
			fmt.Printf("Error visiting ETF site: %s", err)
		}
	} else if url == AdditionalURL {
		travelInfo := TravelInfo{} // Initialize with empty fields
		c := setupCollector()

		// Set headers before making a request
		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Accept-Language", "en-US;q=0.9")
			fmt.Printf("Visiting %s\n", r.URL)
		})

		// Handle any errors during scraping
		c.OnError(func(r *colly.Response, e error) {
			fmt.Printf("Error while scraping: %s\n", e.Error())
		})

		// Implement scraping logic for travel info here
		c.OnHTML("h1.headline", func(e *colly.HTMLElement) {
			fmt.Println("Title found:", e.Text)
			travelInfo.Title = e.Text
		})

		c.OnHTML(".dekText", func(e *colly.HTMLElement) {
			fmt.Println("Description found:", e.Text)
			travelInfo.DescriptionOfPlace = e.Text
		})

		// After scraping is complete, send the result back through the channel
		c.OnScraped(func(r *colly.Response) {
			resultChan <- travelInfo
		})

		// Visit the travel URL
		err := c.Visit(url)
		if err != nil {
			fmt.Printf("Error visiting travel site: %s", err)
		}
	}
}

func main() {
	urls := []string{BaseURL, AdditionalURL}
	var wg sync.WaitGroup
	resultChan := make(chan interface{}, TotalURLs*URLsPerCategory)

	// Launch worker goroutines
	for i := 0; i < Workers; i++ {
		wg.Add(1)
		go func() {
			for _, url := range urls {
				for j := 0; j < URLsPerCategory; j++ {
					scrapeEtfInfo(url, &wg, resultChan)
				}
			}
		}()
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from the channel
	var etfInfos []EtfInfo
	var travelInfos []TravelInfo
	for result := range resultChan {
		switch data := result.(type) {
		case EtfInfo:
			etfInfos = append(etfInfos, data)
		case TravelInfo:
			travelInfos = append(travelInfos, data)
		}
	}

	// Encode and print the scraped data
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")

	// Output ETF information
	enc.Encode(etfInfos)

	// Output travel information
	enc.Encode(travelInfos)
}
