package dal

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// Log struct models the data structure of a log entry in the database
//
// This code defines a Go struct named "Log" with fields for
// log ID, status code, message, GoEngineArea, and date-time information.
type Log struct {
	LogID        string
	status_code  string
	Message      string
	GoEngineArea string
	DateTime     []uint8
}

// This code defines a Go struct named LogStatusCodes with two fields, "StatusCode" and
// "StatusMessage," to represent status code and associated status messages.
type LogStatusCodes struct {
	StatusCode    string
	StatusMessage string
}

// Function to insert a log entry into the database
//
// It  inserts a log entry into a database using a SQL stored procedure, handling any errors that may occur during the execution.
func InsertLog(statusCode, message, goEngineArea string) {
	_, err := DB.Exec("CALL insert_log(?, ?, ?)", statusCode, message, goEngineArea)
	if err != nil {
		log.Println("Error inserting log:", err)
	}
}

// This function creates & adds the log entries to a TextFile if the database is down
func init() {
	// Initialize the database first
	if err := InitDB(); err != nil {
		InsertLog("400", "Failed to initialize database", "init()")
		log.Fatal(err)
	}

	file, err := os.OpenFile("Logging.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		InsertLog("400", "Failed to open file", "init()")
		log.Fatal(err)
	} else {
		InsertLog("200", "INIT Open File Success", "init()")
	}

	log.SetOutput(file)
}

// WriteLog writes a log entry to the database
//
// This code defines a function WriteLog that validates a status code, inserts a log entry into a database,
// and logs the execution process, handling potential errors along the way.
func WriteLog(logID string, status_code string, message string, goEngineArea string, dateTime time.Time) error {
	// Validate the statusCode by checking if it exists in the `log_status_codes` table
	var existingStatusCode string
	err := DB.QueryRow("SELECT status_code FROM log_status_codes WHERE status_code = ?", status_code).Scan(&existingStatusCode)
	if err != nil {
		InsertLog("400", "Failed to query row", "WriteLog()")
		if err == sql.ErrNoRows {
			InsertLog("400", "Invalid statusCode", "WriteLog()")
			return fmt.Errorf("Invalid statusCode: %s", status_code)
		} else {
			InsertLog("200", "Successfully validated status code", "WriteLog()")
		}

		return err
	}
	// Prepare the SQL statement for inserting into the log table
	stmt, err := DB.Prepare("INSERT INTO log(log_ID, status_code, message, go_engine_area, date_time) VALUES (? ,? ,? ,? ,?)")
	if err != nil {
		InsertLog("400", "Failed to prepare SQL statement", "WriteLog()")
		return err
	} else {
		InsertLog("200", "Successfully prepared SQL statement", "WriteLog()")
	}

	defer stmt.Close()

	// Execute the SQL statement
	_, errExec := stmt.Exec(logID, existingStatusCode, message, goEngineArea, dateTime)
	if errExec != nil {
		InsertLog("400", "Failed to execute SQL statement", "WriteLog()")
		return errExec
	} else {
		InsertLog("200", "Successfully executed SQL statement", "WriteLog()")
	}

	return nil
}

// GetLog - Reads the log
//
// This Go code defines a function, "GetLog," that prepares and queries a database for logs, logging both successful and failed operations,
// and returns a log objects along with potential errors.
func GetLog() ([]Log, error) {
	stmt, err := DB.Prepare("CALL select_all_logs()")
	if err != nil {
		InsertLog("400", "Failed to prepare SQL statement", "GetLog()")
		return nil, err
	} else {
		InsertLog("200", "Successfully prepared SQL statement", "GetLog()")
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		InsertLog("400", "Failed to query SQL statement", "GetLog()")
		return nil, err
	} else {
		log.Println("Successfully queried SQL statement")
	}
	defer rows.Close()

	var logs []Log
	for rows.Next() {
		var logItem Log
		var dateTimeStr []uint8
		err := rows.Scan(&logItem.LogID, &logItem.status_code, &logItem.Message, &logItem.GoEngineArea, &dateTimeStr)
		if err != nil {
			return nil, err
		} else {
			InsertLog("200", "Successfully scanned rows", "GetLog()")
		}

		logItem.DateTime = dateTimeStr
		logs = append(logs, logItem)
		InsertLog("200", "Successfully appended log item", "GetLog()")
	}

	if err = rows.Err(); err != nil {
		InsertLog("400", "Failed to iterate over rows", "GetLog()")
		return nil, err
	} else {
		InsertLog("200", "Successfully iterated over rows", "GetLog()")
	}

	return logs, nil
}

// It defines  defines a function that executes a SQL stored procedure "insert_or_update_status_code" with provided parameters "statusCode"
// and "statusMessage" using the "DB" database connection and returns any potential errors.
func InsertOrUpdateStatusCode(statusCode, statusMessage string) error {
	_, err := DB.Exec("CALL insert_or_update_status_code(?, ?)", statusCode, statusMessage)
	return err
}

// GetSuccess - Uses a Procedure to gather all the 'Success' rows in the DB
//
// The code defines a function GetSuccess that retrieves log entries with a "Success" status code from a database, logs various status messages.
func GetSuccess() ([]Log, error) {
	stmt, err := DB.Prepare("CALL select_all_logs_by_status_code(?)")
	if err != nil {
		InsertLog("400", "Failed to prepare SQL statement", "GetSuccess()")
		return nil, err
	} else {
		InsertLog("200", "Successfully prepared SQL statement", "GetSuccess()")
	}
	defer stmt.Close()

	rows, err := stmt.Query("200")
	if err != nil {
		InsertLog("400", "Failed to query SQL statement", "GetSuccess()")
		return nil, err
	} else {
		InsertLog("200", "Successfully queried SQL statement", "GetSuccess()")
	}
	defer rows.Close()
	InsertLog("200", "Successfully closed rows", "GetSuccess()")
	var logs []Log
	for rows.Next() {
		var logItem Log
		var dateTimeStr []uint8
		err := rows.Scan(&logItem.LogID, &logItem.status_code, &logItem.Message, &logItem.GoEngineArea, &dateTimeStr)
		if err != nil {
			InsertLog("400", "Failed to scan rows successfully", "GetSuccess()")
			return nil, err
		} else {
			InsertLog("200", "Successfully scanned rows", "GetSuccess()")
		}
		logItem.DateTime = dateTimeStr
		logs = append(logs, logItem)
		InsertLog("200", "Successfully appended log item", "GetSuccess()")
	}

	if err = rows.Err(); err != nil {
		InsertLog("400", "Failed to iterate over rows", "GetSuccess()")
		return nil, err
	} else {
		InsertLog("200", "Successfully iterated over rows", "GetSuccess()")
	}

	return logs, nil
}

// This code prepares and executes a SQL statement to store log information in a database, logging the status of the SQL operations during the process
func StoreLog(status_code string, message string, goEngineArea string) error {
	stmt, err := DB.Prepare("CALL insert_log(?,?,?)")
	if err != nil {
		InsertLog("400", "Failed to prepare SQL statement", "StoreLog()")
		return err
	} else {
		InsertLog("200", "Successfully prepared SQL statement", "StoreLog()")
	}
	defer stmt.Close()

	_, errExec := stmt.Exec(status_code, message, goEngineArea)
	if errExec != nil {
		InsertLog("400", "Failed to iterate over rows", "GetSuccess()")
		return errExec
	} else {
		InsertLog("200", "Successfully executed SQL statement", "StoreLog()")
	}

	return nil
}
