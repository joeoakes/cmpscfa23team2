package crab_test

import (
	"cmpscfa23team2/crab"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGetURLsToCrawl(t *testing.T) {
	expected := []crab.URLData{
		{URL: "https://books.toscrape.com/"},
	}
	result := crab.GetURLsToCrawl() // Change from main.GetURLsToCrawl to crab.GetURLsToCrawl

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("getURLsToCrawl() = %v, want %v", result, expected)
	}
}

func TestInsertData(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	data := crab.ItemData{ /* initialize with test data */ }
	err = crab.InsertData(data, tmpfile.Name())

	if err != nil {
		t.Errorf("InsertData() error = %v, wantErr %v", err, false)
	}

	// Optionally, read back the file and check contents
}
func TestCreateSiteMap(t *testing.T) {
	urls := []crab.URLData{ /* initialize with test data */ }
	tmpfile, err := ioutil.TempFile("", "sitemap")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	err = crab.CreateSiteMap(urls)
	if err != nil {
		t.Errorf("createSiteMap() error = %v", err)
	}

	// Optionally, read back the file and check contents
}

func TestIsURLAllowedByRobotsTXT(t *testing.T) {
	// Example URL
	url := "https://books.toscrape.com/"

	allowed := crab.IsURLAllowedByRobotsTXT(url)
	if !allowed {
		t.Errorf("isURLAllowedByRobotsTXT(%s) = %v, want %v", url, allowed, true)
	}
}
