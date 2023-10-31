package main

import (
	"cmpscfa23team2/DAL"
	// "cmpscfa23team2/DAL"
	_ "encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	_ "os"
)

func main() {

	fmt.Printf("Web Scraper\n")

	// to create a new collector using default settings
	CollectWebsite := colly.NewCollector()

	// request must be made before the site is visited
	CollectWebsite.OnRequest(func(req *colly.Request) {
		fmt.Println("Visiting: ", req.URL)
	})

	// called if an error occurs during the request
	CollectWebsite.OnError(func(_ *colly.Response, err error) {
		log.Println("An error occurred while processing request: ", err)
	})

	// Called after getting a response from the server
	CollectWebsite.OnResponse(func(req *colly.Response) {
		fmt.Println("Page(s) visited: ", req.Request.URL)
	})

	// targets the <a> tag for hyper links
	CollectWebsite.OnHTML("a", func(element *colly.HTMLElement) {
		// prints all the URLs associated with the hyperlink destination
		link := element.Attr("href")
		if link != "" {
			fmt.Println(link)
		}
	})

	CollectWebsite.OnScraped(func(req *colly.Response) {
		fmt.Println(req.Request.URL, " scraped!")
	})

	// code to connect to DAL web crawler
	sourceURL := `https://www.housevalues.com/?hvnp=1&msclkid=77c1d49a0b91197bebd8b1a4ade8f956#utm_source=bing&utm_medium=cpc&utm_campaign=bing&utm_content=RegionalHotCPA-HousePrices-Broad&utm_term=A`
	crawlerID, err := DAL.CreateWebCrawler(sourceURL)
	if err != nil {
		log.Fatal("Failed to create web crawler: ", err)
	} else {
		log.Printf("Web crawler created: %s", crawlerID)
	}

	url := "https://www.housevalues.com/?hvnp=1&msclkid=77c1d49a0b91197bebd8b1a4ade8f956#utm_source=bing&utm_medium=cpc&utm_campaign=bing&utm_content=RegionalHotCPA-HousePrices-Broad&utm_term=A"
	domain := "housevalues.com"
	tags := map[string]interface{}{
		"Property Type": "House",
	}

	// Insert the URL from DAL with associated information
	id, err := DAL.InsertURL(url, domain, tags)
	if err != nil {
		log.Fatal("Failed to insert URL: ", err)
	} else {
		log.Printf("URL inserted with ID: %s", id)
	}

	// downloading the target HTML page (housing prices)
	err = CollectWebsite.Visit("https://www.housevalues.com/?hvnp=1&msclkid=77c1d49a0b91197bebd8b1a4ade8f956#utm_source=bing&utm_medium=cpc&utm_campaign=bing&utm_content=RegionalHotCPA-HousePrices-Broad&utm_term=A")
	if err != nil {
		return
	}

	CollectWebsite.OnHTML("li.product", func(element *colly.HTMLElement) {
		// TO DO ...
	})

	// keeps the program running
	// select {}
}

// reference: https://www.zenrows.com/blog/web-scraping-golang#set-up-go-project
