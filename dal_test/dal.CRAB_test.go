package dal

import (
	"cmpscfa23team2/dal"
	"reflect"
	"testing"
)

func TestCreateWebCrawler(t *testing.T) {
	sourceURL := "http://example.com"

	_, err := dal.CreateWebCrawler(sourceURL)
	if err != nil {
		t.Errorf("Unexpected error : %v", err)
	}
}

func TestCreateScraperEngine(t *testing.T) {
	engineName := "testengine"
	engineDescription := "test description"

	_, err := dal.CreateScraperEngine(engineName, engineDescription)
	if err != nil {
		t.Errorf("Couldn't make the engine: %v", err)

	}
}

func TestInsertURL(t *testing.T) {
	url := "http://example.com"
	domain := "example.com"
	tags := map[string]interface{}{"tag1": "value1", "tag2": "value2"}
	_, err := dal.InsertURL(url, domain, tags)
	if err != nil {
		t.Errorf("Couldn't insert url: %v", err)
	}
}

func TestUpdateURL(t *testing.T) {
	id := "123"
	url := "https://updatedurl.com"
	domain := "updated-example.com"
	tags := map[string]interface{}{"updated_tag1": "value1", "updated_tag2": "value2"}

	err := dal.UpdateURL(id, url, domain, tags)
	if err != nil {
		t.Errorf("Couldn't update url: %v", err)
	}

}

func TestGetURLTagsAndDomain(t *testing.T) {
	id := "5422f84d-7f6c-11ee-aa3b-6c2b59772aba"
	expectedTags := map[string]interface{}{"tag1": "value1", "tag2": "value2"}
	expectedDomain := "example.com"

	tags, domain, err := dal.GetURLTagsAndDomain(id)
	if err != nil {
		t.Errorf("Couldn't get tags and domain: %v", err)
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
		t.Errorf("Unexpected error: %v", err)

	}
	//if !reflect.DeepEqual(urls, expectedURLs) {
	//	t.Errorf("Expected urls: %v, got: %v", expectedURLs, urls)
	//}
}

func TestGetUUIDFromURLAndDomain(t *testing.T) {
	url := "http://example.com"
	domain := "example.com"
	//expectedUUID := "abc123"

	_, err := dal.GetUUIDFromURLAndDomain(url, domain)
	if err != nil {
		t.Errorf("Couldn't get UUID: %v", err)
	}
}

// function passes test if domain is valid, otherwise, if url contains null value for domain,
// this function does not pass the test
func TestGetRandomURL(t *testing.T) {
	_, err := dal.GetRandomURL()
	if err != nil {
		t.Errorf("Couldn't get random URL: %v", err)
	}
}
func TestGetURLsOnly(t *testing.T) {
	_, err := dal.GetURLsOnly()
	if err != nil {
		t.Errorf("Couldn't get URL: %v", err)
	}
}

func TestGetURLsAndTags(t *testing.T) {
	_, err := dal.GetURLsAndTags()
	if err != nil {
		t.Errorf("Couldn't get url and tag: %v", err)

	}
}
