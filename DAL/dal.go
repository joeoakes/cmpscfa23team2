package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

type JSON_Data_Connect struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Hostname string `json:"Hostname"`
	Database string `json:"Database"`
}

var db *sql.DB

// Read database credentials from a JSON file
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

func InitDB() error {
	config, err := readJSONConfig("config.json")
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.Hostname, config.Database)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func CloseDb() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing the database: %v", err)
		}
	}
}

func main() {
	// Initialize the database connection
	err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %s", err)
	}
	defer CloseDb()

	// Assuming these functions are elsewhere in your code:
	// Test CreateUser
	err = CreateUser("test_name", "test_login", "adm", "test_password", true)
	if err != nil {
		log.Printf("Failed to create a user: %s", err)
	} else {
		log.Println("Successfully created user.")
	}

	// Test UpdateUser
	err = UpdateUser("some_user_id", "updated_name", "login", "dev", "updated_password")
	if err != nil {
		log.Printf("Failed to update user: %s", err)
	} else {
		log.Println("Successfully updated user.")
	}

	// Test DeleteUser
	err = DeleteUser("some_user_id")
	if err != nil {
		log.Printf("Failed to delete user: %s", err)
	} else {
		log.Println("Successfully deleted user.")
	}

	// Additional functionality goes here
	predictionResult := performMLPrediction("Test Data")
	log.Printf(predictionResult)

	// testing converting prediction to JSON
	result, err := convertPredictionToJSON(predictionResult)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Converting prediction to JSON is successful! %s", result)
	}
}
