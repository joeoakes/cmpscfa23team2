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

	// Test: CreateUser
	userID, err := CreateUser("John Doe", "jdoe", "STD", "password123", true)
	if err != nil {
		log.Printf("Error creating user: %s", err)
	} else {
		log.Printf("User created with ID: %s", userID)
	}
	// Test: FetchUserIDByName
	fetchedUserID, err := FetchUserIDByName("Joesph Oakes")
	if err != nil {
		log.Printf("Error fetching user ID by name: %s", err)
	} else {
		log.Printf("Fetched User ID: %s", fetchedUserID)
	}
	// Test: GetUserByID
	user, err := GetUserByID("8e46234a-631e-11ee-8fa9-30d042e80ac3")
	if err != nil {
		log.Printf("Error fetching user by ID: %s", err)
	} else {
		log.Printf("User Details: %+v", user)
	}

	// Test: GetUsersByRole
	users, err := GetUsersByRole("DEV")
	if err != nil {
		log.Printf("Error fetching users by role: %s", err)
	} else {
		for _, user := range users {
			log.Printf("User by Role: %+v", user)
		}
	}

	// Test: GetAllUsers
	allUsers, err := GetAllUsers()
	if err != nil {
		log.Printf("Error fetching all users: %s", err)
	} else {
		log.Println("All users:")
		for _, user := range allUsers {
			log.Printf("User: %+v", user)
		}
	}

	// Test: ValidateUserCredentials
	isValid, err := ValidateUserCredentials("jdoe", "password123")
	if err != nil {
		log.Printf("Error validating user: %s", err)
	} else if isValid {
		log.Println("User user credentials are valid!")
	} else {
		log.Println("User credentials are invalid!")
	}

	// Test: UpdateUser
	err = UpdateUser(userID, "John Updated", "jupdated", "FAC", ("newpassword123"))
	if err != nil {
		log.Printf("Error updating user: %s", err)
	} else {
		log.Println("User details updated successfully!")
	}

	// Validate the user with updated credentials
	isValid, err = ValidateUserCredentials("jupdated", "newpassword123")
	if err != nil {
		log.Printf("Error validating user after update: %s", err)
	} else if isValid {
		log.Println("User's updated credentials are valid!")
	} else {
		log.Println("User's updated credentials are invalid!")
	}

	// Test: DeleteUser
	err = DeleteUser(userID)
	if err != nil {
		log.Printf("Error deleting user: %s", err)
	} else {
		log.Printf("User with ID %s deleted successfully!", userID)
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

	// testing creating the web crawler
	err = CreateWebCrawler("http://www.abc.com")
	if err != nil {
		fmt.Println("Error creating web crawler:", err)
	} else {
		log.Printf("Success creating the web crawler!")
	}

	// the below tests work but the crab needs to be modified so that it matches the script

	//// testing creating the scraper engine
	//err = CreateScraperEngine("ScraperTest", "This is a test scraper")
	//if err != nil {
	//	fmt.Println("Error creating ScraperEngine:", err)
	//} else {
	//	log.Printf("Success creating the scraper engine!")
	//}
	//
	//// testing insert scraped data
	//// Insert some scraped data (this is just an example, your actual scraping logic will go here)
	//err = InsertScrapedData("http://www.abc.com", "This is some scraped data from site.")
	//if err != nil {
	//	fmt.Println("Error inserting ScrapedData:", err)
	//} else {
	//	log.Printf("Success inserting scraped data")
	//}
}
