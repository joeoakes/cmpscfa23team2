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

var DB *sql.DB

// DBConfig holds the database configuration
type DBConfig struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Port     string `json:"Port"`
	DBName   string `json:"dbname"`
}

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

func CreateWebCrawler(sourceURL string) error {
	_, err := db.Exec("CALL create_webcrawler(?)", sourceURL)
	return err
}

func CreateScraperEngine(name, description string) error {
	_, err := db.Exec("CALL create_scraper_engine(?, ?)", name, description)
	return err
}

func InsertScrapedData(url, data string) error {
	_, err := db.Exec("INSERT INTO scraped_data (url, data) VALUES (?, ?)", url, data)
	return err
}

func main() {
	err := CreateWebCrawler("http://www.abc.com")
	if err != nil {
		fmt.Println("Error creating web crawler:", err)
	}

	err = CreateScraperEngine("ScraperTest", "This is a test scraper")
	if err != nil {
		fmt.Println("Error creating ScraperEngine:", err)
	}

	// Insert some scraped data (this is just an example, your actual scraping logic will go here)
	err = InsertScrapedData("http://www.abc.com", "This is some scraped data from site.")
	if err != nil {
		fmt.Println("Error inserting ScrapedData:", err)
	}
}
