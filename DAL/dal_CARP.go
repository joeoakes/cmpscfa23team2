package DAL

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strings"
)

/*
So far this is a template web scraper that has placeholders for goengine.

Evans Notes
What is needed:
more logging as joe has stated we should be logging everything we can
functional replacements for current placeholders
pagination handling will definately be needed and therefore a future refactor will also be needed
maybe some sort of user-agent rotation but kind of advanced. i was able to get this to work on another project but this
project has a lot more moving parts that should be prioritized

*/
import (
	_"database/sql"
	_"github.com/go-sql-driver/mysql"
)

// Initialize the logger
var logger = logrus.New()

func main() {
	// Set log level
	logger.SetLevel(logrus.InfoLevel)

	// Log messages
	logger.Info("Starting the web scraper...")

	// ...

	// Handle errors with logging
	if err != nil {
		logger.Error("An error occurred: ", err)
	}
}

import (
	_ "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/*
constants for magic values
and just using this as a placeholder we will be adjusting down the line
*/
const (
	BaseURL  = "https://www.trackingdifferences.com/ETF/ISIN/"
	Children = 3
)

// EtfInfo struct holds information about an etf.
type EtfInfo struct {
	Title              string
	Replication        string
	Earnings           string
	TotalExpenseRatio  string
	TrackingDifference string
	FundSize           string
}

// initialize and configure the colly collector. required from colly docs
func setupCollector() *colly.Collector {
	return colly.NewCollector(colly.AllowedDomains("www.trackingdifferences.com", "trackingdifferences.com"))
}

func main() {
	// initial connection to MySQL database
	//rework to read right from config.json
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/goengine")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// prepare SQL statement
	//need adjustment for settUp.go and the actual naming conventions current are placeholders
	stmt, err := db.Prepare("INSERT INTO etf_info(title, replication, earnings, total_expense_ratio, tracking_difference, fund_size) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	//end of the db

	//examples for the curent structure
	isins := []string{"IE00B1XNHC34", "IE00B4L5Y983", "LU1838002480"}
	etfInfo := EtfInfo{}
	etfInfos := make([]EtfInfo, 0, len(isins))

	c := setupCollector()

	// set headers before making the colly requests
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Printf("Visiting %s\n", r.URL)
	})

	// errors handling during scraping
	c.OnError(func(r *colly.Response, e error) {
		// Wrap the original error with additional context
		err := errors.New("Error during scraping: " + e.Error())
		fmt.Println(err)
	})

	// scrape ETF title again placeholder for current working example
	c.OnHTML("h1.page-title", func(h *colly.HTMLElement) {
		etfInfo.Title = h.Text
	})

	// scrapes various ETF attributes
	//again placeholder for current working example
	c.OnHTML("div.descfloat p.desc", func(h *colly.HTMLElement) {
		selection := h.DOM

		// gets all of the child nodes of the selection initiated in magic vars
		childNodes := selection.Children().Nodes

		// checks the number of child nodes (in this case three) to make sure we're looking at the right elements
		if len(childNodes) == 3 {
			description := cleanDesc(selection.Find("span.desctitle").Text())
			value := selection.FindNodes(childNodes[2]).Text()

			// populate etfInfo based on scraped information
			//again placeholder for current working example
			switch description {
			case "Replication":
				etfInfo.Replication = value
			case "TER":
				etfInfo.TotalExpenseRatio = value
			case "TD":
				etfInfo.TrackingDifference = value
			case "Earnings":
				etfInfo.Earnings = value
			case "Fund size":
				etfInfo.FundSize = value
			}
		}
	})

	c.OnScraped(func(r *colly.Response) {
		etfInfos = append(etfInfos, etfInfo)

		// going back to the db after the scrape. here we insert scraped data into goengine
		_, err := stmt.Exec(etfInfo.Title, etfInfo.Replication, etfInfo.Earnings, etfInfo.TotalExpenseRatio, etfInfo.TrackingDifference, etfInfo.FundSize)
		if err != nil {
			log.Printf("Failed to insert data: %s", err)
		}

		etfInfo = EtfInfo{} // Reset etfInfo for the next round
	})

	// this loops through ISINs and perform scraping operations
	for _, isin := range isins {
		if err := c.Visit(scrapeUrl(isin)); err != nil {
			// Wrap the original error with additional context
			err = errors.New("Error visiting site: " + err.Error())
			fmt.Println(err)
			continue
		}
	}

	// encode and print the scraped data
	if err := encodeAndPrint(etfInfos); err != nil {
		// Handle the error appropriately
		fmt.Println("Error encoding and printing ETF information: ", err)
	}
}

// encodeAndPrint encodes the given slice of EtfInfo and prints it
// returns an error if any step fails
func encodeAndPrint(etfInfos []EtfInfo) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	if err := enc.Encode(etfInfos); err != nil {
		return errors.New("Failed to encode ETF information: " + err.Error())
	}
	return nil
}

// cleaner
func cleanDesc(s string) string {
	return strings.TrimSpace(s)
}

// this constructs the URL for scraping based on the ISIN from the top
func scrapeUrl(isin string) string {
	return BaseURL + isin
}
