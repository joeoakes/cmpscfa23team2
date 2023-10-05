package main

import (
	"encoding/json"
	_ "errors"
	_ "github.com/go-sql-driver/mysql"
)

// Function to create a new web crawler
func CreateWebCrawler(sourceURL string) (string, error) {
	var crawlerID string
	err := db.QueryRow("CALL create_webcrawler(?)", sourceURL).Scan(&crawlerID)
	if err != nil {
		return "", err
	}
	return crawlerID, nil
}

// Function to create a new scraper engine
func CreateScraperEngine(engineName, engineDescription string) (string, error) {
	var engineID string
	err := db.QueryRow("CALL create_scraper_engine(?, ?)", engineName, engineDescription).Scan(&engineID)
	if err != nil {
		return "", err
	}
	return engineID, nil
}

// Function to insert a new URL
func InsertURL(url, domain string, tags map[string]interface{}) (string, error) {
	var id string
	jsonTags, err := json.Marshal(tags)
	if err != nil {
		return "", err
	}

	err = db.QueryRow("CALL insert_url(?, ?, ?)", url, string(jsonTags), domain).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Function to update an existing URL
func UpdateURL(id, url, domain string, tags map[string]interface{}) error {
	jsonTags, err := json.Marshal(tags)
	if err != nil {
		return err
	}

	_, err = db.Exec("CALL update_url(?, ?, ?, ?)", id, url, string(jsonTags), domain)
	return err
}

// Function to fetch URL tags and domain by ID
func GetURLTagsAndDomain(id string) (map[string]interface{}, string, error) {
	var tagsStr, domain string
	err := db.QueryRow("CALL get_url_tags_and_domain(?)", id).Scan(&tagsStr, &domain)
	if err != nil {
		return nil, "", err
	}

	var tags map[string]interface{}
	err = json.Unmarshal([]byte(tagsStr), &tags)
	if err != nil {
		return nil, "", err
	}

	return tags, domain, nil
}

// Function to fetch URLs from a specific domain
func GetURLsFromDomain(domain string) ([]string, error) {
	rows, err := db.Query("CALL get_urls_from_domain(?)", domain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var id, url, tags, domain string
		var createdTime []byte // <-- Change this line

		if err := rows.Scan(&id, &url, &tags, &domain, &createdTime); err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	return urls, rows.Err()
}

// Function to fetch UUID from URL and domain
func GetUUIDFromURLAndDomain(url, domain string) (string, error) {
	var id string
	err := db.QueryRow("CALL get_Uuid_from_URL_and_domain(?, ?)", url, domain).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// Function to fetch a random URL
func GetRandomURL() (string, error) {
	var id, url, tags, domain string
	var createdTime []byte // As per our earlier fix
	err := db.QueryRow("CALL get_random_url()").Scan(&id, &url, &tags, &domain, &createdTime)
	if err != nil {
		return "", err
	}
	return url, nil
}

// Function to fetch all URLs (just the 'url' column)
func GetURLsOnly() ([]string, error) {
	rows, err := db.Query("CALL get_urls_only()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, rows.Err()
}

// Function to fetch all URLs with their tags
func GetURLsAndTags() (map[string]map[string]interface{}, error) {
	rows, err := db.Query("CALL get_urls_and_tags()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urlsAndTags := make(map[string]map[string]interface{})
	for rows.Next() {
		var url string
		var tagsStr string
		if err := rows.Scan(&url, &tagsStr); err != nil {
			return nil, err
		}

		var tags map[string]interface{}
		err = json.Unmarshal([]byte(tagsStr), &tags)
		if err != nil {
			return nil, err
		}

		urlsAndTags[url] = tags
	}

	return urlsAndTags, rows.Err()
}
