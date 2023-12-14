package crab_test

import (
	"cmpscfa23team2/dal"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup: Initialize the database
	err := dal.InitDB()
	if err != nil {
		panic("Failed to initialize the database: " + err.Error())
	}

	// Run all tests in the package
	code := m.Run()

	// Teardown: Close the database
	dal.CloseDb()

	os.Exit(code)
}
