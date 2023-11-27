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
		fmt.Printf("Book %d Price: %.2f\n", i+1, price)
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

func main() {
	filePath := "book.json"
	items := readJSON(filePath)

	// Extract prices
	prices := extractPrices(items)

	// Print prices for troubleshooting
	fmt.Println("Prices:", prices)

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

	// // Output the prediction for a new x value
	// newX := indices      // Example new x value for prediction
	// newY := a*newX + b // Predict y based on the regression model
	fmt.Println("New X Values:")
	for _, xVal := range newX {
		fmt.Printf("%.2f ", xVal)
	}
	fmt.Println("\nPredicted Y Values:")
	for _, yVal := range newY {
		fmt.Printf("%.2f ", yVal)
	}
	// // Print the output
	// fmt.Printf("X Value: %.2f\n", newX)
	// fmt.Printf("Y Value: %.2f\n", newY)
	// fmt.Printf("Prediction: %.2f\n", newY)

	// fmt.Println("\nAnalysis complete.")
}
