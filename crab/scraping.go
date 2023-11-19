package main

//Error creating ETFs table: Error 1045 (28000): Access denied for user 'root'@'localhost' (using password: YES) ?
// how can we address this error?
import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

// //////////////////////////////////////////////////////////////////////////////////////////////
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

// ////////////////////////////////////////////////////////////////////////////////////////////////
// (storing data proceedures/functions):
// the new functions perform a combination of operations or provide similar functionality
// as the previous functions, but they are not direct duplicates of each other. The
// differences lie in the specific operations they perform or the stored procedures they call.
// (put call in the name of each function so that it is easier to know that they were made for the stored proceedures)

func callInsertETFDataSP(db *sql.DB, title, replication, earnings, totalExpenseRatio, trackingDifference, fundSize, isin string) error {
	query := "CALL InsertETFData(?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, title, replication, earnings, totalExpenseRatio, trackingDifference, fundSize, isin)
	return err
}

func callInsertOrUpdateETFDataSP(db *sql.DB, title, replication, earnings, totalExpenseRatio, trackingDifference, fundSize, isin string) error {
	query := "CALL InsertOrUpdateETFData(?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, title, replication, earnings, totalExpenseRatio, trackingDifference, fundSize, isin)
	return err
}

func callFetchETFByISINSP(db *sql.DB, isin string) (*sql.Rows, error) {
	query := "CALL FetchETFByISIN(?)"
	return db.Query(query, isin)
}

func callUpdateETFByISINSP(db *sql.DB, isin, title, replication, earnings, totalExpenseRatio, trackingDifference, fundSize string) error {
	query := "CALL UpdateETFByISIN(?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, isin, title, replication, earnings, totalExpenseRatio, trackingDifference, fundSize)
	return err
}

func callDeleteETFByISINSP(db *sql.DB, isin string) error {
	query := "CALL DeleteETFByISIN(?)"
	_, err := db.Exec(query, isin)
	return err
}

// Function to display ETF data from the database
func displayETFDataFromDatabase(db *sql.DB) {
	fmt.Println("Displaying ETF data from the database:")

	rows, err := listAllETFs(db)
	if err != nil {
		fmt.Printf("Error fetching ETF data: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var etfInfo EtfInfo
		err := rows.Scan(&etfInfo.Title, &etfInfo.Replication, &etfInfo.Earnings, &etfInfo.TotalExpenseRatio, &etfInfo.TrackingDifference, &etfInfo.FundSize)
		if err != nil {
			fmt.Printf("Error scanning ETF data: %v\n", err)
			return
		}

		fmt.Printf("ETF Title: %s\n", etfInfo.Title)
		fmt.Printf("Replication: %s\n", etfInfo.Replication)
		fmt.Printf("Earnings: %s\n", etfInfo.Earnings)
		fmt.Printf("Total Expense Ratio: %s\n", etfInfo.TotalExpenseRatio)
		fmt.Printf("Tracking Difference: %s\n", etfInfo.TrackingDifference)
		fmt.Printf("Fund Size: %s\n", etfInfo.FundSize)
		fmt.Println("--------------")
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error iterating over ETF data: %v\n", err)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:Pane1901.@tcp(localhost:3306)/mysql")
	if err != nil {
		fmt.Printf("Error opening database connection: %v\n", err)
		return
	}
	defer db.Close()
	urls := []string{BaseURL, AdditionalURL}
	var wg sync.WaitGroup
	resultChan := make(chan interface{}, TotalURLs*URLsPerCategory)

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
	// Create the ETFs table if it doesn't exist
	if err := createETFsTable(db); err != nil {
		fmt.Printf("Error creating ETFs table: %v\n", err)
		return
	}

	// Display ETF data from the database
	displayETFDataFromDatabase(db)
}
