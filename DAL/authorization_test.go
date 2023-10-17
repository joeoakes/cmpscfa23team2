package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var testDb *sql.DB

type Config struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Hostname string `json:"Hostname"`
	Database string `json:"Database"`
}

func readConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func init() {
	config, err := readConfig("config.json")
	if err != nil {
		panic(err)
	}

	dbConnStr := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.Hostname, config.Database)
	testDb, err = sql.Open("mysql", dbConnStr)
	if err != nil {
		panic(err)
	}
}
