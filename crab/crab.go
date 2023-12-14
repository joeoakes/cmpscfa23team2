package main

import "fmt"

// begin main ===========================================================================================================
func main() {

	//// Run the Python scraper script
	//cmd := exec.Command("python", "C:\\Users\\mathe\\GolandProjects\\cmpscfa23team2\\crab\\job_listings_scraper.py")
	//
	//var stdout, stderr bytes.Buffer
	//cmd.Stdout = &stdout
	//cmd.Stderr = &stderr
	//
	//cmd.Run()
	//
	//fmt.Println("Scraper executed successfully. Check the output directory for the JSON file.")
	//fmt.Println("Stdout:", stdout.String())
	//begin crawler

	//airdatatest()
	//scrapeInflationData()
	//scrapeGasInflationData()
	//scrapeHousingData()
	//end crawler

	//begin scraper
	fmt.Println("Available domains:")
	for domainName := range domainConfigurations {
		fmt.Printf("- %s\n", domainName)
	}

	// Ask the user to choose a domain
	var domainName string
	fmt.Print("Enter the domain you want to scrape: ")
	fmt.Scanln(&domainName)

	// Check if the chosen domain is valid
	_, exists := domainConfigurations[domainName]
	if !exists {
		fmt.Printf("Invalid domain name provided: %s\n", domainName)
		return
	}
	//
	//// Perform the scraping for the chosen domain
	testScrape(domainName)
	//
	////csvread
	//filePath := "crab/csv"
	//properties, err := ReadCSV(filePath)
	//if err != nil {
	//	fmt.Printf("Error reading CSV file: %s\n", err)
	//	return
	//}

	// Print the PropertyData for demonstration purposes
	//for _, property := range properties {
	//	fmt.Printf("%+v\n", property)
	//}
}

//end main =============================================================================================================
