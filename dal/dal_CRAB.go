package dal

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
		InsertLog("400", "Error creating web crawler: "+err.Error(), "CreateWebCrawler()")
		return "", err
	} else {
		InsertLog("200", "Web crawler created: "+crawlerID, "CreateWebCrawler()")
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
		InsertLog("400", "Error creating scraper engine: "+err.Error(), "CreateScraperEngine()")
		return "", err
	} else {
		InsertLog("200", "Scraper engine created: "+engineID, "CreateScraperEngine()")
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
		InsertLog("400", "Error marshalling tags: "+err.Error(), "InsertURL()")
		return "", err
	} else {
		InsertLog("200", "URL inserted successfully", "InsertURL()")
		log.Printf("URL inserted with tags: %v", tags)
	}

	err = DB.QueryRow("CALL insert_url(?, ?, ?)", url, string(jsonTags), domain).Scan(&id)
	if err != nil {
		InsertLog("400", "Error inserting URL: "+err.Error(), "InsertURL()")
		return "", err
	} else {
		InsertLog("200", "URL inserted with ID: "+id, "InsertURL()")
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
		InsertLog("400", "Error marshalling tags: "+err.Error(), "UpdateURL()")
		return err
	} else {
		InsertLog("200", "URL updated with tags sucessfully", "UpdateURL()")
		log.Printf("URL updated with tags: %v", tags)
	}

	_, err = DB.Exec("CALL update_url(?, ?, ?, ?)", id, url, string(jsonTags), domain)
	if err != nil {
		InsertLog("400", "Error updating URL: "+err.Error(), "UpdateURL()")
	}
	return err
}

// Function to fetch URL tags and domain by ID
//
// It defines a function that retrieves tags and a domain from a database using a specified ID, logs the results, and returns them in a map and a string along with potential errors.
func GetURLTagsAndDomain(id string) (map[string]interface{}, string, error) {
	var tagsStr, domain string
	err := DB.QueryRow("CALL get_url_tags_and_domain(?)", id).Scan(&tagsStr, &domain)
	if err != nil {
		InsertLog("400", "Error getting URL tags and domain: "+err.Error(), "GetURLTagsAndDomain()")
		return nil, "", err
	} else {
		InsertLog("200", "Tags retrieved successfully", "GetURLTagsAndDomain()")
		log.Printf("Tags: %v, Domain: %s", tagsStr, domain)
	}
	var tags map[string]interface{}
	err = json.Unmarshal([]byte(tagsStr), &tags)
	if err != nil {
		InsertLog("400", "Error unmarshalling tags: "+err.Error(), "GetURLTagsAndDomain()")
		return nil, "", err
	} else {
		InsertLog("200", "Tags marshalled successfully", "GetURLTagsAndDomain()")
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
		InsertLog("400", "Error getting URLs from domain: "+err.Error(), "GetURLsFromDomain()")
		return nil, err
	}
	log.Println("Closing Rows: %+v", rows)
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var id, url, tags, domain string
		var createdTime []byte // <-- Change this line
		if err := rows.Scan(&id, &url, &tags, &domain, &createdTime); err != nil {
			InsertLog("400", "Error scanning rows: "+err.Error(), "GetURLsFromDomain()")
			return nil, err
		} else {
			InsertLog("200", "URLs from domain extracted successfully", "GetURLsFromDomain()")
			log.Printf("URLs from domain: %+v", urls)
		}
		urls = append(urls, url)
	}
	return urls, rows.Err()
}
