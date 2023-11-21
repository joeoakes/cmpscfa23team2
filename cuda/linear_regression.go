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
)

// JSONData represents the structure of your entire JSON data
type JSONData struct {
	Items []Item `json:"items"`
}

// Item represents the structure of your JSON data
type Item struct {
	Domain string `json:"domain"`
	X      float64
	Y      float64
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
	return jsonData.Items
}

// processColumnValues processes the JSON records and returns a map of columns and their values.
func processColumnValues(items []Item) map[int]plotter.Values {
	columnsValues := make(map[int]plotter.Values)
	for i, item := range items {
		columnsValues[0] = append(columnsValues[0], float64(i+1))
		columnsValues[1] = append(columnsValues[1], item.X)
		columnsValues[2] = append(columnsValues[2], item.Y)
	}
	return columnsValues
}

// createHistogram creates and saves a histogram plot for the given values.
func createHistogram(values plotter.Values, title string, filename string) {
	p := plot.New()
	p.Title.Text = title

	h, err := plotter.NewHist(values, 16)
	if err != nil {
		log.Fatal(err)
	}
	h.Normalize(1)
	p.Add(h)

	if err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, filename); err != nil {
		log.Fatal(err)
	}
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
	filePath := "crab/template.json" // Update with your actual data file path and file type
	items := readJSON(filePath)

	// Domain to axis labels mapping
	axisLabels := map[string][2]string{
		"RealEstate": {"Square Footage", "Price (Thousands)"},
		"Healthcare": {"Age (Years)", "Blood Pressure"},
		"Weather":    {"Temperature (°C)", "Ice Cream Sales"},
	}

	// Group data by domain
	domainData := make(map[string][]Item)
	for _, item := range items {
		domain := item.Domain
		domainData[domain] = append(domainData[domain], item)
	}

	for domain, data := range domainData {
		fmt.Printf("\nAnalyzing domain: %s\n", domain)

		xValues := make([]float64, len(data))
		yValues := make([]float64, len(data))

		for i, item := range data {
			x := item.X
			y := item.Y

			xValues[i] = x
			yValues[i] = y
		}

		// Calculate linear regression coefficients
		a, b := linearRegression(xValues, yValues)
		fmt.Printf("Linear Regression Model for %s: y = %.2fx + %.2f\n", domain, a, b)

		// Retrieve axis labels for the domain
		labels := axisLabels[domain]

		// Create scatter plot with domain-specific axis labels
		plotFileName := fmt.Sprintf("%s_scatter_plot.png", domain)
		createScatterPlot(xValues, yValues, a, b, fmt.Sprintf("%s Linear Regression", domain), plotFileName, labels[0], labels[1])

		// Example prediction
		newX := 50.0       // Example new x value for prediction
		newY := a*newX + b // Predict y based on the regression model
		fmt.Printf("Prediction for %s: For x = %.2f, predicted y = %.2f\n\n", domain, newX, newY)
	}

	fmt.Println("Analysis complete.")
}
