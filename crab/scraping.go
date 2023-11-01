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
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"sync"
)

const (
	BaseURL         = "https://www.trackingdifferences.com/ETF/ISIN/"
	AdditionalURL   = "https://www.travelpulse.com/gallery/features/the-15-travel-trends-that-will-define-2023"
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

//type travelInfo struct {
//	Title              string
//	DescriptionOfPlace string
// jjjjjj
//}

func setupCollector() *colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains("www.trackingdifferences.com", "trackingdifferences.com", "www.travelpulse.com"),
	)
}

func scrapeEtfInfo(url string, wg *sync.WaitGroup, resultChan chan EtfInfo) {
	defer wg.Done()

	etfInfo := EtfInfo{}
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

	// Implement other scraping logic here

	// After scraping is complete, send the result back through the channel
	c.OnScraped(func(r *colly.Response) {
		resultChan <- etfInfo
	})

	// Visit the URL
	err := c.Visit(url)
	if err != nil {
		fmt.Printf("Error visiting site: %s", err)
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////////////
// (ETF functions and table creations)
// Function to Create the ETFs Table:
func createETFsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS ETFs (
            etf_id INT AUTO_INCREMENT PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            replication VARCHAR(255),
            earnings VARCHAR(255),
            total_expense_ratio VARCHAR(255),
            tracking_difference VARCHAR(255),
            fund_size VARCHAR(255),
            isin VARCHAR(255) UNIQUE NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `
	_, err := db.Exec(query)
	return err
}

// Function to Insert or Update ETF Data:
func UpdateETFData(db *sql.DB, title, replication, earnings, totalExpenseRatio, trackingDifference, fundSize, isin string) error {
	query := `
        INSERT INTO ETFs (title, replication, earnings, total_expense_ratio, tracking_difference, fund_size, isin)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
        title = VALUES(title),
        replication = VALUES(replication),
        earnings = VALUES(earnings),
        total_expense_ratio = VALUES(total_expense_ratio),
        tracking_difference = VALUES(tracking_difference),
        fund_size = VALUES(fund_size);
    `
	_, err := db.Exec(query, title, replication, earnings, totalExpenseRatio, trackingDifference, fundSize, isin)
	return err
}

// Function to Fetch ETF Data by ISIN:
func fetchETFByISIN(db *sql.DB, isin string) (*sql.Rows, error) {
	query := "CALL FetchETFByISIN(?)"
	return db.Query(query, isin)
}

// Function to Update the Fund Size of an ETF by ISIN:
func updateFundSizeByISIN(db *sql.DB, isin string, fundSize string) error {
	query := "UPDATE ETFs SET fund_size = ? WHERE isin = ?"
	_, err := db.Exec(query, fundSize, isin)
	return err
}

// Function to Count the Number of ETFs:
func countETFs(db *sql.DB) (int, error) {
	query := "SELECT COUNT(*) FROM ETFs"
	var count int
	err := db.QueryRow(query).Scan(&count)
	return count, err
}

// Function to Delete ETF Data by ISIN:
func deleteETFByISIN(db *sql.DB, isin string) error {
	query := "CALL DeleteETFByISIN(?)"
	_, err := db.Exec(query, isin)
	return err
}

// Function to List All ETFs:
func listAllETFs(db *sql.DB) (*sql.Rows, error) {
	query := "CALL ListAllETFs()"
	return db.Query(query)
}

//////////////////////////////////////////////////////////////////////////////////////////////////

func main() {
	urls := []string{BaseURL, AdditionalURL}
	var wg sync.WaitGroup
	resultChan := make(chan EtfInfo, TotalURLs*URLsPerCategory)

	// Launch worker goroutineshh
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
	for etfInfo := range resultChan {
		etfInfos = append(etfInfos, etfInfo)
	}

	// Encode and print the scraped data
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(etfInfos)
}
