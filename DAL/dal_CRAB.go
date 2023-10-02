package main

// Struct for a URL and its tags
type URLTag struct {
	URL  string
	Tags string
}

// Create a new web crawler with a source URL
func CreateWebCrawler(sourceURL string) error {
	_, err := db.Exec("CALL create_webcrawler(?)", sourceURL)
	return err
}

// Create a new scraper engine with name and description
func CreateScraperEngine(name, description string) error {
	_, err := db.Exec("CALL create_scraper_engine(?, ?)", name, description)
	return err
}
