package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DBConfig struct {
	Username string
	Password string
	HostName string
	Database string
}

var db *sql.DB

func InitDB(config DBConfig) error {
	connDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.HostName, config.Database))
	if err != nil {
		return err
	}
	err = connDB.Ping()
	if err != nil {
		return err
	}
	db = connDB
	return nil
}

func CloseDb() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatalf("Could not establish a connection with the database: %v", err)
		}
	}
}

func main() {
	config := DBConfig{
		Username: "root",
		Password: "Password",
		HostName: "localhost:3306",
		Database: "goengine",
	}
	err := InitDB(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
	}
	defer CloseDb()

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
	//testing inserting engine id
	//sampleEngineID := "sample_engine_id4"
	//sampleEngineName := "Sample Engine"
	//sampleEngineDescription := "This is a sample engine."
	//
	//err = insertSampleEngine(sampleEngineID, sampleEngineName, sampleEngineDescription)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Sample engine with ID %s inserted successfully.\n", sampleEngineID)
	// testing if engine id exists
	//exists, err := engineIDExists(sampleEngineID)
	//if exists {
	//	fmt.Printf("Engine with ID %s exists.\n", "sample_engine_id")
	//} else {
	//	fmt.Printf("Engine with ID %s does not exist.\n", "sample_engine_id")
	//}
	//Testing for inserting prediction
	err = insertPrediction("sample_engine_id", "{\"my\": \"prediction\"}") //need to talk to Hansi about the json
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Prediction for engine %s inserted successfully.\n", "sample_engine_id")

}
