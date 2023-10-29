package main

import (
	"fmt"
	"log"

	/* "encoding/csv"
	"log"
	"os"*/
	"github.com/gocolly/colly"
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

	// downloading the target HTML page (housing prices)
	err := CollectWebsite.Visit("https://www.housevalues.com/?hvnp=1&msclkid=77c1d49a0b91197bebd8b1a4ade8f956#utm_source=bing&utm_medium=cpc&utm_campaign=bing&utm_content=RegionalHotCPA-HousePrices-Broad&utm_term=A")
	if err != nil {
		return
	}

	// put scraping elements here

	// keeps the program running
	select {}
}

// reference: https://www.zenrows.com/blog/web-scraping-golang#set-up-go-project
