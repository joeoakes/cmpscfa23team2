package DAL

import (
	"database/sql"
	"encoding/json"
	_ "errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
)

/*
Constants for magic values
And just using this as a placeholder we will be adjusting down the line
*/

// Declare db at the package level for global use
var db *sql.DB

// Log struct models the data structure of a log entry in the database
type Log struct {
	LogID        string
	status_code  string
	Message      string
	GoEngineArea string
	DateTime     []uint8
}

// Prediction struct models the data structure of a prediction in the database
type Prediction struct {
	PredictionID string
	EngineID     string
	InputData    string
	ScrapInfo    string
	ScrapTime    string
}

// JSON_Data_Connect struct models the structure of database credentials in config.json
type JSON_Data_Connect struct {
	Username string
	Password string
	Hostname string
	Database string
}

// init initializes the program, reading the database configuration and establishing a connection
func init() {
	config, err := readJSONConfig("config.json")
	if err != nil {
		log.Fatal("Error reading JSON config:", err)
	}

	var connErr error
	db, connErr = Connection(config)
	if connErr != nil {
		log.Fatal("Error establishing database connection:", connErr)
	}
}

// Connection establishes a new database connection based on provided credentials
func Connection(config JSON_Data_Connect) (*sql.DB, error) {
	connDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.Hostname, config.Database))
	if err != nil {
		return nil, err
	}

	err = connDB.Ping()
	if err != nil {
		return nil, err
	}

	return connDB, nil
}

// readJSONConfig reads database credentials from a JSON file
func readJSONConfig(filename string) (JSON_Data_Connect, error) {
	var config JSON_Data_Connect
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// Function to check if the engine_id exists in scraper_engine table
func engineIDExists(engineID string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM scraper_engine WHERE engine_id=?)"
	err := db.QueryRow(query, engineID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

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
