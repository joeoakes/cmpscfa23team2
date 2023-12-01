package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// This code defines a Go struct named "JSON_Data_Connect" with fields for username, password, hostname,
// and database, each tagged for JSON serialization.
type JSON_Data_Connect struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Hostname string `json:"Hostname"`
	Database string `json:"Database"`
}

var DB *sql.DB

// Read database credentials from a JSON file
func readJSONConfig(filename string) (JSON_Data_Connect, error) {
	var config JSON_Data_Connect
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading config file '%s': %s", filename, err)
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Printf("Error unmarshalling JSON data from file '%s': %s", filename, err)
		return config, err
	}
	log.Println("Successfully read and parsed config file.")
	return config, nil
}

// It initializes a database connection using configuration data from a JSON file and logs any errors encountered during the process.
func InitDB() error {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current working directory: %s", err)
		return err
	}

	// if you are running this from goFrontEnd
	// Construct the path to the config file
	//path := filepath.Join(cwd, "/../../mysql/config.json")

	// if you testing dal:
	path := filepath.Join(cwd, "/../mysql/config.json")

	config, err := readJSONConfig(path)
	if err != nil {
		log.Printf("Error initializing DB from config: %s", err)
		return err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.Hostname, config.Database)
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error opening database with DSN '%s': %s", dsn, err)
		return err
	}

	err = DB.Ping()
	if err != nil {
		log.Printf("Error pinging database: %s", err)
		return err
	}

	log.Println("Database initialized and connected successfully.")
	return nil
}

// defines a function to close a database connection
// and logs any errors or a success message if the connection is closed successfully.
func CloseDb() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Error closing database connection: %s", err)
		} else {
			log.Println("Database connection closed successfully!")
		}
	}
}

