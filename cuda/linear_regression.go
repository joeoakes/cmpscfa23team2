package main

import (
	"encoding/csv"
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

// readCSV reads data from a CSV file and returns a slice of records.
func readCSV(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for i, record := range records {

		if len(record) != 21 { // Expecting 3 columns: domain, x, y
			log.Fatalf("Error in line %d: Expected 3 fields, got %d", i+1, len(record))
		}
		fmt.Printf("Line %d: %v\n", i+1, record)
	}
	return records
}

// processColumnValues processes the CSV records and returns a map of columns and their values.
func processColumnValues(records [][]string, startColumn int) map[int]plotter.Values {
	columnsValues := make(map[int]plotter.Values)
	for i, record := range records {
		for c := startColumn; c < len(record); c++ {
			if _, found := columnsValues[c]; !found {
				columnsValues[c] = make(plotter.Values, len(records))
			}
			floatVal, err := strconv.ParseFloat(record[c], 64)
			if err != nil {
				log.Fatal(err)
			}
			columnsValues[c][i] = floatVal
		}
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
	filePath := "cuda\\data2.csv"
	records := readCSV(filePath)

	// Domain to axis labels mapping
	axisLabels := map[string][2]string{
		"RealEstate": {"Square Footage", "Price (Thousands)"},
		"Healthcare": {"Age (Years)", "Blood Pressure"},
		"Weather":    {"Temperature (°C)", "Ice Cream Sales"},
	}

	// Group data by domain
	domainData := make(map[string][][]string)
	for _, record := range records[1:] { // Skipping header
		domain := record[0]
		domainData[domain] = append(domainData[domain], record[1:21]) // Taking only x and y
	}

	for domain, data := range domainData {
		fmt.Printf("\nAnalyzing domain: %s\n", domain)

		xValues := make([]float64, len(data))
		yValues := make([]float64, len(data))

		for i, record := range data {
			xStr := strings.TrimSpace(record[0])
			yStr := strings.TrimSpace(record[1])

			x, err := strconv.ParseFloat(xStr, 64)
			if err != nil {
				log.Fatalf("Error parsing x value in domain %s: %v", domain, err)
			}
			y, err := strconv.ParseFloat(yStr, 64)
			if err != nil {
				log.Fatalf("Error parsing y value in domain %s: %v", domain, err)
			}

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
