package main

// Import required packages
import (
	"encoding/json"                    // For JSON handling
	"fmt"                              // For formatted I/O
	_ "github.com/go-sql-driver/mysql" // Import mysql driver
	"log"                              // For logging
	"time"                             // For simulating machine learning model processing time
)

// Prediction struct models the data structure of a prediction in the database
//
// This code defines a struct named "Prediction" with fields for PredictionID, EngineID, InputData, PredictionInfo, and PredictionTime.
type Prediction struct {
	PredictionID   string
	EngineID       string
	InputData      string
	PredictionInfo string
	PredictionTime string
}

// init initializes the program, reading the database configuration and establishing a connection
//func init() {
//	config, err := readJSONConfig("config.json")
//	if err != nil {
//		log.Fatal("Error reading JSON config:", err)
//	} else {
//		log.Println("Successfully read JSON config.")
//	}
//
//	var connErr error
//	db, connErr = Connection(config)
//	if connErr != nil {
//		log.Fatal("Error establishing database connection:", connErr)
//	} else {
//		log.Println("Successfully connected to database.")
//	}
//}

// Connection establishes a new database connection based on provided credentials
//func Connection(config JSON_Data_Connect) (*sql.DB, error) {
//	connDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.Hostname, config.Database))
//	if err != nil {
//		return nil, err
//	} else {
//		log.Println("Successfully opened database connection.")
//	}
//
//	err = connDB.Ping()
//	if err != nil {
//		return nil, err
//	} else {
//		log.Println("Successfully pinged database.")
//	}
//
//	return connDB, nil
//}

// readJSONConfig reads database credentials from a JSON file
//func readJSONConfig(filename string) (JSON_Data_Connect, error) {
//	var config JSON_Data_Connect
//	file, err := ioutil.ReadFile(filename)
//	if err != nil {
//		return config, err
//	} else {
//		log.Println("Successfully read JSON config file.")
//	}
//
//	err = json.Unmarshal(file, &config)
//	if err != nil {
//		return config, err
//	} else {
//		log.Println("Successfully unmarshalled JSON config.")
//	}
//
//	return config, nil
//}

// Function to check if the engine_id exists in scraper_engine table
//
// This function checks if a given engine ID exists in a databse table and returns a boolean indicating existence or an error.
func EngineIDExists(engineID string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM scraper_engine WHERE engine_id=?)"
	err := DB.QueryRow(query, engineID).Scan(&exists)
	if err != nil {
		return false, err
	} else {
		log.Println("Successfully checked if engine ID exists.")
	}
	return exists, nil
}

// Function to insert a new prediction
// The function InsertPrediction, that checks the existence of an engineID, logs the result, and inserts predictionInfo into a database table if the engineID exists, handling errors along the way.
func InsertPrediction(engineID string, predictionInfo string) error {
	exists, err := EngineIDExists(engineID)
	if err != nil {
		return fmt.Errorf("Error checking engine ID: %v", err)
	} else {
		log.Println("Successfully checked if engine ID exists.")
	}
	if !exists {
		return fmt.Errorf("engine_id %s does not exist", engineID)
	} else {
		log.Println("Engine ID exists.")
	}

	query := "INSERT INTO predictions (engine_id, prediction_info) VALUES (?, ?)"
	_, err = DB.Exec(query, engineID, predictionInfo)
	if err != nil {
		return fmt.Errorf("Error storing prediction: %v", err)
	} else {
		log.Println("Successfully inserted prediction.")
	}
	return nil
}

// Function to insert a sample engine ID into scraper_engine table
//
// Function inserts a sample engine's information into a database table, logs success, and returns any encountered errors.
func InsertSampleEngine(engineID, engineName, engineDescription string) error {
	query := "INSERT INTO scraper_engine (engine_id, engine_name, engine_description) VALUES (?, ?, ?)"
	_, err := DB.Exec(query, engineID, engineName, engineDescription)
	if err != nil {
		return fmt.Errorf("Error inserting sample engine: %v", err)
	} else {
		log.Println("Successfully inserted sample engine.")
	}
	return nil
}

// Simulated ML model prediction function
//
// It definesa function that simulates an ML model prediction with a 2-second delay
// and logs a success message before returning a prediction result as a formatted string.
func PerformMLPrediction(inputData string) string {
	// Simulate some delay for ML model prediction
	time.Sleep(2 * time.Second)
	log.Println("Successfully inserted sample engine.")
	return fmt.Sprintf("Prediction result for %s", inputData)
}

// Convert prediction result to JSON
//
// defines a function that converts a given prediction result string into a JSON format, logging a success message and returning the JSON string or an error.
func ConvertPredictionToJSON(predictionResult string) (string, error) {
	predictionMap := map[string]string{"result": predictionResult}
	predictionJSON, err := json.Marshal(predictionMap)
	if err != nil {
		return "", err
	} else {
		log.Println("Successfully converted prediction to JSON.")
	}
	return string(predictionJSON), nil
}

//
//func main() {
//	if db == nil {
//		log.Fatal("Database connection is not initialized.")
//	} else {
//		log.Println("Database connection is initialized.")
//	}
//
//	// Using a WaitGroup for multi-threading
//	var wg sync.WaitGroup
//
//	// Insert a sample engine ID
//	sampleEngineID := "sample_engine_id"
//	sampleEngineName := "Sample Engine"
//	sampleEngineDescription := "This is a sample engine."
//	exists, err := engineIDExists(sampleEngineID)
//	if err != nil {
//		log.Fatalf("Error checking if engine ID exists: %v", err)
//	} else {
//		log.Println("Successfully checked if engine ID exists.")
//	}
//
//	if !exists {
//		err = insertSampleEngine(sampleEngineID, sampleEngineName, sampleEngineDescription)
//		if err != nil {
//			log.Fatalf("Failed to insert sample engine: %v", err)
//		} else {
//			log.Println("Successfully inserted sample engine.")
//		}
//	}
//
//	// Simulate getting some prediction data and performing ML prediction
//	predictionResult := performMLPrediction("Test Data")
//
//	// Convert the prediction result to JSON
//	predictionMap := map[string]string{"result": predictionResult}
//	predictionJSON, err := json.Marshal(predictionMap)
//	if err != nil {
//		log.Fatalf("Failed to convert prediction to JSON: %v", err)
//	}
//	predictionInfo := string(predictionJSON)
//
//	// Use goroutine to insert prediction
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		err := insertPrediction(sampleEngineID, predictionInfo)
//		if err != nil {
//			log.Fatalf("Failed to insert prediction: %v", err)
//		} else {
//			log.Println("Successfully inserted prediction.")
//		}
//	}()
//
//	// Wait for all goroutines to complete
//	wg.Wait()
//}
