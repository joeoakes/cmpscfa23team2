package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var db *sql.DB

type Log struct {
	LogID        string
	StatusCode   string
	Message      string
	GoEngineArea string
	DateTime     time.Time
}

type JSON_Data_Connect struct {
	Username string
	Password string
	Hostname string
	Database string
}

type Prediction struct {
	ID          int
	Model       string
	InputData   string
	Prediction  string
	PredictedAt time.Time
}

func init() {
	config, err := readJSONConfig("config.json")
	if err != nil {
		log.Fatal("Error reading JSON config:", err)
		return
	}

	var connErr error
	db, connErr = Connection(config)
	if connErr != nil {
		log.Fatal("Error establishing database connection:", connErr)
	}
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

func makePrediction(model string, inputData string) Prediction {
	// For demonstration, let's assume the prediction result is fixed.
	// Normally, this would involve running the ML model on inputData.
	prediction := "Some Prediction Result"
	currentTime := time.Now()

	return Prediction{
		Model:       model,
		InputData:   inputData,
		Prediction:  prediction,
		PredictedAt: currentTime,
	}
}

func storePrediction(pred Prediction) error {
	query := "INSERT INTO predictions (model, input_data, prediction, predicted_at) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, pred.Model, pred.InputData, pred.Prediction, pred.PredictedAt)
	return err
}

func main() {
	if db == nil {
		log.Fatal("Database connection is not initialized.")
	}

	// For multi-threading, let's use WaitGroup from the "sync" package
	var wg sync.WaitGroup

	models := []string{"model1", "model2", "model3"}
	inputData := "Some Input Data"

	for _, model := range models {
		wg.Add(1)

		go func(model string, inputData string) {
			defer wg.Done()

			prediction := makePrediction(model, inputData)
			err := storePrediction(prediction)
			if err != nil {
				log.Println("Error storing prediction:", err)
			}
		}(model, inputData)
	}

	wg.Wait()
}
