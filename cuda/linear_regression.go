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

// extractPricesAndYears extracts numerical values from price strings and returns prices, years, and cpiValues.
func extractPricesYearsAndCPI(items []GasolineData) ([]float64, []float64, []float64) {
	var prices []float64
	var years []float64
	var cpiValues []float64

	for _, item := range items {
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

		cpiStr := strings.ReplaceAll(item.AverageAnnualCPIForGas, "$", "")
		cpi, err := strconv.ParseFloat(cpiStr, 64)
		if err != nil {
			log.Fatal(err)
		}

		prices = append(prices, price)
		years = append(years, year)
		cpiValues = append(cpiValues, cpi)

		// fmt.Printf("Item %d Year: %.2f, Price: %.2f, CPI: %.2f\n", i+1, year, price, cpi)
	}

	return prices, years, cpiValues
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

// LINEAR REGRESSION ---------------------------------------------------

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

// linearRegressionThreeVariables calculates the coefficients for a multiple linear regression model (y = ax1 + bx2 + c).
func linearRegressionThreeVariables(x1, x2, y []float64) (a, b, c float64) {
	var sumX1, sumX2, sumY, sumX1Y, sumX2Y, sumX1X2, sumX1Squared, sumX2Squared float64
	n := float64(len(x1))

	for i := 0; i < len(x1); i++ {
		sumX1 += x1[i]
		sumX2 += x2[i]
		sumY += y[i]
		sumX1Y += x1[i] * y[i]
		sumX2Y += x2[i] * y[i]
		sumX1X2 += x1[i] * x2[i]
		sumX1Squared += x1[i] * x1[i]
		sumX2Squared += x2[i] * x2[i]
	}

	denominator := n*sumX1Squared*sumX2Squared - sumX1*sumX1*sumX2*sumX2 - sumX2*sumX2*sumX1Squared + 2*sumX1*sumX2*sumX1X2 - sumX1X2*sumX1X2
	a = (sumY*sumX1Squared*sumX2Squared - sumX1*sumX1Y*sumX2Squared - sumX2*sumX2Y*sumX1Squared +
		2*sumX1Y*sumX2*sumX1X2 - sumX2Y*sumX1X2*sumX1 - sumY*sumX2*sumX1X2) / denominator
	b = (n*sumX2Y*sumX1Squared - sumX2*sumY*sumX1Squared - sumX2Y*sumX1*sumX1 +
		2*sumY*sumX1*sumX1X2 - sumY*sumX1X2*sumX2) / denominator
	c = (-sumX2Squared*sumX1Y + sumX2*sumX1*sumY + sumX2Squared*sumX1*sumX1Y -
		sumX1X2*sumY*sumX2 + sumX1X2*sumX2Y*sumX1 - sumX1*sumX2*sumX2Y) / denominator

	return a, b, c
}

// SCATTER PLOT --------------------------------------------------------

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

		// Extract prices, years, and CPI values
		prices, years, cpiValues := extractPricesYearsAndCPI(GasolineData)

		// Perform linear regression
		a, b, c := linearRegressionThreeVariables(years, cpiValues, prices)

		// Extend the time range for prediction (next year)
		var newX []float64
		for i := 1; i <= 1; i++ {
			newX = append(newX, years[len(years)-1]+float64(i))
		}

		// Predict y based on the regression model for the extended x values
		var newY []float64
		for _, xVal := range newX {
			yVal := a*xVal + b
			newY = append(newY, yVal)
		}

		// Predict gas price for 2023
		year2023 := 2023.0
		cpi2023 := 349.189
		price2023 := (0.10 * (a*year2023 + b*cpi2023 + c))

		fmt.Printf("Predicted Gas Price for 2023: $%.2f\n", price2023)

		// Append the predicted gas price to the existing prices slice
		prices = append(prices, price2023)

		// Create and save the scatter plot with the extended x values
		title := "Gas Price Prediction Scatter Plot (Extended)"
		filename := "gas_scatter_plot_extended.png"
		xLabel := "Year"
		yLabel := "Average Gasoline Prices"
		createScatterPlot(append(years, newX...), prices, a, b, title, filename, xLabel, yLabel)

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
