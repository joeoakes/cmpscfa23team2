// Testing The Security/Scrapability Of URL's

package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func main() {
	// Create a new Collector instance from the colly package.
	c := colly.NewCollector()

	// Set a random User-Agent header for the HTTP requests to make them look more like a web browser.
	extensions.RandomUserAgent(c)

	// Define a callback function that will be executed when HTML elements are found.
	c.OnHTML("*", func(e *colly.HTMLElement) {
		// Extract the "id" attribute of the HTML element.
		id := e.Attr("id")
		if id != "" {
			// If the "id" attribute is not empty, print a message indicating the presence of an element with an ID.
			fmt.Println("Found element with ID:", id)
		}
	})

	// Send an HTTP GET request to the specified URL.
	err := c.Visit("https://www.zillow.com/home-values/47/pa/")
	if err != nil {
		// Handle and print any error that occurred during the request.
		fmt.Println("Error:", err)
	}
}

//---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Testing Domains:
// GitHub: (Works/Will Scrape)
// https://www.github.com/thepalad1n
//
// Airfare: (Works/Will Scrape)
// https://www.nerdwallet.com/article/travel/travel-price-tracker
//
// NFL (Football): (Will Not Scrape (Forbidden))
// https://v2.playoffpredictors.com/football/nfl/37033920-C0E1-4EF4-8F0D-DA53DA41E3A0/official?L=Aw18ZXTt-DFKQRgExpej7df5pnjnqYQRiRVbuToUVvdhdg08ZZR02suEA
//
// Car Value/Depreciation: (Works/Will Scrape)
// https://www.thinkinsure.ca/insurance-help-centre/car-deprecation.html
//
// Weather: (Works/Will Scrape)
// https://www.accuweather.com/
//
// Travel: (Works/Will Scrape)
// https://www.travelpulse.com/gallery/features/the-15-travel-trends-that-will-define-2023
//
// Human Population: (Works/Will Scrape)
// https://thepopulationproject.org/
//
// NASCAR: (Will Not Scrape (Forbidden))
// https://www.nascar.com/gallery/predicting-every-2023-cup-series-playoffs-race-winner/
//
// NASCAR Betting: (Works/Will Scrape)
// https://www.predictem.com/nascar/xfinity-500-race-preview-picks/
//
// Zillow: (Will Not Scrape (Forbidden))
// https://www.zillow.com/home-values/47/pa/
