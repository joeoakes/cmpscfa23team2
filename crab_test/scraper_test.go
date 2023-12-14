package crab_test

import (
	"cmpscfa23team2/crab"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

func TestReadCSV(t *testing.T) {
	// Create a temporary file
	file, err := ioutil.TempFile("", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name()) // clean up

	// Write test data to the file
	_, err = file.WriteString("Status,Bedrooms,Bathrooms,...\nFor Sale,3,2,...")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	// Test ReadCSV
	properties, err := main.ReadCSV(file.Name())
	if err != nil {
		t.Errorf("ReadCSV() error = %v, wantErr %v", err, false)
	}
	if len(properties) != 1 {
		t.Errorf("ReadCSV() got %v items, want %v", len(properties), 1)
	}

	// Additional checks can be added to validate the contents of 'properties'
}

func TestScrape(t *testing.T) {
	// This is a basic structure. You'll need to mock network calls for a proper test.
	// Alternatively, you can set up a local HTTP server that serves test HTML pages.

	// Set up a test DomainConfig and a WaitGroup
	domainConfig := main.DomainConfig{ /* ... */ }
	var wg sync.WaitGroup

	wg.Add(1)
	go main.Scrape("http://example.com", domainConfig, &wg)
	wg.Wait()

	// Assertions go here
	// For example, check if the file was created or the data format is as expected
}
