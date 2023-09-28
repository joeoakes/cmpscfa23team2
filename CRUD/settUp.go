package main

/*
	CMPSC 488
	Fall 2023
	PredictAI
*/

// Import required packages
import (
	"database/sql"  // For SQL database interaction
	"encoding/json" // For JSON handling
	"fmt"           // For formatted I/O
	"io/ioutil"     // For I/O utility functions
	"log"           // For logging
	"time"          // For time manipulation

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

// Declare db at the package level for global use
var db *sql.DB

// Log struct models the data structure of a log entry in the database
type Log struct {
	LogID        string
	StatusCode   string
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

// WriteLog writes a log entry to the database
func WriteLog(logID string, statusCode string, message string, goEngineArea string, dateTime time.Time) error {
	// Validate the statusCode by checking if it exists in the `logstatuscodes` table
	var existingStatusCode string
	err := db.QueryRow("SELECT statusCode FROM logstatuscodes WHERE statusCode = ?", statusCode).Scan(&existingStatusCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Invalid statusCode: %s", statusCode)
		}
		return err
	}

	// Prepare the SQL statement for inserting into the log table
	stmt, err := db.Prepare("INSERT INTO log(logID, statusCode, message, goEngineArea, dateTime) VALUES (? ,? ,? ,? ,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, errExec := stmt.Exec(logID, existingStatusCode, message, goEngineArea, dateTime)
	if errExec != nil {
		return errExec
	}

	return nil
}

// InsertStatusCode inserts a new status code into the `logstatuscodes` table
func InsertStatusCode(statusCode, description string) error {
	config, err := readJSONConfig("config.json")
	if err != nil {
		return err
	}
	db, err := Connection(config)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO logstatuscodes(statusCode, statusMessage) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, errExec := stmt.Exec(statusCode, description)
	if errExec != nil {
		return errExec
	}

	return nil
}

// GetLog retrieves all logs from the database
func GetLog() ([]Log, error) {
	stmt, err := db.Prepare("CALL SelectAllLogs()")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []Log
	for rows.Next() {
		var logItem Log
		var dateTimeStr []uint8
		err := rows.Scan(&logItem.LogID, &logItem.StatusCode, &logItem.Message, &logItem.GoEngineArea, &dateTimeStr)
		if err != nil {
			return nil, err
		}
		logItem.DateTime = dateTimeStr
		logs = append(logs, logItem)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetSuccess fetches all logs with a "Success" status code
func GetSuccess() ([]Log, error) {
	stmt, err := db.Prepare("CALL SelectAllLogsByStatusCode(?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("Success")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []Log
	for rows.Next() {
		var logItem Log
		var dateTimeStr []uint8
		err := rows.Scan(&logItem.LogID, &logItem.StatusCode, &logItem.Message, &logItem.GoEngineArea, &dateTimeStr)
		if err != nil {
			return nil, err
		}
		logItem.DateTime = dateTimeStr
		logs = append(logs, logItem)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// InsertOrUpdateStatusCode either inserts a new status code or updates an existing one
func InsertOrUpdateStatusCode(statusCode, description string) error {
	// Check if the status code already exists
	var existingStatusCode string
	err := db.QueryRow("SELECT statusCode FROM logstatuscodes WHERE statusCode = ?", statusCode).Scan(&existingStatusCode)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		// Insert new status code
		stmt, err := db.Prepare("INSERT INTO logstatuscodes(statusCode, statusMessage) VALUES (?, ?)")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(statusCode, description)
		stmt.Close()
		return err
	} else {
		// Update existing status code
		stmt, err := db.Prepare("UPDATE logstatuscodes SET statusMessage = ? WHERE statusCode = ?")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(description, statusCode)
		stmt.Close()
		return err
	}
}

// StoreLog stores a log entry using a stored procedure
func StoreLog(statusCode string, message string, goEngineArea string) error {
	stmt, err := db.Prepare("CALL InsertLog(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, errExec := stmt.Exec(statusCode, message, goEngineArea)
	if errExec != nil {
		return errExec
	}

	return nil
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

// UpdateUser updates an existing user in the database
func UpdateUser(id, name, login, role, password string) error {
	stmt, err := db.Prepare("CALL goengine.update_user(?, ?, ?, ?, ?)") // Updated to match SQL
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, errExec := stmt.Exec(id, name, login, role, password)
	if errExec != nil {
		return errExec
	}

	return nil
}

// DeleteUser removes a user from the database
func DeleteUser(id string) error {
	stmt, err := db.Prepare("CALL goengine.delete_user(?)") // Updated to match SQL
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

	// Insert or Update status code
	err := InsertOrUpdateStatusCode("Pos", "Positive Status")
	if err != nil {
		log.Println("Failed to insert or update status code:", err)
	}

	// Write log
	currentTime := time.Now()
	// This needs to be updated so it adds to a new line each time
	err = WriteLog("3", "Pos", "Message logged successfully", "Engine1", currentTime) // Notice changed logID to "2"
	if err != nil {
		log.Println("Failed to write log:", err)
	}

	// Get and print all logs
	logs, err := GetLog()
	if err != nil {
		log.Println("Failed to get logs:", err)
	} else {
		for _, logItem := range logs {
			fmt.Println(logItem)
		}
	}

	// Store log using a stored procedure (uncomment if needed)
	// err = StoreLog("Success", "Stored using procedure", "Engine1")
	// if err != nil {
	// 	log.Println("Failed to store log using stored procedure:", err)
	// }

	// Insert a new status code
	err = InsertStatusCode("200", "OK")
	if err != nil {
		log.Println("Failed to insert new status code:", err)
	}

	// Create a new user
	err = CreateUser("John", "john123", "admin", "password", true)
	if err != nil {
		log.Println("Failed to create a new user:", err)
	}

	// Update an existing user
	err = UpdateUser("1", "John Doe", "john_doe", "admin", "newpassword")
	if err != nil {
		log.Println("Failed to update user:", err)
	}

	// Delete a user
	err = DeleteUser("1")
	if err != nil {
		log.Println("Failed to delete user:", err)
	}

	// Get and print all "Success" logs
	//successLogs, err := GetSuccess()
	//if err != nil {
	//	log.Println("Failed to get success logs:", err)
	//} else {
	//	fmt.Println("Success Logs:")
	//	for _, logItem := range successLogs {
	//		fmt.Println(logItem)
	//	}
	//}
}
