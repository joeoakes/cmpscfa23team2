package cuda_test

import (
	"cmpscfa23team2/cuda"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"testing"
)

// testing
// TestLinearRegression tests the linearRegression function.
func TestLinearRegression(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 5, 4, 5}

	// Expected coefficients
	expectedA := 0.6
	expectedB := 2.2

	// Call the function
	a, b := cuda.LinearRegression(x, y)

	// Check the results
	if a != expectedA || b != expectedB {
		t.Errorf("Linear regression failed. Expected (a, b): (%f, %f), got (a, b): (%f, %f)", expectedA, expectedB, a, b)
	}
}

// TestLinearRegressionThreeVariables tests the linearRegressionThreeVariables function.
func TestLinearRegressionThreeVariables(t *testing.T) {
	x1 := []float64{1, 2, 3, 4, 5}
	x2 := []float64{2, 4, 5, 4, 5}
	y := []float64{3, 5, 7, 8, 10}

	// Expected coefficients
	expectedA := 0.560012
	expectedB := 0.138779
	expectedC := -3.854065

	// Tolerance for floating-point comparisons
	delta := 1e-6

	// Call the function
	a, b, c := cuda.LinearRegressionThreeVariables(x1, x2, y)

	// Check the results with tolerance
	if !almostEqual(a, expectedA, delta) || !almostEqual(b, expectedB, delta) || !almostEqual(c, expectedC, delta) {
		t.Errorf("Linear regression failed. Expected (a, b, c): (%f, %f, %f), got (a, b, c): (%f, %f, %f)", expectedA, expectedB, expectedC, a, b, c)
	}
}

// almostEqual checks if two floating-point numbers are almost equal within a given delta.
func almostEqual(a, b, delta float64) bool {
	return math.Abs(a-b) < delta
}

// TestReadJSON tests the readJSON function.
func TestReadJSON(t *testing.T) {
	filePath := "/Users/Sara/GolandProjects/cmpscfa23team2/gasoline_data.json"
	items := cuda.ReadGasJSON(filePath)

	// Perform your assertions on items, e.g., check the length of the data.
	if len(items) == 0 {
		t.Errorf("Expected non-empty items data, got empty data")
	}
	// Add more assertions based on your data structure.
}

// TestExtractPricesYearsAndCPI tests the extractPricesYearsAndCPI function.
func TestExtractPricesYearsAndCPI(t *testing.T) {
	// Create a sample GasolineData for testing
	sampleData := []cuda.GasolineData{
		{Year: "2020", AverageGasolinePrices: "2.242", AverageAnnualCPIForGas: "194.130", GasPricesAdjustedForInfl: "$4.02"},
	}

	prices, years, cpiValues := cuda.ExtractPricesYearsAndCPI(sampleData)

	// Perform your assertions on prices, years, and cpiValues.
	if len(prices) != 1 || prices[0] != 2.242 {
		t.Errorf("Expected prices to be [2.242], got %v", prices)
	}

	if len(years) != 1 || years[0] != 2020 {
		t.Errorf("Expected years to be [2020], got %v", years)
	}

	if len(cpiValues) != 1 || cpiValues[0] != 194.13 {
		t.Errorf("Expected cpiValues to be [194.13], got %v", cpiValues)
	}
}

// TestCreateScatterPlot tests the createScatterPlot function.
func TestCreateScatterPlot(t *testing.T) {
	// Sample data
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 5, 4, 5}

	// Regression line coefficients
	a := 0.6
	b := 2.2

	// Plot information
	title := "Scatter Plot with Regression Line"
	filename := "scatter_plot_test.png"
	xLabel := "X Axis Label"
	yLabel := "Y Axis Label"

	// Call the function
	cuda.CreateScatterPlot(x, y, a, b, title, filename, xLabel, yLabel)

	// Assertions
	assert.FileExists(t, filename, "The plot file should exist.")
	fileInfo, err := os.Stat(filename)
	assert.NoError(t, err, "Error getting file information.")
	assert.True(t, fileInfo.Size() > 0, "The plot file should not be empty.")
}
