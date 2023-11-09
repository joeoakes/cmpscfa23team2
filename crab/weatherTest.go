package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"io/ioutil"
	"time"
)

// WeatherData represents the weather information for a single day
type WeatherData struct {
	Date            string `json:"date"`
	HighTemperature string `json:"high_temperature"`
	LowTemperature  string `json:"low_temperature"`
}

func main() {
	city := "chalfont" // Use a default city name
	zip := "18914"     // Use a default zip code

	// Container for all weather data
	allWeatherData := []WeatherData{}

	// all data
	// startDate := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	// endDate := time.Date(2023, time.November, 30, 0, 0, 0, 0, time.UTC)

	startDate := time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, time.November, 30, 0, 0, 0, 0, time.UTC)
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	// Iterate over the months and years between start and end dates
	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 1, 0) {
		currentYear, currentMonth, _ := d.Date()
		monthName := currentMonth.String()
		startingURL := fmt.Sprintf("https://www.accuweather.com/en/us/%s/%s/%s-weather/2215752?year=%d", city, zip, monthName, currentYear)
		fmt.Println("Constructed URL:", startingURL)

		// Container for the current month's data
		var monthlyWeatherData []WeatherData

		monthlyCollector := c.Clone()
		extensions.RandomUserAgent(monthlyCollector)

		var currentWeatherData WeatherData // Temporarily hold the data for each day
		monthlyCollector.OnHTML("div.date", func(e *colly.HTMLElement) {
			currentWeatherData = WeatherData{Date: e.Text} // Initialize a new WeatherData
		})
		monthlyCollector.OnHTML("div.high", func(e *colly.HTMLElement) {
			currentWeatherData.HighTemperature = e.Text // Set high temperature
		})
		monthlyCollector.OnHTML("div.low", func(e *colly.HTMLElement) {
			currentWeatherData.LowTemperature = e.Text                          // Set low temperature
			monthlyWeatherData = append(monthlyWeatherData, currentWeatherData) // Add to slice
		})

		// Visit the URL with retry logic
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			err := monthlyCollector.Visit(startingURL)
			if err == nil {
				break // No error, break the retry loop
			}
			fmt.Printf("Error visiting %s: %s, retrying (%d/%d)\n", startingURL, err, i+1, maxRetries)
			if i < maxRetries-1 {
				time.Sleep(time.Second * 10) // Wait before retrying
			}
		}

		// Append the month's data to `allWeatherData`
		allWeatherData = append(allWeatherData, monthlyWeatherData...)

		// Sleep to prevent rate-limiting issues
		time.Sleep(time.Second * 5)
	}

	// Marshal the data into JSON
	jsonData, err := json.MarshalIndent(allWeatherData, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling data to JSON:", err)
		return
	}

	// Write the data to a file
	err = ioutil.WriteFile("weatherData.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
	}
}

// go run main.go --city="warrington" --location="18974" --month="october" --year="2023"
