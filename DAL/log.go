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

type JsonDataConnect struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Database string `json:"database"`
}

type Log struct {
	LogID        string
	StatusCode   string
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
		return err
	}

	encryptedPassword, err := encryptAES([]byte(config.Password), key)
	if err != nil {
		return err
	}

	encryptedConfig := JsonDataConnect{
		Username: base64.StdEncoding.EncodeToString(encryptedUsername),
		Password: base64.StdEncoding.EncodeToString(encryptedPassword),
		Hostname: config.Hostname,
		Database: config.Database,
	}

	// Marshal the encrypted configuration struct to JSON
	data, err := json.Marshal(encryptedConfig)
	if err != nil {
		return err
	}

	// Write the encrypted JSON configuration to file
	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// AES encryption function
func encryptAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	paddedData := padPKCS7(data, blockSize)

	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, key[:blockSize])
	mode.CryptBlocks(ciphertext, paddedData)

	return ciphertext, nil
}

// PKCS7 padding function
func padPKCS7(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	paddedData := append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
	return paddedData
}

// Function to read the JSON configuration file
func ReadJSONConfig(configFile string, key []byte) (*JsonDataConnect, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var encryptedConfig JsonDataConnect
	err = json.Unmarshal(data, &encryptedConfig)
	if err != nil {
		return nil, err
	}

	decryptedUsername, err := base64.StdEncoding.DecodeString(encryptedConfig.Username)
	if err != nil {
		return nil, err
	}

	decryptedPassword, err := base64.StdEncoding.DecodeString(encryptedConfig.Password)
	if err != nil {
		return nil, err
	}

	username, err := decryptAES(decryptedUsername, key)
	if err != nil {
		return nil, err
	}

	password, err := decryptAES(decryptedPassword, key)
	if err != nil {
		return nil, err
	}

	decryptedConfig := JsonDataConnect{
		Username: string(username),
		Password: string(password),
		Hostname: encryptedConfig.Hostname,
		Database: encryptedConfig.Database,
	}

	return &decryptedConfig, nil
}

// AES decryption function
func decryptAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(data)%blockSize != 0 {
		return nil, errors.New("ciphertext length is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decryptedData := make([]byte, len(data))
	mode.CryptBlocks(decryptedData, data)

	// Remove padding
	decryptedData = unpadPKCS7(decryptedData)

	return decryptedData, nil
}

// PKCS7 unpadding function
func unpadPKCS7(data []byte) []byte {
	padding := int(data[len(data)-1])
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
		log.Fatal(err)
	}

	log.SetOutput(file)
}

// readJSONConfig - Reads the JSON config file
func readJSONConfig(filename string) (JsonDataConnect, error) {
	var config JsonDataConnect
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

// WriteLog writes a log entry to the database
func WriteLog(logID string, statusCode string, message string, goEngineArea string, dateTime time.Time) error {
	config, err := ReadJSONConfig("config.json", []byte("IST440WSRA440WGE"))
	if err != nil {
		return err
	}

	// Create the connection string using the configuration values
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.Username, config.Password, config.Hostname, config.Database)

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	defer func() {
		// Close the database connection when finished
		closeErr := db.Close()
		if closeErr != nil {
			log.Println("Error closing database connection:", closeErr)
		}
	}()

	stmt, err := db.Prepare("INSERT INTO log(log_ID, status_code, message, go_engine_area, date_Time) VALUES (? ,? ,? ,? ,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, errExec := stmt.Exec(logID, statusCode, message, goEngineArea, dateTime)
	if errExec != nil {
		return errExec
	}

	return nil
}

// GetLog - Reads the log
func GetLog() ([]Log, error) {
	config, err := readJSONConfig("config.json")
	if err != nil {
		return nil, err
	}

	db, err := Connection(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("CALL select_all_logs()")
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

		// Assign the dateTimeStr directly to the DateTime field
		logItem.DateTime = dateTimeStr

		logs = append(logs, logItem)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetSuccess - Uses a Procedure to gather all the 'Success' rows in the DB
func GetSuccess() ([]Log, error) {
	config, err := readJSONConfig("config.json")
	if err != nil {
		return nil, err
	}
	db, err := Connection(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("CALL select_all_logs_by_status_code(?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Modify the argument to the desired status code
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

		// Assign the dateTimeStr directly to the DateTime field
		logItem.DateTime = dateTimeStr

		logs = append(logs, logItem)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func StoreLog(statusCode string, message string, goEngineArea string) error {
	config, err := readJSONConfig("config.json")
	if err != nil {
		return err
	}

	db, err := Connection(config)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("CALL insert_log(?, ?, ?)")
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

// Connection - Establishes connection to the database
func Connection(config JsonDataConnect) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.Hostname, config.Database))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
