package main

/*
	CMPSC 488
	Fall 2023
	PredictAI
*/

// Import required packages
import (
	"database/sql"                     // For SQL database interaction
	"encoding/json"                    // For JSON handling
	"fmt"                              // For formatted I/O
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"io/ioutil"                        // For I/O utility functions
	"log"                              // For logging
)

// Declare db at the package level for global use
var db *sql.DB

// Log struct models the data structure of a log entry in the database
type Log struct {
	LogID        string
	status_code  string
	Message      string
	GoEngineArea string
	DateTime     []uint8
}

// JSON_Data_Connect struct models the structure of database credentials in config.json
type JSON_Data_Connect struct {
	Username string
	Password string
	Hostname string
	Database string
}

// init initializes the program, reading the database configuration and establishing a connection
func init() {
	// Read database credentials from config.json
	config, err := readJSONConfig("config.json")
	if err != nil {
		log.Fatal("Error reading JSON config:", err)
		return
	}

	// Establish a new database connection
	var connErr error
	db, connErr = Connection(config)
	if connErr != nil {
		log.Fatal("Error establishing database connection:", connErr)
	}
}

// Connection establishes a new database connection based on provided credentials
func Connection(config JSON_Data_Connect) (*sql.DB, error) {
	connDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.Hostname, config.Database))
	if err != nil {
		return nil, err
	}

	err = connDB.Ping()
	if err != nil {
		return nil, err
	}

	return connDB, nil
}

// readJSONConfig reads database credentials from a JSON file
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

// User Management Functions

// CreateUser creates a new user in the database
func CreateUser(name, login, role, password string, active bool) error {
	stmt, err := db.Prepare("CALL goengine.create_user(?, ?, ?, ?, ?)") // Updated to match SQL
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, errExec := stmt.Exec(name, login, role, password, active)
	if errExec != nil {
		return errExec
	}

	return nil
}

func FetchUserID(login string) (string, error) {
	var userID string
	query := "SELECT user_id FROM users WHERE user_login = ?"
	fmt.Printf("Executing query: %s with login = %s\n", query, login)
	err := db.QueryRow(query, login).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user with login: %s", login)
		}
		return "", err
	}
	return userID, nil
}

// UpdateUser updates an existing user in the database.
func UpdateUser(name, login, role, password string) error {
	userID, err := FetchUserID(login)
	if err != nil {
		return fmt.Errorf("failed to fetch user ID: %w", err)
	}
	stmt, err := db.Prepare("CALL goengine.update_user(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, errExec := stmt.Exec(userID, name, login, role, password)
	if errExec != nil {
		return fmt.Errorf("failed to update user: %w", errExec)
	}
	return nil
}

// DeleteUser removes a user from the database
func DeleteUser(id string) error {
	stmt, err := db.Prepare("CALL delete_user(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, errExec := stmt.Exec(id)
	if errExec != nil {
		return errExec
	}

	return nil
}

// main function to test all existing methods
func main() {

	// Initialize database connection
	if db == nil {
		log.Fatal("Database connection is not initialized.")
	}

	//Create a new user
	err := CreateUser("Adam", "adam123", "DEV", "aaabbsssbbcc", true)
	if err != nil {
		log.Println("Failed to create a new user:", err)
	}

	_, err = FetchUserID("adam123")
	if err != nil {
		log.Fatalf("Failed to fetch user ID: %v", err)
	}

	// Update User
	err = UpdateUser("NewName", "adam123", "ADM", "newpassword")
	if err != nil {
		fmt.Printf("Failed to update user: %s\n", err)
	} else {
		fmt.Println("Successfully updated user")
	}

	// Delete User
	//err = DeleteUser("NewName")
	//if err != nil {
	//	fmt.Printf("Failed to delete user: %s\n", err)
	//} else {
	//	fmt.Println("Successfully deleted user")
	//}
}
