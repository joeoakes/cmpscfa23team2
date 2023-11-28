package main

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
func InsertLog(statusCode, message, goEngineArea string) error {
	_, err := DB.Exec("CALL insert_log(?, ?, ?)", statusCode, message, goEngineArea)
	if err != nil {
		log.Println("Error inserting log:", err)
	}
	return err
}

// This function creates & adds the log entries to a TextFile if the database is down
func init() {
	// Initialize the database first
	if err := InitDB(); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("Logging_to.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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

		InsertLog("200", "Successfully validated status code", "WriteLog()")
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

	rows, err := stmt.Query("Success")
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

//
//func main() {
//	// Initialize the database connection
//	err := InitDB()
//	if err != nil {
//		log.Fatalf("Failed to initialize the database: %v", err)
//	}
//	// The rest of your code
//	defer CloseDb()
//
//	// Insert or Update status code
//	err = InsertOrUpdateStatusCode("POS", "noth")
//	if err != nil {
//		log.Println("Failed to insert or update status code:", err)
//		InsertLog("200", "Failed to insert or update status code", "main()")
//	} else {
//		InsertLog("200", "Successfully inserted or updated status code", "main()")
//	}
//
//	_, err = FetchUserIDByName("Joesph Oakes")
//	if err != nil {
//		log.Fatalf("Failed to fetch user ID: %v", err)
//		InsertLog("200", "Failed to fetch user ID", "main()")
//	} else {
//		InsertLog("200", "Successfully fetched user ID", "main()")
//	}
//
//	// Update User
//	err = UpdateUser("NewName", "jxo19", "ADM", "ADM", "password")
//	if err != nil {
//		fmt.Printf("Failed to update user: %s\n", err)
//		InsertLog("200", "Failed to update user", "main()")
//	} else {
//		InsertLog("200", "Successfully updated user", "main()")
//	}
//
//	// Delete User
//	err = DeleteUser("jxo19")
//	if err != nil {
//		fmt.Printf("Failed to delete user: %s\n", err)
//		InsertLog("200", "Failed to delete user", "main()")
//	} else {
//		InsertLog("200", "Successfully deleted user", "main()")
//	}
//
//	//Generate a unique logID
//	uniqueLogID := uuid.New().String()
//
//	//Write log
//	currentTime := time.Now()
//	err = WriteLog(uniqueLogID, "Pos", "Message logged successfully", "Engine1", currentTime)
//	if err != nil {
//		log.Println("Failed to write log:", err)
//		InsertLog("200", "Failed to write log", "main()")
//	} else {
//		InsertLog("200", "Successfully wrote log", "main()")
//	}
//
//	// Get and print all logs
//	logs, err := GetLog()
//	if err != nil {
//		log.Println("Failed to get logs:", err)
//		InsertLog("200", "Failed to get logs", "main()")
//	} else {
//		for _, logItem := range logs {
//			fmt.Println(logItem)
//			err := InsertLog("200", "Successfully got logs", "main()")
//			if err != nil {
//				return
//			}
//		}
//	}
//
//	//Store log using a stored procedure (uncomment if needed)
//	err = StoreLog("200", "Stored using procedure", "Engine1")
//	if err != nil {
//		log.Println("Failed to store log using stored procedure:", err)
//		InsertLog("200", "Failed to store log using stored procedure", "main()")
//	} else {
//		InsertLog("200", "Successfully stored log using stored procedure", "main()")
//	}
//
//	//Create a new user
//	_, err = CreateUser("John", "john123", "ADM", "password", true)
//	if err != nil {
//		log.Println("Failed to create a new user:", err)
//		InsertLog("200", "Failed to create a new user", "main()")
//	} else {
//		InsertLog("200", "Successfully created a new user", "main()")
//	}
//
//}

// readJSONConfig - Reads the JSON config file
//func readJSONConfig(filename string) (JsonDataConnect, error) {
//	var config JsonDataConnect
//	file, err := ioutil.ReadFile(filename)
//	if err != nil {
//		InsertLog(db, "200", "Failed to read JSON config file", "readJSONConfig()")
//		return config, err
//	} else {
//		InsertLog(db, "200", "Successfully read JSON config file", "readJSONConfig()")
//	}
//
//	err = json.Unmarshal(file, &config)
//	if err != nil {
//		InsertLog(db, "200", "Failed to unmarshal JSON config", "readJSONConfig()")
//		return config, err
//	} else {
//		InsertLog(db, "200", "Successfully unmarshalled JSON config", "readJSONConfig()")
//	}
//
//	return config, nil
//}

// Function to read the JSON configuration file
//func readJSONConfig(configFile string, key []byte) (*JsonDataConnect, error) {
//	data, err := ioutil.ReadFile(configFile)
//	if err != nil {
//		InsertLog("200", "Failed to read JSON config file", "ReadJSONConfig()")
//		return nil, err
//	} else {
//		InsertLog("200", "Successfully read JSON config file", "ReadJSONConfig()")
//	}
//
//	var encryptedConfig JsonDataConnect
//	err = json.Unmarshal(data, &encryptedConfig)
//	if err != nil {
//		InsertLog("200", "Failed to unmarshal JSON config", "ReadJSONConfig()")
//		return nil, err
//	} else {
//		InsertLog("200", "Successfully unmarshalled JSON config", "ReadJSONConfig()")
//	}
//
//	decryptedUsername, err := base64.StdEncoding.DecodeString(encryptedConfig.Username)
//	if err != nil {
//		InsertLog("200", "Failed to decode username", "ReadJSONConfig()")
//		return nil, err
//	} else {
//		InsertLog("200", "Successfully decoded username", "ReadJSONConfig()")
//	}
//
//	decryptedPassword, err := base64.StdEncoding.DecodeString(encryptedConfig.Password)
//	if err != nil {
//		InsertLog("200", "Failed to decode password", "ReadJSONConfig()")
//		return nil, err
//	} else {
//		InsertLog("200", "Successfully decoded password", "ReadJSONConfig()")
//	}
//
//	username, err := decryptAES(decryptedUsername, key)
//	if err != nil {
//		return nil, err
//	} else {
//		InsertLog("200", "Successfully decrypted username", "ReadJSONConfig()")
//	}
//
//	password, err := decryptAES(decryptedPassword, key)
//	if err != nil {
//		InsertLog("200", "Failed to decrypt password", "ReadJSONConfig()")
//		return nil, err
//	} else {
//		InsertLog("200", "Successfully decrypted password", "ReadJSONConfig()")
//	}
//
//	decryptedConfig := JsonDataConnect{
//		Username: string(username),
//		Password: string(password),
//		Hostname: encryptedConfig.Hostname,
//		Database: encryptedConfig.Database,
//	}
//	InsertLog("200", "Successfully decrypted config", "ReadJSONConfig()")
//
//	return &decryptedConfig, nil
//}

// Function to write the JSON configuration file with encrypted username and password
//func WriteJSONConfig(configFile string, config *JsonDataConnect, key []byte) error {
//	encryptedUsername, err := encryptAES([]byte(config.Username), key)
//	if err != nil {
//		return err
//	} else {
//		InsertLog("200", "Successfully encrypted username", "WriteJSONConfig()")
//	}
//
//	encryptedPassword, err := encryptAES([]byte(config.Password), key)
//	if err != nil {
//		InsertLog("200", "Failed to encrypt password", "WriteJSONConfig()")
//		return err
//	} else {
//		InsertLog("200", "Successfully encrypted password", "WriteJSONConfig()")
//	}
//
//	encryptedConfig := JsonDataConnect{
//		Username: base64.StdEncoding.EncodeToString(encryptedUsername),
//		Password: base64.StdEncoding.EncodeToString(encryptedPassword),
//		Hostname: config.Hostname,
//		Database: config.Database,
//	}
//	InsertLog("200", "Successfully encrypted config", "WriteJSONConfig()")
//	// Marshal the encrypted configuration struct to JSON
//	data, err := json.Marshal(encryptedConfig)
//	if err != nil {
//		InsertLog("200", "Failed to marshal encrypted config", "WriteJSONConfig()")
//		return err
//	} else {
//		InsertLog("200", "Successfully marshalled encrypted config", "WriteJSONConfig()")
//	}
//
//	// Write the encrypted JSON configuration to file
//	err = ioutil.WriteFile(configFile, data, 0644)
//	if err != nil {
//		InsertLog("200", "Failed to write encrypted config to file", "WriteJSONConfig()")
//		return err
//	} else {
//		InsertLog("200", "Successfully wrote encrypted config to file", "WriteJSONConfig()")
//	}
//
//	InsertLog("200", "Successfully wrote JSON config", "WriteJSONConfig()")
//	return nil
//}
//
//// AES encryption function
//func encryptAES(data []byte, key []byte) ([]byte, error) {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		InsertLog("200", "Failed to create new cipher", "encryptAES()")
//		return nil, err
//	} else {
//		InsertLog("200", "Successfully created new cipher", "encryptAES()")
//	}
//
//	blockSize := block.BlockSize()
//	paddedData := padPKCS7(data, blockSize)
//	InsertLog("200", "Successfully padded data", "encryptAES()")
//	ciphertext := make([]byte, len(paddedData))
//	mode := cipher.NewCBCEncrypter(block, key[:blockSize])
//	mode.CryptBlocks(ciphertext, paddedData)
//	InsertLog("200", "Successfully encrypted data", "encryptAES()")
//	return ciphertext, nil
//}
//
//// PKCS7 padding function
//func padPKCS7(data []byte, blockSize int) []byte {
//	padding := blockSize - (len(data) % blockSize)
//	paddedData := append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
//	InsertLog("200", "Successfully padded data", "padPKCS7()")
//	return paddedData
//}
//
//// AES decryption function
//func decryptAES(data []byte, key []byte) ([]byte, error) {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		InsertLog("400", "Failed to create new cipher", "decryptAES()")
//		WriteLog("encryptAES_fail_1", "200", "Failed to create new cipher", "decryptAES()", time.Now())
//		return nil, err
//	} else {
//		InsertLog("200", "Successfully created new cipher", "decryptAES()")
//		WriteLog("encryptAES_success_1", "200", "Successfully created new cipher", "decryptAES()", time.Now())
//	}
//
//	blockSize := block.BlockSize()
//	if len(data)%blockSize != 0 {
//		InsertLog("200", "Ciphertext length is not a multiple of the block size", "decryptAES()")
//		return nil, errors.New("ciphertext length is not a multiple of the block size")
//	} else {
//		InsertLog("200", "Ciphertext length is a multiple of the block size", "decryptAES()")
//	}
//	InsertLog("200", "Successfully checked ciphertext length", "decryptAES()")
//	mode := cipher.NewCBCDecrypter(block, key[:blockSize])
//	decryptedData := make([]byte, len(data))
//	mode.CryptBlocks(decryptedData, data)
//	InsertLog("200", "Successfully decrypted data", "decryptAES()")
//	// Remove padding
//	decryptedData = unpadPKCS7(decryptedData)
//	InsertLog("200", "Successfully removed padding", "decryptAES()")
//	return decryptedData, nil
//}
//
//// PKCS7 unpadding function
//func unpadPKCS7(data []byte) []byte {
//	padding := int(data[len(data)-1])
//	InsertLog("200", "Successfully got padding", "unpadPKCS7()")
//	return data[:len(data)-padding]
//}
