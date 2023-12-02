package main

import "time"

// URLData represents the URL to be crawled and its related data.
type URLData struct {
	URL     string    // The URL to be crawled
	Created time.Time // Timestamp of URL creation or retrieval
	Links   []string  // URLs found on this page
}

// MonthData represents data for each month.
type MonthData struct {
	Month string `json:"month"`
	Rate  string `json:"rate"`
}

// AirfareData represents airfare data structure.
type AirfareData struct {
	Domain string `json:"domain"`
	URL    string `json:"url"`
	Data   struct {
		Title          string   `json:"title"`
		Year           string   `json:"year"`
		Location       string   `json:"location"`
		Features       []string `json:"features"`
		AdditionalInfo struct {
			Country    string      `json:"country"`
			MonthsData []MonthData `json:"months_data"`
		} `json:"additional_info"`
		Metadata struct {
			Source    string `json:"source"`
			Timestamp string `json:"timestamp"`
		} `json:"metadata"`
	} `json:"data"`
}

// YearData represents data for each year.
type YearData struct {
	Year string `json:"year"`
	Jan  string `json:"jan"`
	Feb  string `json:"feb"`
	Mar  string `json:"mar"`
	Apr  string `json:"apr"`
	May  string `json:"may"`
	Jun  string `json:"jun"`
	July string `json:"july"`
	Aug  string `json:"aug"`
	Sept string `json:"sept"`
	Oct  string `json:"oct"`
	Nov  string `json:"nov"`
	Dec  string `json:"dec"`
	Avg  string `json:"avg"`
}

// GasolineData represents gasoline data structure.
type GasolineData struct {
	Year                     string `json:"year"`
	AverageGasolinePrices    string `json:"average_gasoline_prices"`
	AverageAnnualCPIForGas   string `json:"average_annual_cpi_for_gas"`
	GasPricesAdjustedForInfl string `json:"gas_prices_adjusted_for_inflation"`
}

// PropertyData represents property data structure.
type PropertyData struct {
	Status    string `json:"status"`
	Bedrooms  string `json:"bedrooms"`
	Bathrooms string `json:"bathrooms"`
	AcreLot   string `json:"acre_lot"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
	HouseSize string `json:"house_size"`
	SoldDate  string `json:"prev_sold_date"`
	Price     string `json:"price"`
}

// ScraperConfig represents configuration for the scraper.
type ScraperConfig struct {
	StartingURLs []string
}

// DomainConfig represents configuration for different domains.
type DomainConfig struct {
	Name                            string
	ItemSelector                    string
	TitleSelector                   string
	URLSelector                     string
	DescriptionSelector             string
	PriceSelector                   string
	FactorsSelector                 string
	DepreciationRatesSelector       string
	ModelsLeastDepreciationSelector string
	ModelsMostDepreciationSelector  string
}

// Metadata represents metadata for scraped data.
type Metadata struct {
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

// GenericData represents generic data structure for items.
type GenericData struct {
	Title                   string            `json:"title"`
	URL                     string            `json:"url"`
	Description             string            `json:"description"`
	Price                   string            `json:"price"`
	Factors                 []string          `json:"factors"`
	DepreciationRates       map[string]string `json:"depreciation_rates"`
	ModelsLeastDepreciation []string          `json:"models_least_depreciation"`
	ModelsMostDepreciation  []string          `json:"models_most_depreciation"`
	Metadata                Metadata          `json:"metadata"`
}

// ItemData represents data for an item.
type ItemData struct {
	Domain string        `json:"domain"`
	Data   []GenericData `json:"data"`
}
