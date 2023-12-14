package dal_test

import (
	"cmpscfa23team2/dal"
	"reflect"
	"testing"
)

func TestCreateWebCrawler(t *testing.T) {
	sourceURL := "http://example.com"

	_, err := dal.CreateWebCrawler(sourceURL)
	if err != nil {
		dal.InsertLog("400", "Failed to create web crawler", "TestCreateWebCrawler()")
		t.Errorf("Unexpected error : %v", err)
	} else {
		dal.InsertLog("200", "Successfully created web crawler", "TestCreateWebCrawler()")
	}
}

func TestCreateScraperEngine(t *testing.T) {
	engineName := "testengine"
	engineDescription := "test description"

	_, err := dal.CreateScraperEngine(engineName, engineDescription)
	if err != nil {
		dal.InsertLog("400", "Failed to create scraper engine", "TestCreateScraperEngine()")
		t.Errorf("Couldn't make the engine: %v", err)
	} else {
		dal.InsertLog("200", "Successfully created scraper engine", "TestCreateScraperEngine()")
	}
}

func TestInsertURL(t *testing.T) {
	url := "http://example.com"
	domain := "example.com"
	tags := map[string]interface{}{"tag1": "value1", "tag2": "value2"}
	_, err := dal.InsertURL(url, domain, tags)
	if err != nil {
		dal.InsertLog("400", "Failed to insert URL", "TestInsertURL()")
		t.Errorf("Couldn't insert URL: %v", err)
	} else {
		dal.InsertLog("200", "Successfully inserted URL", "TestInsertURL()")
	}
}

func TestUpdateURL(t *testing.T) {
	id := "123"
	url := "https://updatedurl.com"
	domain := "updated-example.com"
	tags := map[string]interface{}{"updated_tag1": "value1", "updated_tag2": "value2"}

	err := dal.UpdateURL(id, url, domain, tags)
	if err != nil {
		dal.InsertLog("400", "Failed to update URL", "TestUpdateURL()")
		t.Errorf("Couldn't update URL: %v", err)
	} else {
		dal.InsertLog("200", "Successfully updated URL", "TestUpdateURL()")
	}
}

func TestGetURLTagsAndDomain(t *testing.T) {
	id := "20303a5b-8ff4-11ee-ae02-30d042e80ac3"
	expectedTags := map[string]interface{}{"tag1": "value1", "tag2": "value2"}
	expectedDomain := "example.com"

	tags, domain, err := dal.GetURLTagsAndDomain(id)
	if err != nil {
		dal.InsertLog("400", "Failed to get URL tags and domain", "TestGetURLTagsAndDomain()")
		t.Errorf("Couldn't get tags and domain: %v", err)
	} else {
		dal.InsertLog("200", "Successfully got URL tags and domain", "TestGetURLTagsAndDomain()")
	}

	if !reflect.DeepEqual(tags, expectedTags) {
		t.Errorf("Expected tags: %v, got: %v", expectedTags, tags)
	}
	if domain != expectedDomain {
		t.Errorf("Expected domain: %s, got: %s", expectedDomain, domain)
	}
}

func TestGetURLsFromDomain(t *testing.T) {
	domain := "example.com"
	//expectedURLs := []string{"http://example.com/page1", "http://example.com/page2"}
	_, err := dal.GetURLsFromDomain(domain)
	if err != nil {
		dal.InsertLog("400", "Failed to get URLs from domain", "TestGetURLsFromDomain()")
		t.Errorf("Unexpected error: %v", err)
	} else {
		dal.InsertLog("200", "Successfully got URLs from domain", "TestGetURLsFromDomain()")
	}
	//if !reflect.DeepEqual(urls, expectedURLs) {
	//	t.Errorf("Expected urls: %v, got: %v", expectedURLs, urls)
	//}
}
