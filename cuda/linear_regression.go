package cuda

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
)

// for airfare: -------------------------------------------------------------------------------------
// JSONData represents the structure of your entire JSON data
type JSONAirfareData struct {
	Domain string      `json:"domain"`
	Data   AirfareData `json:"data"`
}

// AirfareData represents the structure of each item in the JSON data for airfare
type AirfareData struct {
	Title          string            `json:"title"`
	Year           string            `json:"year"`
	Location       string            `json:"location"`
	Features       []string          `json:"features"`
	AdditionalInfo AirfareAdditional `json:"additional_info"`
	Metadata       AirfareMetadata   `json:"metadata"`
}

// AirfareAdditional represents additional information for airfare data
type AirfareAdditional struct {
	Country    string         `json:"country"`
	MonthsData []AirfareMonth `json:"months_data"`
}

// AirfareMonth represents each month's data for airfare
type AirfareMonth struct {
	Month string `json:"month"`
	Rate  string `json:"rate"`
	Year  string `json:"year"`
}

// AirfareMetadata represents metadata for airfare data
type AirfareMetadata struct {
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

// readAirfareJSON reads airfare JSON data from a file and returns a slice of AirfareData
func readAirfareJSON(filePath string) AirfareData {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data JSONAirfareData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data.Data
}

// extractPricesYearsAndCPIAirfare extracts numerical values from AirfareData and returns prices, months, and years
func extractPricesYearsAndMonths(data AirfareData) ([]float64, []string, []string) {
	var prices []float64
	var months []string
	var years []string

	for _, monthData := range data.AdditionalInfo.MonthsData {
		priceStr := strings.ReplaceAll(monthData.Rate, "$", "")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			// Skip invalid entries
			continue
		}

		prices = append(prices, price)
		months = append(months, monthData.Month)
		years = append(years, data.Year)
	}

	return prices, months, years
}

// for gas: -------------------------------------------------------------------------------------
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
func ReadGasJSON(filePath string) []GasolineData {
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
func ExtractPricesYearsAndCPI(items []GasolineData) ([]float64, []float64, []float64) {
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
	}

	return prices, years, cpiValues
}

// for books --------------------------------------------------------------------------------------

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

