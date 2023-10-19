package DAL

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var db *sql.DB

type JsonDataConnect struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Database string `json:"database"`
}

// Log struct models the data structure of a log entry in the database
type Log struct {
	LogID        string
	status_code  string
	Message      string
	GoEngineArea string
	DateTime     []uint8
}

type LogStatusCodes struct {
	StatusCode    string
	StatusMessage string
}

// URL Struct to hold the URL details
type URL struct {
	ID          []uint8 `json:"id"`
	Tags        string  `json:"tags"`
	CreatedTime string  `json:"created_time"`
	URL         string  `json:"url"`
}

type SqlDatabase struct {
	db *sql.DB
}

type TagData struct {
	URL  string            `json:"url"`
	Tags map[string]string `json:"tags"`
}

// GetTagsString Helper function to get a formatted string representation of tags
func GetTagsString(tagData TagData) string {
	tagsString := ""
	for tag, attributes := range tagData.Tags {
		tagsString += fmt.Sprintf("%s=%s ", tag, attributes)
	}
	return tagsString
}

func GetURLs(d *SqlDatabase) ([]URL, error) {
	// Query the URLs from the database
	rows, err := d.db.Query("CALL get_urls_only()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []URL

	// Iterate over the rows
	for rows.Next() {
		var url URL
		err := rows.Scan(&url.ID, &url.URL, &url.Tags)
		if err != nil {
			log.Println(err)
			continue
		}

		urls = append(urls, url)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

// Function to write the JSON configuration file with encrypted username and password
func WriteJSONConfig(configFile string, config *JsonDataConnect, key []byte) error {
	encryptedUsername, err := encryptAES([]byte(config.Username), key)
	if err != nil {
		InsertLog(db, "200", "Failed to encrypt username", "WriteJSONConfig()")
		return err
	} else {InsertLog(db, "200", "Successfully encrypted username", "WriteJSONConfig()")}

	encryptedPassword, err := encryptAES([]byte(config.Password), key)
	if err != nil {
		InsertLog(db, "200", "Failed to encrypt password", "WriteJSONConfig()")
		return err
	} else {InsertLog(db, "200", "Successfully encrypted password", "WriteJSONConfig()")}

	encryptedConfig := JsonDataConnect{
		Username: base64.StdEncoding.EncodeToString(encryptedUsername),
		Password: base64.StdEncoding.EncodeToString(encryptedPassword),
		Hostname: config.Hostname,
		Database: config.Database,
	}
	InsertLog(db, "200", "Successfully encrypted config", "WriteJSONConfig()"
	// Marshal the encrypted configuration struct to JSON
	data, err := json.Marshal(encryptedConfig)
	if err != nil {
		InsertLog(db, "200", "Failed to marshal encrypted config", "WriteJSONConfig()"
		return err
	} else {InsertLog(db, "200", "Successfully marshalled encrypted config", "WriteJSONConfig()")}

	// Write the encrypted JSON configuration to file
	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		InsertLog(db, "200", "Failed to write encrypted config to file", "WriteJSONConfig()"
		return err
	} else {InsertLog(db, "200", "Successfully wrote encrypted config to file", "WriteJSONConfig()")}

	InsertLog(db, "200", "Successfully wrote JSON config", "WriteJSONConfig()")
	return nil
}

// AES encryption function
func encryptAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		InsertLog(db, "200", "Failed to create new cipher", "encryptAES()")
		return nil, err
	} else {InsertLog(db, "200", "Successfully created new cipher", "encryptAES()")}

	blockSize := block.BlockSize()
	paddedData := padPKCS7(data, blockSize)
	InsertLog(db, "200", "Successfully padded data", "encryptAES()"
	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, key[:blockSize])
	mode.CryptBlocks(ciphertext, paddedData)
	InsertLog(db, "200", "Successfully encrypted data", "encryptAES()"
	return ciphertext, nil
}

// PKCS7 padding function
func padPKCS7(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	paddedData := append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
	InsertLog(db, "200", "Successfully padded data", "padPKCS7()")
	return paddedData
}

// Function to read the JSON configuration file
func ReadJSONConfig(configFile string, key []byte) (*JsonDataConnect, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		InsertLog(db, "200", "Failed to read JSON config file", "ReadJSONConfig()"
		return nil, err
	} else{InsertLog(db, "200", "Successfully read JSON config file", "ReadJSONConfig()")}

	var encryptedConfig JsonDataConnect
	err = json.Unmarshal(data, &encryptedConfig)
	if err != nil {
		InsertLog(db, "200", "Failed to unmarshal JSON config", "ReadJSONConfig()")
		return nil, err
	} else{InsertLog(db, "200", "Successfully unmarshalled JSON config", "ReadJSONConfig()")}

	decryptedUsername, err := base64.StdEncoding.DecodeString(encryptedConfig.Username)
	if err != nil {
		InsertLog(db, "200", "Failed to decode username", "ReadJSONConfig()")
		return nil, err
	} else{InsertLog(db, "200", "Successfully decoded username", "ReadJSONConfig()")}

	decryptedPassword, err := base64.StdEncoding.DecodeString(encryptedConfig.Password)
	if err != nil {
		InsertLog(db, "200", "Failed to decode password", "ReadJSONConfig()")
		return nil, err
	} else{InsertLog(db, "200", "Successfully decoded password", "ReadJSONConfig()")}

	username, err := decryptAES(decryptedUsername, key)
	if err != nil {
		return nil, err
	} else{InsertLog(db, "200", "Successfully decrypted username", "ReadJSONConfig()")}

	password, err := decryptAES(decryptedPassword, key)
	if err != nil {
		InsertLog(db, "200", "Failed to decrypt password", "ReadJSONConfig()")
		return nil, err
	} else{InsertLog(db, "200", "Successfully decrypted password", "ReadJSONConfig()")}


	decryptedConfig := JsonDataConnect{
		Username: string(username),
		Password: string(password),
		Hostname: encryptedConfig.Hostname,
		Database: encryptedConfig.Database,
	}
	InsertLog(db, "200", "Successfully decrypted config", "ReadJSONConfig()")

	return &decryptedConfig, nil
}

// AES decryption function
func decryptAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		InsertLog(db, "200", "Failed to create new cipher", "decryptAES()")
		return nil, err
	} else{InsertLog(db, "200", "Successfully created new cipher", "decryptAES()")}

	blockSize := block.BlockSize()
	if len(data)%blockSize != 0 {
		InsertLog(db, "200", "Ciphertext length is not a multiple of the block size", "decryptAES()"
		return nil, errors.New("ciphertext length is not a multiple of the block size")
	} else{InsertLog(db, "200", "Ciphertext length is a multiple of the block size", "decryptAES()")}
	InsertLog(db, "200", "Successfully checked ciphertext length", "decryptAES()"
	mode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decryptedData := make([]byte, len(data))
	mode.CryptBlocks(decryptedData, data)
	InsertLog(db, "200", "Successfully decrypted data", "decryptAES()"
	// Remove padding
	decryptedData = unpadPKCS7(decryptedData)
	InsertLog(db, "200", "Successfully removed padding", "decryptAES()"
	return decryptedData, nil
}

// PKCS7 unpadding function
func unpadPKCS7(data []byte) []byte {
	padding := int(data[len(data)-1])
	InsertLog(db, "200", "Successfully got padding", "unpadPKCS7()"
	return data[:len(data)-padding]
}

// Function to insert a log entry into the database
func InsertLog(db *sql.DB, statusCode, message, goEngineArea string) error {
	// logID := uuid.New().String()

	_, err := db.Exec("CALL insert_log(?, ?, ?)", statusCode, message, goEngineArea)
	if err != nil {
		log.Println("Error inserting log:", err)
	}
	return err
}

// This function creates & adds the log entries to a TextFile if the database is down
func init() {
	file, err := os.OpenFile("CRAB_Logging_to.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		InsertLog(db, "200", "Failed to open file", "init()"
		log.Fatal(err)
	} else{
		InsertLog(db, "200", "INIT Open File Success", "init()")}

	log.SetOutput(file)
}

// readJSONConfig - Reads the JSON config file
func readJSONConfig(filename string) (JsonDataConnect, error) {
	var config JsonDataConnect
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		InsertLog(db, "200", "Failed to read JSON config file", "readJSONConfig()"
		return config, err
	} else{
		InsertLog(db, "200", "Successfully read JSON config file", "readJSONConfig()")
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		InsertLog(db, "200", "Failed to unmarshal JSON config", "readJSONConfig()"
		return config, err
	} else{InsertLog(db, "200", "Successfully unmarshalled JSON config", "readJSONConfig()")}

	return config, nil
}

// WriteLog writes a log entry to the database
func WriteLog(logID string, status_code string, message string, goEngineArea string, dateTime time.Time) error {
	// Validate the statusCode by checking if it exists in the `log_status_codes` table
	var existingStatusCode string
	err := db.QueryRow("SELECT status_code FROM log_status_codes WHERE status_code = ?", status_code).Scan(&existingStatusCode)
	if err != nil {
		InsertLog(db, "200", "Failed to query row", "WriteLog()")
		if err == sql.ErrNoRows {
			InsertLog(db, "200", "Invalid statusCode", "WriteLog()")
			return fmt.Errorf("Invalid statusCode: %s", status_code)
		} else {
			InsertLog(db, "200", "Successfully validated status code", "WriteLog()")
		}

		InsertLog(db, "200", "Successfully validated status code", "WriteLog()")
		return err
	}
	// Prepare the SQL statement for inserting into the log table
	stmt, err := db.Prepare("INSERT INTO log(log_ID, status_code, message, go_engine_area, date_time) VALUES (? ,? ,? ,? ,?)")
	if err != nil {
		InsertLog(db, "200", "Failed to prepare SQL statement", "WriteLog()")
		return err
	} else{InsertLog(db, "200", "Successfully prepared SQL statement", "WriteLog()")}

	defer stmt.Close()

	// Execute the SQL statement
	_, errExec := stmt.Exec(logID, existingStatusCode, message, goEngineArea, dateTime)
	if errExec != nil {
		InsertLog(db, "200", "Failed to execute SQL statement", "WriteLog()"
		return errExec
	} else{InsertLog(db, "200", "Successfully executed SQL statement", "WriteLog()")}

	return nil
}

// GetLog - Reads the log
func GetLog() ([]Log, error) {
	stmt, err := db.Prepare("CALL select_all_logs()")
	if err != nil {
		InsertLog(db, "200", "Failed to prepare SQL statement", "GetLog()"
		return nil, err
	} else{InsertLog(db, "200", "Successfully prepared SQL statement", "GetLog()")}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		InsertLog(db, "200", "Failed to query SQL statement", "GetLog()"
		return nil, err
	} else{log.Println("Successfully queried SQL statement")}
	defer rows.Close()

	var logs []Log
	for rows.Next() {
		var logItem Log
		var dateTimeStr []uint8
		err := rows.Scan(&logItem.LogID, &logItem.status_code, &logItem.Message, &logItem.GoEngineArea, &dateTimeStr)
		if err != nil {
			return nil, err
		} else{InsertLog(db, "200", "Successfully scanned rows", "GetLog()")}

		logItem.DateTime = dateTimeStr
		logs = append(logs, logItem)
		InsertLog(db, "200", "Successfully appended log item", "GetLog()")
	}

	if err = rows.Err(); err != nil {
		InsertLog(db, "200", "Failed to iterate over rows", "GetLog()"
		return nil, err
	} else{InsertLog(db, "200", "Successfully iterated over rows", "GetLog()")}

	return logs, nil
}

// GetSuccess - Uses a Procedure to gather all the 'Success' rows in the DB
func GetSuccess() ([]Log, error) {
	stmt, err := db.Prepare("CALL select_all_logs_by_status_code(?)")
	if err != nil {
		InsertLog(db, "200", "Failed to prepare SQL statement", "GetSuccess()"
		return nil, err
	} else{InsertLog(db, "200", "Successfully prepared SQL statement", "GetSuccess()")}
	defer stmt.Close()

	rows, err := stmt.Query("Success")
	if err != nil {
		return nil, err
	} else{InsertLog(db, "200", "Successfully queried SQL statement", "GetSuccess()")}
	defer rows.Close()
	InsertLog(db, "200", "Successfully closed rows", "GetSuccess()"
	var logs []Log
	for rows.Next() {
		var logItem Log
		var dateTimeStr []uint8
		err := rows.Scan(&logItem.LogID, &logItem.status_code, &logItem.Message, &logItem.GoEngineArea, &dateTimeStr)
		if err != nil {
			return nil, err
		} else{InsertLog(db, "200", "Successfully scanned rows", "GetSuccess()")}
		logItem.DateTime = dateTimeStr
		logs = append(logs, logItem)
		InsertLog(db, "200", "Successfully appended log item", "GetSuccess()")
	}

	if err = rows.Err(); err != nil {
		return nil, err
	} else{InsertLog(db, "200", "Successfully iterated over rows", "GetSuccess()")}

	return logs, nil
}

func StoreLog(status_code string, message string, goEngineArea string) error {
	stmt, err := db.Prepare("CALL insert_log(?,?,?)")
	if err != nil {
		return err
	} else{InsertLog(db, "200", "Successfully prepared SQL statement", "StoreLog()")}
	defer stmt.Close()

	_, errExec := stmt.Exec(status_code, message, goEngineArea)
	if errExec != nil {
		return errExec
	} else{InsertLog(db, "200", "Successfully executed SQL statement", "StoreLog()")}

	return nil
}

// Connection - Establishes connection to the database
func Connection(config JsonDataConnect) (*sql.DB, error) {
	connDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.Hostname, config.Database))
	if err != nil {
		return nil, err
	} else{InsertLog(db, "200", "Successfully opened database connection", "Connection()")}

	err = connDB.Ping()
	if err != nil {
		return nil, err
	} else{InsertLog(db, "200", "Successfully pinged database", "Connection()"))}

	log.Printf("Database connection: %+v", connDB)
	return connDB, nil
}

// Added the main function where you need to test all the methods above, some of the testing structure is below
// Note: Don't worry about some that are missing you can comment them out
// Or, you can find those methods in the CRUD/setup.go
// main function to test all existing methods
func main() {

	// Initialize database connection
	if db == nil {
		log.Fatal("Database connection is not initialized.")
	} else{
		InsertLog(db, "200", "Database connection is initialized", "main()")
	}

	// Insert or Update status code
	err := InsertOrUpdateStatusCode("POS", "noth")
	if err != nil {
		log.Println("Failed to insert or update status code:", err)
		InsertLog(db, "200", "Failed to insert or update status code", "main()")
	} else{
		InsertLog(db, "200", "Successfully inserted or updated status code", "main()")
	}

	_, err = FetchUserID("jxo19")
	if err != nil {
		log.Fatalf("Failed to fetch user ID: %v", err)
		InsertLog(db, "200", "Failed to fetch user ID", "main()")
	} else{
		InsertLog(db, "200", "Successfully fetched user ID", "main()")
	}

	// Update User
	err = UpdateUser("NewName", "jxo19", "ADM", "newpassword")
	if err != nil {
		fmt.Printf("Failed to update user: %s\n", err)
		InsertLog(db, "200", "Failed to update user", "main()"
	} else {
		InsertLog(db, "200", "Successfully updated user", "main()"
	}

	// Delete User
	err = DeleteUser("jxo19")
	if err != nil {
		fmt.Printf("Failed to delete user: %s\n", err)
		InsertLog(db, "200", "Failed to delete user", "main()"
	} else {
		InsertLog(db, "200", "Successfully deleted user", "main()"
	}

	//Generate a unique logID
	uniqueLogID := uuid.New().String()

	//Write log
	currentTime := time.Now()
	err = WriteLog(uniqueLogID, "Pos", "Message logged successfully", "Engine1", currentTime)
	if err != nil {
		log.Println("Failed to write log:", err)
		InsertLog(db, "200", "Failed to write log", "main()")
	} else{InsertLog(db, "200", "Successfully wrote log", "main()")}

	// Get and print all logs
	logs, err := GetLog()
	if err != nil {
		log.Println("Failed to get logs:", err)
		InsertLog(db, "200", "Failed to get logs", "main()")
	} else {
		for _, logItem := range logs {
			fmt.Println(logItem)
			err := InsertLog(db, "200", "Successfully got logs", "main()")
			if err != nil {
				return
			}
		}
	}

	//Store log using a stored procedure (uncomment if needed)
	err = StoreLog("Success", "Stored using procedure", "Engine1")
	if err != nil {
		log.Println("Failed to store log using stored procedure:", err)
		InsertLog(db, "200", "Failed to store log using stored procedure", "main()"
	} else{InsertLog(db, "200", "Successfully stored log using stored procedure", "main()")}

	//Insert a new status code
	err = InsertStatusCode("200", "OK")
	if err != nil {
		log.Println("Failed to insert new status code:", err)
		InsertLog(db, "200", "Failed to insert new status code", "main()")
	} else{InsertLog(db, "200", "Successfully inserted new status code", "main()")}

	//Create a new user
	err = CreateUser("John", "john123", "ADM", "password", true)
	if err != nil {
		log.Println("Failed to create a new user:", err)
		InsertLog(db, "200", "Failed to create a new user", "main()")
	} else{InsertLog(db, "200", "Successfully created a new user", "main()")}

	//Delete a user
	//err = DeleteUser("john123")
	//if err != nil {
	//	log.Println("Failed to delete user:", err)
	//}
	//
	////Get and print all "Success" logs
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
