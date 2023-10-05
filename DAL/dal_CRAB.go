package main

import (
	"database/sql"
	"encoding/json"
	_ "errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
)

// DBConfig holds the database configuration
type DBConfig struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Port     string `json:"Port"`
	DBName   string `json:"DBName"`
}

var db *sql.DB

func init() {
	// Read the config.json file
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	// Parse the JSON to the DBConfig struct
	var config DBConfig
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Error parsing config file:", err)
	}

	// Initialize the database connection here
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.Username, config.Password, config.Port, config.DBName)
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
}

func CreateWebCrawler(urls string) error {
	_, err := db.Exec("CALL create_webcrawler(?)", urls)
	return err
}

func CreateScraperEngine(urls, id, tags string) error {
	_, err := db.Exec("CALL create_scraper_engine(?, ?)", id, tags)
	return err
}

func InsertScrapedData(url, id, tags string) error {
	_, err := db.Exec("INSERT INTO urls (id, url, tags) VALUES (?, ?)", id, tags)
	return err
}

func main() {
	err := CreateWebCrawler("http://www.abc.com")
	if err != nil {
		fmt.Println("Error creating web crawler:", err)
	}

	err = CreateScraperEngine("http://www.abc.com", "1", "ScraperTest")
	if err != nil {
		fmt.Println("Error creating ScraperEngine:", err)
	}

	// Insert some scraped data (this is just an example, your actual scraping logic will go here)
	err = InsertScrapedData("http://www.abc.com", "1", "ScraperTest")
	if err != nil {
		fmt.Println("Error inserting ScrapedData:", err)
	}
}
