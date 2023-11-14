// Only had about an hour or two to implement. This should be a more agnostic way to do this. Someone with the updated
// crawler should try an implementation with this for testing.
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // Adjust the import based on your database
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// URLData struct for storing scraped data
type URLData struct {
	URL     string
	Tags    map[string]interface{}
	Domain  string
	Created time.Time
}

// GetURLsOnly calls a stored procedure to retrieve a list of URLs.
func GetURLsOnly(db *sql.DB) ([]string, error) {
	var urls []string

	// Call the stored procedure that returns URLs. Make sure the procedure returns a single column of URLs.
	rows, err := db.Query("CALL GetURLsOnly()") // Replace GetURLList with your actual stored procedure name.
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows and add URLs to the list.
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

// GetURLTagsAndDomain calls a stored procedure to retrieve tags and the domain for a URL.
func GetURLTagsAndDomain(db *sql.DB, url string) (map[string]interface{}, string, error) {
	var (
		tagsJSON string // Assuming the stored procedure returns tags as a JSON string.
		domain   string
	)

	// Call the stored procedure passing the URL as a parameter. Make sure to replace GetURLTagsDomain with your actual stored procedure name.
	err := db.QueryRow("CALL GetURLTagsDomain(?)", url).Scan(&tagsJSON, &domain)
	if err != nil {
		return nil, "", err
	}

	// Decode the JSON string back to a map.
	var tags map[string]interface{}
	err = json.Unmarshal([]byte(tagsJSON), &tags)
	if err != nil {
		return nil, "", err
	}

	return tags, domain, nil
}

// insertIntoDB inserts a URLData record into the database
func insertIntoDB(db *sql.DB, data URLData) error {
	// Replace with your actual DB logic and SQL query
	query := `INSERT INTO scraped_data (url, tags, domain, created_at) VALUES (?, ?, ?, ?)`

	// Convert Tags to JSON for storage
	tagsJSON, err := json.Marshal(data.Tags)
	if err != nil {
		return err
	}

	// Execute the SQL query for inserting data
	_, err = db.Exec(query, data.URL, tagsJSON, data.Domain, data.Created)
	return err
}

func main() {
	// Setup the database connection
	db, err := sql.Open("goengine", "user:password@/dbname")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	var wg sync.WaitGroup

	// OnHTML callback to process each visited page
	c.OnHTML("body", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()

		tags, domain, err := GetURLTagsAndDomain(url)
		if err != nil {
			log.Printf("Error getting tags and domain for %s: %s", url, err)
			return
		}

		data := URLData{
			URL:     url,
			Tags:    tags,
			Domain:  domain,
			Created: time.Now(),
		}

		// Insert the data into the database
		err = insertIntoDB(db, data)
		if err != nil {
			log.Printf("Failed to insert data into DB: %s", err)
			return
		}

		// Print the collected data
		fmt.Printf("Scraped URLData: %+v\n", data)
		wg.Done()
	})

	// Get list of URLs to scrape using the stored procedure
	urls, err := GetURLsOnly(db)
	if err != nil {
		log.Fatal("Error getting URLs:", err)
	}

	// Start crawling each URL
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			err := c.Visit(url)
			if err != nil {
				log.Printf("Error visiting %s: %s", url, err)
				wg.Done()
			}
		}(url)
	}

	// Wait for all crawling to complete
	wg.Wait()
}

//Example SPROCS
//CREATE TABLE scraped_data (
//id INT AUTO_INCREMENT PRIMARY KEY,
//url VARCHAR(2048) NOT NULL,
//tags JSON,
//domain VARCHAR(255),
//created_at DATETIME
//);
//
//DELIMITER //
//
//CREATE PROCEDURE GetURLList()
//BEGIN
//SELECT url FROM urls_to_scrape;
//END //
//
//DELIMITER ;
//
//DELIMITER //
//
//CREATE PROCEDURE InsertScrapedData(IN _url VARCHAR(2048), IN _tags JSON, IN _domain VARCHAR(255), IN _created_at DATETIME)
//BEGIN
//INSERT INTO scraped_data (url, tags, domain, created_at) VALUES (_url, _tags, _domain, _created_at);
//END //
//
//DELIMITER ;
