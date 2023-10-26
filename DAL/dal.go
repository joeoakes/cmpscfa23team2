package DAL

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
)

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

func InitDB() error {
	config, err := readJSONConfig("../config.json")
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
func main() {
	// Initialize the database connection
	err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %s", err)
	}
	defer CloseDb()
	// dal_CARP.go Functions Testing
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
	user, err := GetUserByID(fetchedUserID)
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

	// Test: UpdateUser
	err = UpdateUser(userID, "John Updated", "jupdated", "FAC", ("newpassword123"))
	if err != nil {
		log.Printf("Error updating user: %s", err)
	} else {
		log.Println("User details updated successfully!")
	}

	// Test: DeleteUser
	err = DeleteUser(userID)
	if err != nil {
		log.Printf("Error deleting user: %s", err)
	} else {
		log.Printf("User with ID %s deleted successfully!", userID)
	}
	// CUDA Testing

	//testing inserting engine id
	sampleEngineID := "sample_engine_id10"
	sampleEngineName := "Sample Engine"
	sampleEngineDescription := "This is a sample engine."

	err = InsertSampleEngine(sampleEngineID, sampleEngineName, sampleEngineDescription)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sample engine with ID %s inserted successfully.\n", sampleEngineID)
	//testing if engine id exists
	exists, err := EngineIDExists(sampleEngineID)
	if exists {
		fmt.Printf("Engine with ID %s exists.\n", "sample_engine_id")
	} else {
		fmt.Printf("Engine with ID %s does not exist.\n", "sample_engine_id")
	}
	engineID := "sample_engine_id5"          // Replace with an existing engine ID
	predictionInfo := "{\"key\": \"value\"}" // Replace with valid JSON data

	err = InsertPrediction(engineID, predictionInfo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Prediction for engine %s inserted successfully.\n", engineID)
	// Additional functionality goes here
	predictionResult := PerformMLPrediction("Test Data")
	log.Printf(predictionResult)

	// testing converting prediction to JSON
	result, err := ConvertPredictionToJSON(predictionResult)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Converting prediction to JSON is successful! %s", result)
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

	// dal.CRAB function testing
	// Test: CreateWebCrawler
	crawlerID, err := CreateWebCrawler("http://example.com")
	if err != nil {
		log.Printf("Error creating web crawler: %s", err)
	} else {
		log.Printf("Web crawler created with ID: %s", crawlerID)
	}

	// Test: CreateScraperEngine
	engineID, err = CreateScraperEngine("EngineName", "ScraperEngine")
	if err != nil {
		log.Printf("Error creating scraper engine: %s", err)
	} else {
		log.Printf("Scraper engine created with ID: %s", engineID)
	}

	// Test: InsertURL
	tags := map[string]interface{}{
		"tag1": "value1",
		"tag2": "value2",
	}
	urlID, err := InsertURL("http://example.com/page1", "example.com", tags)
	if err != nil {
		log.Printf("Error inserting URL: %s", err)
	} else {
		log.Printf("URL inserted with ID: %s", urlID)
	}

	// Test: UpdateURL
	updatedTags := map[string]interface{}{
		"tag1": "updatedValue1",
		"tag2": "updatedValue2",
	}
	err = UpdateURL(urlID, "http://example.com/updatedPage", "example.com", updatedTags)
	if err != nil {
		log.Printf("Error updating URL: %s", err)
	} else {
		log.Println("URL updated successfully!")
	}

	// Then, let's fetch it to check the update
	tags, domain, err := GetURLTagsAndDomain(urlID)
	if err != nil {
		log.Printf("Error fetching updated URL details: %s", err)
	} else {
		log.Printf("Domain: %s, Tags: %v", domain, tags)
	}

	// Test: GetURLTagsAndDomain
	returnedTags, returnedDomain, err := GetURLTagsAndDomain(urlID)
	if err != nil {
		log.Printf("Error getting URL tags and domain: %s", err)
	} else {
		log.Printf("Tags: %v, Domain: %s", returnedTags, returnedDomain)
	}

	// Test: GetURLsFromDomain
	urls, err := GetURLsFromDomain("example.com")
	if err != nil {
		log.Printf("Error getting URLs from domain: %s", err)
	} else {
		log.Printf("URLs from domain: %v", urls)
	}

	// Test: GetUUIDFromURLAndDomain
	//uuid, err := GetUUIDFromURLAndDomain("http://example.com/page1", "example.com")
	//if err != nil {
	//	log.Printf("Error getting UUID from URL and domain: %s", err)
	//} else {
	//	log.Printf("UUID for given URL and domain: %s", uuid)
	//}

	// Test: GetRandomURL
	//randomURL, err := GetRandomURL()
	//if err != nil {
	//	log.Printf("Error getting random URL: %s", err)
	//} else {
	//	log.Printf("Random URL: %s", randomURL)
	//}

	// Test: GetURLsOnly
	allURLs, err := GetURLsOnly()
	if err != nil {
		log.Printf("Error getting all URLs: %s", err)
	} else {
		log.Printf("All URLs: %v", allURLs)
	}

	// Test: GetURLsAndTags
	urlsAndTags, err := GetURLsAndTags()
	if err != nil {
		log.Printf("Error getting URLs and tags: %s", err)
	} else {
		log.Printf("All URLs and tags: %v", urlsAndTags)
	}

}
