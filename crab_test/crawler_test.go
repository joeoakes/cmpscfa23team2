package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGetURLsToCrawl(t *testing.T) {
	expected := []URLData{
		{URL: "https://www.kaggle.com/search?q=housing+prices"},
	}
	result := GetURLsToCrawl()

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

	data := ItemData{ /* initialize with test data */ }
	err = InsertData(data, tmpfile.Name())

	if err != nil {
		t.Errorf("InsertData() error = %v, wantErr %v", err, false)
	}

	// Optionally, read back the file and check contents
}
func TestCreateSiteMap(t *testing.T) {
	urls := []URLData{ /* initialize with test data */ }
	tmpfile, err := ioutil.TempFile("", "sitemap")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	err = CreateSiteMap(urls)
	if err != nil {
		t.Errorf("createSiteMap() error = %v", err)
	}

	// Optionally, read back the file and check contents
}

func TestIsURLAllowedByRobotsTXT(t *testing.T) {
	// Example URL
	url := "http://example.com"

	allowed := IsURLAllowedByRobotsTXT(url)
	if !allowed {
		t.Errorf("isURLAllowedByRobotsTXT(%s) = %v, want %v", url, allowed, true)
	}
}
