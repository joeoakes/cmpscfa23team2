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

var db *sql.DB

type Log struct {
	LogID        string
	StatusCode   string
	Message      string
	GoEngineArea string
	DateTime     []uint8
}

type Scrape struct {
	Title string
	Data  string
	URL   string
}

type JSON_Data_Connect struct {
	Username string
	Password string
	Hostname string
	Database string
}

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

func InsertLog(log Log) error {
	query := "INSERT INTO logs (LogID, StatusCode, Message, GoEngineArea, DateTime) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, log.LogID, log.StatusCode, log.Message, log.GoEngineArea, log.DateTime)
	return err
}

func UpdateLog(log Log) error {
	query := "UPDATE logs SET StatusCode = ?, Message = ?, GoEngineArea = ?, DateTime = ? WHERE LogID = ?"
	_, err := db.Exec(query, log.StatusCode, log.Message, log.GoEngineArea, log.DateTime, log.LogID)
	return err
}

func DeleteLog(logID string) error {
	query := "DELETE FROM logs WHERE LogID = ?"
	_, err := db.Exec(query, logID)
	return err
}

func InsertScrape(scrape Scrape) error {
	query := "INSERT INTO scrapes (Title, Data, URL) VALUES (?, ?, ?)"
	_, err := db.Exec(query, scrape.Title, scrape.Data, scrape.URL)
	return err
}

func UpdateScrape(scrape Scrape) error {
	query := "UPDATE scrapes SET Title = ?, Data = ?, URL = ? WHERE Title = ?"
	_, err := db.Exec(query, scrape.Title, scrape.Data, scrape.URL, scrape.Title)
	return err
}

func DeleteScrape(title string) error {
	query := "DELETE FROM scrapes WHERE Title = ?"
	_, err := db.Exec(query, title)
	return err
}

func main() {
	//test for each of the log functionalities
	newLog := Log{
		LogID:        "1",
		StatusCode:   "200",
		Message:      "Operation successful",
		GoEngineArea: "Auth",
		DateTime:     []uint8("2023-10-01"),
	}
	err := InsertLog(newLog)
	if err != nil {
		log.Fatal("Error inserting new log:", err)
	}
	fmt.Println("Inserted new log.")

	// Update an existing log
	updatedLog := Log{
		LogID:        "1",
		StatusCode:   "201",
		Message:      "updated",
		GoEngineArea: "auth",
		DateTime:     []uint8("2023-10-03"),
	}
	err = UpdateLog(updatedLog)
	if err != nil {
		log.Fatal("Error updating log:", err)
	}
	fmt.Println("Updated log.")

	// Delete a log
	err = DeleteLog("1")
	if err != nil {
		log.Fatal("Error deleting log:", err)
	}
	fmt.Println("Deleted log.")

	//test for each of the scrape functionalities
	newScrape := Scrape{
		Title: "webtitle",
		Data:  "data scraped new",
		URL:   "https://abcupdate.com",
	}
	err = InsertScrape(newScrape)
	if err != nil {
		log.Fatal("Error inserting new scrape:", err)
	}
	fmt.Println("Inserted new scrape.")

	// Update an existing scrape
	updatedScrape := Scrape{
		Title: "webtitle update",
		Data:  "data scraped update",
		URL:   "https://abcupdate.com",
	}
	err = UpdateScrape(updatedScrape)
	if err != nil {
		log.Fatal("Error updating scrape:", err)
	}
	fmt.Println("Updated scrape.")

	// Delete a scrape
	err = DeleteScrape("webtitle")
	if err != nil {
		log.Fatal("Error deleting scrape:", err)
	}
	fmt.Println("Deleted scrape.")
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
