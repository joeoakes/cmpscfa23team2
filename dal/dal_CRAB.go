package main

import (
	"encoding/json"
	_ "errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// Function to create a new web crawler
//
// It creates a web crawler with a specified source URL and logs the crawler's ID if successful.
func CreateWebCrawler(sourceURL string) (string, error) {
	var crawlerID string
	err := DB.QueryRow("CALL create_webcrawler(?)", sourceURL).Scan(&crawlerID)
	if err != nil {
		return "", err
	} else {
		log.Printf("Web crawler created: %s", crawlerID)
	}
	return crawlerID, nil
}

// Function to create a new scraper engine
//
// defines a function called "CreateScraperEngine" that creates a scraper engine in a database, and it returns the engine's ID or an error.
func CreateScraperEngine(engineName, engineDescription string) (string, error) {
	var engineID string
	err := DB.QueryRow("CALL create_scraper_engine(?, ?)", engineName, engineDescription).Scan(&engineID)
	if err != nil {
		return "", err
	} else {
		log.Printf("Scraper engine created: %s", engineID)
	}
	return engineID, nil
}

// Function to insert a new URL
//
// Function "InsertURL," inserts a URL into a database along with associated tags and logs the operation, returning the generated ID or an error.
func InsertURL(url, domain string, tags map[string]interface{}) (string, error) {
	var id string
	jsonTags, err := json.Marshal(tags)
	if err != nil {
		return "", err
	} else {
		log.Printf("URL inserted with tags: %v", tags)
	}

	err = DB.QueryRow("CALL insert_url(?, ?, ?)", url, string(jsonTags), domain).Scan(&id)
	if err != nil {
		return "", err
	} else {
		log.Printf("URL inserted with tags: %v", tags)
	}
	return id, nil
}

// Function to update an existing URL
//
// It defines a function UpdateURL that updates a URL record in a database, converting tags into JSON format and logging the update action.
func UpdateURL(id, url, domain string, tags map[string]interface{}) error {
	jsonTags, err := json.Marshal(tags)
	if err != nil {
		return err
	} else {
		log.Printf("URL updated with tags: %v", tags)
	}

	_, err = DB.Exec("CALL update_url(?, ?, ?, ?)", id, url, string(jsonTags), domain)
	return err
}

// Function to fetch URL tags and domain by ID
//
// It defines a function that retrieves tags and a domain from a database using a specified ID, logs the results, and returns them in a map and a string along with potential errors.
func GetURLTagsAndDomain(id string) (map[string]interface{}, string, error) {
	var tagsStr, domain string
	err := DB.QueryRow("CALL get_url_tags_and_domain(?)", id).Scan(&tagsStr, &domain)
	if err != nil {
		return nil, "", err
	} else {
		log.Printf("Tags: %v, Domain: %s", tagsStr, domain)
	}
	var tags map[string]interface{}
	err = json.Unmarshal([]byte(tagsStr), &tags)
	if err != nil {
		return nil, "", err
	} else {
		log.Printf("Tags: %v, Domain: %s", tags, domain)
	}

	return tags, domain, nil
}

// Function to fetch URLs from a specific domain
//
// Defines a function that queries a database to retrieve URLs associated with a given domain, processes the results, and returns the URLs in a slice while handling potential errors and logging.
func GetURLsFromDomain(domain string) ([]string, error) {
	rows, err := DB.Query("CALL get_urls_from_domain(?)", domain)
	if err != nil {
		return nil, err
	} else {
		log.Printf("URLs from domain: %v", rows)
	}
	log.Println("Closing Rows: %+v", rows)
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var id, url, tags, domain string
		var createdTime []byte // <-- Change this line
		if err := rows.Scan(&id, &url, &tags, &domain, &createdTime); err != nil {
			return nil, err
		} else {
			log.Printf("URLs from domain: %v", urls)
		}
		urls = append(urls, url)
	}
	return urls, rows.Err()
}

// Function to fetch UUID from URL and domain
//
// This function retrieves a UUID from a database by calling a stored procedure with a given URL and domain, and it logs the result if successful.
func GetUUIDFromURLAndDomain(url, domain string) (string, error) {
	var id string
	err := DB.QueryRow("CALL get_Uuid_from_URL_and_domain(?, ?)", url, domain).Scan(&id)
	if err != nil {
		return "", err
	} else {
		log.Printf("UUID for given URL and domain: %s", id)
	}
	return id, nil
}

// Function to fetch a random URL
//
// It defines a function that retrieves a random URL from a database, logs the URL, and returns it as a string, handling any potential errors during the database query.
func GetRandomURL() (string, error) {
	var id, url, tags, domain string
	var createdTime []byte // As per our earlier fix
	err := DB.QueryRow("CALL get_random_url()").Scan(&id, &url, &tags, &domain, &createdTime)
	if err != nil {
		return "", err
	} else {
		log.Printf("Get Random URL: %s", url)
	}
	return url, nil
}

// Function to fetch all URLs (just the 'url' column)
//
// Defines a Go function called "GetURLsOnly" that queries a database for URLs and returns them in a slice, handling potential errors along the way.
func GetURLsOnly() ([]string, error) {
	rows, err := DB.Query("CALL get_urls_only()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		} else {
			log.Printf("All URLs: %v", urls)
		}
		urls = append(urls, url)
	}

	return urls, rows.Err()
}

// Function to fetch all URLs with their tags
//
// It retrieves URL and tag data from a database, processes it into a map of URLs mapped to corresponding tags, and handles potential errors, logging intermediate results
func GetURLsAndTags() (map[string]map[string]interface{}, error) {
	rows, err := DB.Query("CALL get_urls_and_tags()")
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
		} else {
			log.Printf("All URLs and tags: %v", urlsAndTags)
		}

		var tags map[string]interface{}
		err = json.Unmarshal([]byte(tagsStr), &tags)
		if err != nil {
			return nil, err
		} else {
			log.Printf("All URLs and tags mapped: %v", urlsAndTags)
		}

		urlsAndTags[url] = tags
	}

	return urlsAndTags, rows.Err()
}
