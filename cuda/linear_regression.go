package main

import (
	"encoding/json"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
)

// for gas: ---------------------------------------------------------
// JSONData represents the structure of your entire JSON data
type JSONGasData struct {
	Domain string         `json:"domain"`
	Data   []GasolineData `json:"data"`
}

// GasolineData represents the structure of each item in the JSON data for gas
type GasolineData struct {
	Year                     string `json:"year"`
	AverageGasolinePrices    string `json:"average_gasoline_prices"`
	AverageAnnualCPIForGas   string `json:"average_annual_cpi_for_gas"`
	GasPricesAdjustedForInfl string `json:"gas_prices_adjusted_for_inflation"`
}

// read gas
func readGasJSON(filePath string) []GasolineData {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data []GasolineData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// extractPricesAndYears extracts numerical values from price strings and returns both prices and years.
func extractPricesAndYears(items []GasolineData) ([]float64, []float64) {
	var prices []float64
	var years []float64

	for i, item := range items {
		priceStr := strings.ReplaceAll(item.AverageGasolinePrices, "$", "")
		priceStr = strings.ReplaceAll(priceStr, ",", "")

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			log.Fatal(err)
		}

		year, err := strconv.ParseFloat(item.Year, 64)
		if err != nil {
			log.Fatal(err)
		}

		prices = append(prices, price)
		years = append(years, year)

		fmt.Printf("Item %d Year: %.2f, Price: %.2f\n", i+1, year, price)
	}

	return prices, years
}

// for books ---------------------------------------------------------

// JSONData represents the structure of your entire JSON data
type JSONData struct {
	Domain string `json:"domain"`
	Data   []Item `json:"data"`
}

// Item represents the structure of each item in the JSON data
type Item struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Metadata    struct {
		Source    string `json:"source"`
		Timestamp string `json:"timestamp"`
	} `json:"metadata"`
}

// readJSON reads data from a JSON file and returns a slice of records.
func readJSON(filePath string) []Item {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var jsonData JSONData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonData)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData.Data
}

// extractPrices extracts numerical values from price strings.
func extractPrices(items []Item) []float64 {
	var prices []float64
	for i, item := range items {
		priceStr := strings.ReplaceAll(item.Price, "£", "")
		priceStr = strings.ReplaceAll(priceStr, ",", "")

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			log.Fatal(err)
		}

		prices = append(prices, price)
		fmt.Printf("Item %d Price: %.2f\n", i+1, price)
	}
	return prices
}

// linearRegression calculates the coefficients for a simple linear regression model (y = ax + b).
func linearRegression(x, y []float64) (a, b float64) {
	var sumX, sumY, sumXY, sumX2 float64
	n := float64(len(x))

	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	a = (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	b = (sumY - a*sumX) / n

	return a, b
}

// createScatterPlot creates and saves a scatter plot with the linear regression line.
func createScatterPlot(x, y []float64, a, b float64, title, filename, xLabel, yLabel string) {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	pts := make(plotter.XYs, len(x))
	for i := range x {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}
	p.Add(scatter)

	// Adding a line plot for the regression line
	line := plotter.NewFunction(func(x float64) float64 {
		return a*x + b
	})
	line.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Set the line color to red
	p.Add(line)

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, filename); err != nil {
		log.Fatal(err)
	}
}

// createResidualsPlot creates and saves a plot of residuals.
func createResidualsPlot(x, y []float64, a, b float64, filename, xLabel, yLabel string) {
	p := plot.New()

	p.Title.Text = "Residuals Plot"
	p.X.Label.Text = xLabel
	p.Y.Label.Text = "Residuals"

	// Calculate residuals
	residuals := make(plotter.XYs, len(x))
	for i := range x {
		residuals[i].X = x[i]
		residuals[i].Y = y[i] - (a*x[i] + b)
	}

	// Add a scatter plot for residuals
	scatter, err := plotter.NewScatter(residuals)
	if err != nil {
		log.Fatal(err)
	}
	p.Add(scatter)

	// Save the residuals plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, filename); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Accept user input for domain
	var domain string
	fmt.Print("Enter domain (gas/books): ")
	fmt.Scan(&domain)

	// Handle different domains
	switch domain {
	case "gas":
		filePath := "gasoline_data.json"
		GasolineData := readGasJSON(filePath)

		// Extract prices and years
		prices, years := extractPricesAndYears(GasolineData)

		// Print prices and years for troubleshooting
		fmt.Println("Prices:", prices)
		fmt.Println("Years:", years)

		// Perform linear regression for gas
		yearsForRegression := make([]float64, len(GasolineData))
		for i := 0; i < len(GasolineData); i++ {
			yearsForRegression[i], _ = strconv.ParseFloat(GasolineData[i].Year, 64)
		}
		// Perform linear regression
		a, b := linearRegression(yearsForRegression, prices)

		// Extend the time range for prediction (next 20 years)
		var newX []float64
		for i := 1; i <= 20; i++ {
			newX = append(newX, years[len(years)-1]+float64(i))
		}

		// Output the prediction for the extended x values
		fmt.Println("Extended X Values:")
		for _, xVal := range newX {
			fmt.Printf("%.2f ", xVal)
		}

		// Predict y based on the regression model for the extended x values
		var newY []float64
		for _, xVal := range newX {
			yVal := a*xVal + b
			newY = append(newY, yVal)
		}

		// Output the prediction for the extended x values
		fmt.Println("\nPredicted Values:")
		for _, yVal := range newY {
			fmt.Printf("%.2f ", yVal)
		}

		// Create and save the scatter plot with the extended x values
		title := "Gas Price Prediction Scatter Plot (Extended)"
		filename := "gas_scatter_plot_extended.png"
		residualsPlotFilename := "gas_residuals_plot.png"
		xLabel := "Year"
		yLabel := "Average Gasoline Prices"
		createScatterPlot(append(years, newX...), append(prices, newY...), a, b, title, filename, xLabel, yLabel)

		// Create and save the residuals plot
		createResidualsPlot(years, prices, a, b, residualsPlotFilename, xLabel, yLabel)

	// case books ____________________________________________________
	case "books":
		filePath := "books_data.json"
		items := readJSON(filePath)

		// Extract prices
		prices := extractPrices(items)

		// Print prices for troubleshooting
		fmt.Println("Prices:", prices)

		// Perform linear regression for books
		indices := make([]float64, len(items))
		for i := 0; i < len(items); i++ {
			indices[i] = float64(i + 1)
		}
		// Perform linear regression
		a, b := linearRegression(indices, prices)

		// Output the prediction for new x values
		newX := indices // Example new x values for prediction
		var newY []float64

		for _, xVal := range newX {
			yVal := a*xVal + b // Predict y based on the regression model
			newY = append(newY, yVal)
		}

		// Output the prediction for a new x value
		// Example new x value for prediction
		// Predict y based on the regression model
		fmt.Println("New X Values:")
		for _, xVal := range newX {
			fmt.Printf("%.2f ", xVal)
		}
		fmt.Println("\nPredicted Values:")
		for _, yVal := range newY {
			fmt.Printf("%.2f ", yVal)
		}

		// Create and save the scatter plot
		title := "Book Price Prediction Scatter Plot"
		filename := "book_scatter_plot.png"
		xLabel := "Book Index"
		yLabel := "Price"
		createScatterPlot(indices, prices, a, b, title, filename, xLabel, yLabel)

	default:
		log.Fatal("Unknown domain:", domain)
	}
}
