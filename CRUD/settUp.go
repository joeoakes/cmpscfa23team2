package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB // Declare db at the package level

type Log struct {
	LogID        string
	StatusCode   string
	Message      string
	GoEngineArea string
	DateTime     []uint8
}

type JSON_Data_Connect struct {
	Username string
	Password string
	Hostname string
	Database string
}

func init() {
	config, err := readJSONConfig("config.json")
	if err != nil {
		log.Fatal("Error reading JSON config:", err)
	}
	db, err = Connection(config)
	if err != nil {
		log.Fatal("Error establishing database connection:", err)
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
	err = WriteLog("2", "Pos", "Message logged successfully", "Engine1", currentTime) // Notice changed logID to "2"
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

	// Store log using a stored procedure
	err = StoreLog("Success", "Stored using procedure", "Engine1")
	if err != nil {
		log.Println("Failed to store log using stored procedure:", err)
	}
}