// LINEAR REGRESSION ----------------------------------------------------------------------------------
// linearRegression calculates the coefficients for a simple linear regression model (y = ax + b).
func LinearRegression(x, y []float64) (a, b float64) {
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
func LinearRegressionThreeVariables(x1, x2, y []float64) (a, b, c float64) {
	var sumX1, sumX2, sumY, sumX1Y, sumX2Y, sumX1X2, sumX1Squared, sumX2Squared float64
	n := float64(len(x1))

	fmt.Println("Lengths:", len(x1), len(x2), len(y))
	fmt.Println("x1:", x1)
	fmt.Println("x2:", x2)
	fmt.Println("y:", y)

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

// SCATTER PLOT --------------------------------------------------------------------------------

// createScatterPlot creates and saves a scatter plot with the linear regression line.
func CreateScatterPlot(x, y []float64, a, b float64, title, filename, xLabel, yLabel string) {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	// Print statements for debugging
	fmt.Println("Creating scatter plot...")

	pts := make(plotter.XYs, len(x))
	for i := range x {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatalf("Error creating scatter plot: %v", err)
	}
	p.Add(scatter)

	// Adding a line plot for the regression line
	line := plotter.NewFunction(func(x float64) float64 {
		return a*x + b
	})
	line.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Set the line color to red
	p.Add(line)

	// Save the plot to a PNG file.
	fmt.Printf("Saving plot to file: %s\n", filename)
	if err := p.Save(6*vg.Inch, 4*vg.Inch, filename); err != nil {
		log.Fatalf("Error saving plot: %v", err)
	}

	fmt.Println("Scatter plot created and saved successfully.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Accept user input for domain
	var domain string
	fmt.Print("Enter domain (gas/books/airfare): ")
	fmt.Scan(&domain)

	// Handle different domains
	switch domain {
	case "gas":
		filePath := "gasoline_data.json"
		GasolineData := ReadGasJSON(filePath)

		// Extract prices, years, and CPI values
		prices, years, cpiValues := ExtractPricesYearsAndCPI(GasolineData)

		// Perform linear regression
		a, b, c := LinearRegressionThreeVariables(years, cpiValues, prices)

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
		// Format prediction output with descriptive text
		descriptivePrediction := fmt.Sprintf("The prediction for gas prices in the year 2023 is: $%.2f", price2023)
		// Collect input data from previous years
		var inputDataStrings []string
		for i := 0; i < len(years); i++ {
			inputDataString := fmt.Sprintf("(%f, %f, %f)", years[i], cpiValues[i], prices[i])
			inputDataStrings = append(inputDataStrings, inputDataString)
		}

		// Combine input data strings into a single string
		inputData := strings.Join(inputDataStrings, ",")

		// Connect to the database
		db, err := sql.Open("mysql", "root:FJALkalim123@tcp(127.0.0.1:3306)/goengine")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Insert the prediction into the "predictions" table, including input_data
		insertStatement := "INSERT INTO linear_regression_predictions (prediction_id,query_identifier, input_data, prediction_info) VALUES (?,?, ?, ?)"
		// Insert the prediction with descriptive text into the database
		_, err = db.Exec(insertStatement, 1, "Gas Prices Query 1", inputData, descriptivePrediction)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Predicted Gas Price for 2023: %.2f\n", price2023)

		// Append the predicted gas price to the existing prices slice
		prices = append(prices, price2023)

		// Create and save the scatter plot with the extended x values
		title := "Gas Price Prediction Scatter Plot (Extended)"
		filename := "Gas Prices Prediction for the year 2023_scatter_plot.png"
		xLabel := "Year"
		yLabel := "Average Gasoline Prices"
		CreateScatterPlot(append(years, newX...), prices, a, b, title, filename, xLabel, yLabel)

	// case airfare ---------------------------------------------------
	case "airfare":
		filePath := "airfare_data_price.json"
		airfareData := readAirfareJSON(filePath)

		// Extract prices, months, and inflation rate values
		prices, months, years := extractPricesYearsAndMonths(airfareData)

		// Perform linear regression
		indices := make([]float64, len(months))
		for i := 0; i < len(months); i++ {
			indices[i] = float64(i + 1)
		}

		// Convert years to numeric values
		var numericYears []float64
		for i := 0; i < len(years); i++ {
			numericYears = append(numericYears, float64(i+1))
		}

		// Perform linear regression
		a, b, c := LinearRegressionThreeVariables(indices, prices, numericYears)

		// Output the prediction for new x values (months)
		// Example new x values for prediction
		newX := indices

		var newY []float64
		for _, xVal := range numericYears {
			yVal := a*xVal + b*float64(len(newX)+1) + c*float64(len(newX)) // Predict y based on the regression model
			newY = append(newY, yVal)
		}

		// Connect to the database !
		db, err := sql.Open("mysql", "root:Matthew@1499.@tcp(127.0.0.1:3306)/goengine")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Collect input data from previous months
		var inputDataStrings []string
		for i := 0; i < len(indices); i++ {
			inputDataString := fmt.Sprintf("(%f, %f, %f)", indices[i], numericYears[i], prices[i])
			inputDataStrings = append(inputDataStrings, inputDataString)
		}

		// Combine input data strings into a single string
		inputData := strings.Join(inputDataStrings, ",")

		// Convert the numeric prediction_info value to a JSON-formatted string
		predictionInfo, err := json.Marshal(newY)
		if err != nil {
			log.Fatal(err)
		}

		// Insert the prediction into the "predictions" table, including input_data
		insertStatement := "INSERT INTO linear_regression_predictions (prediction_id,query_identifier, input_data, prediction_info) VALUES (?,?, ?, ?)"
		_, err = db.Exec(insertStatement, 3, "Airefare infaltion rate query 1", inputData, predictionInfo)
		if err != nil {
			log.Fatal(err)
		}

		// Create and save the scatter plot with the extended x valuesf
		title := "Airfare Price Prediction Scatter Plot"
		filename := "airfare_scatter_plot.png"
		xLabel := "Month"
		yLabel := "Average Airfare Prices"
		CreateScatterPlot(indices, prices, a, b, title, filename, xLabel, yLabel)

	// case books __________________________________________________________________________________
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
		a, b := LinearRegression(indices, prices)

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
		CreateScatterPlot(indices, prices, a, b, title, filename, xLabel, yLabel)

	default:
		log.Fatal("Unknown domain:", domain)
	}

}
