// package main
//
// import (
//
//	"database/sql"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"log"
//
// )
//
//	type DBConfig struct {
//		Username string
//		Password string
//		HostName string
//		Database string
//	}
//
// var db *sql.DB
//
//	func InitDB(config DBConfig) error {
//		connDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.HostName, config.Database))
//		if err != nil {
//			return err
//		}
//		err = connDB.Ping()
//		if err != nil {
//			return err
//		}
//		db = connDB
//		return nil
//	}
//
//	func CloseDb() {
//		if db != nil {
//			err := db.Close()
//			if err != nil {
//				log.Fatalf("Could not establish a connection with the database: %v", err)
//			}
//		}
//	}
//
//	func main() {
//		config := DBConfig{
//			Username: "root",
//			Password: "Pane1901.",
//			HostName: "localhost:3306",
//			Database: "goengine",
//		}
//		err := InitDB(config)
//		if err != nil {
//			log.Fatalf("Failed to initialize database: %s", err)
//		}
//		defer CloseDb()
//
//		// Test CreateUser
//		err = CreateUser("test_name", "test_login", "adm", "test_password", true)
//		if err != nil {
//			log.Printf("Failed to create a user: %s", err)
//		} else {
//			log.Println("Successfully created user.")
//		}
//
//		// Test UpdateUser
//		err = UpdateUser("some_user_id", "updated_name", "login", "dev", "updated_password")
//		if err != nil {
//			log.Printf("Failed to update user: %s", err)
//		} else {
//			log.Println("Successfully updated user.")
//		}
//
//		// Test DeleteUser
//		err = DeleteUser("some_user_id")
//		if err != nil {
//			log.Printf("Failed to delete user: %s", err)
//		} else {
//			log.Println("Successfully deleted user.")
//		}
//		//testing inserting engine id
//		//sampleEngineID := "sample_engine_id4"
//		//sampleEngineName := "Sample Engine"
//		//sampleEngineDescription := "This is a sample engine."
//		//
//		//err = insertSampleEngine(sampleEngineID, sampleEngineName, sampleEngineDescription)
//		//if err != nil {
//		// log.Fatal(err)
//		//}
//		//fmt.Printf("Sample engine with ID %s inserted successfully.\n", sampleEngineID)
//		// testing if engine id exists
//		//exists, err := engineIDExists(sampleEngineID)
//		//if exists {
//		// fmt.Printf("Engine with ID %s exists.\n", "sample_engine_id")
//		//} else {
//		// fmt.Printf("Engine with ID %s does not exist.\n", "sample_engine_id")
//		//}
//		//Testing for inserting prediction
//		err = insertPrediction("sample_engine_id", "{\"my\": \"prediction\"}") //need to talk to Hansi about the json
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		fmt.Printf("Prediction for engine %s inserted successfully.\n", "sample_engine_id")
//
//		// Testing for performing ML prediction
//		predictionResult := performMLPrediction("Test Data")
//		log.Printf(predictionResult)
//		//err = performMLPrediction("sample_input_data")
//		//if err != nil {
//		//	log.Fatal(err)
//		//	log.Printf("Failed to perform ML prediction: %s", err)
//		//} else {
//		//	log.print("Performing ML Prediction is successful!")
//		//}
//
//		// testing for converting prediction to JSON
//		err := convertPredictionToJSON(predictionResult)
//		if err != nil {
//			log.Fatal(err)
//			log.Printf("Failed to convert prediction to JSON", err)
//		} else {
//			log.Printf("Converting prediction to JSON is successful!")
//		}
//
// }
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

	//// Test GetUsers
	//users, err := GetUsers()
	//if err != nil {
	//	log.Printf("Failed to fetch users: %s", err)
	//} else {
	//	for _, u := range users {
	//		log.Printf("Fetched User: %+v", u)
	//	}
	//}
	//testing inserting engine id
	//		//sampleEngineID := "sample_engine_id4"
	//		//sampleEngineName := "Sample Engine"
	//		//sampleEngineDescription := "This is a sample engine."
	//		//
	//		//err = insertSampleEngine(sampleEngineID, sampleEngineName, sampleEngineDescription)
	//		//if err != nil {
	//		// log.Fatal(err)
	//		//}
	//		//fmt.Printf("Sample engine with ID %s inserted successfully.\n", sampleEngineID)
	//		// testing if engine id exists
	//		//exists, err := engineIDExists(sampleEngineID)
	//		//if exists {
	//		// fmt.Printf("Engine with ID %s exists.\n", "sample_engine_id")
	//		//} else {
	//		// fmt.Printf("Engine with ID %s does not exist.\n", "sample_engine_id")
	//		//}
	//		//Testing for inserting prediction
	//		err = insertPrediction("sample_engine_id", "{\"my\": \"prediction\"}") //need to talk to Hansi about the json
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		fmt.Printf("Prediction for engine %s inserted successfully.\n", "sample_engine_id")
	//
	// Testing for performing ML prediction
	predictionResult := performMLPrediction("Test Data")
	log.Printf(predictionResult)

	// testing converting prediction to JSON
	result, err := convertPredictionToJSON(predictionResult)
	if err != nil {
		log.Fatal(err)
		//log.Printf("Failed to convert prediction to JSON", err)
	} else {
		log.Printf("Converting prediction to JSON is successful! %s", result)
	}
}
