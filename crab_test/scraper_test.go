package crab_test

import (
	"cmpscfa23team2/crab"
	"sync"
	"testing"
)

//func TestReadCSV(t *testing.T) {
//	// Create a temporary file
//	file, err := ioutil.TempFile("", "test.csv")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer os.Remove(file.Name()) // clean up
//
//	// Write test data to the file
//	_, err = file.WriteString("Status,Bedrooms,Bathrooms,...\nFor Sale,3,2,...")
//	if err != nil {
//		t.Fatal(err)
//	}
//	file.Close()
//
//	// Test ReadCSV
//	properties, err := crab.ReadCSV(file.Name()) // Note the use of crab.
//	if err != nil {
//		t.Errorf("crab.ReadCSV() error = %v, wantErr %v", err, false)
//	}
//	if len(properties) != 1 {
//		t.Errorf("crab.ReadCSV() got %v items, want %v", len(properties), 1)
//	}
//
//}

func TestScrape(t *testing.T) {
	// This is a basic structure. You'll need to mock network calls for a proper test.
	// Alternatively, you can set up a local HTTP server that serves test HTML pages.

	domainConfig := crab.DomainConfig{ /* ... */ } // Note the use of crab.

	var wg sync.WaitGroup

	wg.Add(1)
	go crab.Scrape("https://books.toscrape.com/", domainConfig, &wg) // Note the use of crab.
	wg.Wait()

	// Assertions go here
	// For example, check if the file was created or the data format is as expected
}
