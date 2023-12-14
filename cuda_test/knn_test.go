package cuda_test

import (
	"cmpscfa23team2/cuda/ML"
	"math"
	"testing"
)

func TestKNN(t *testing.T) {
	// Define test cases
	tests := []struct {
		k        int
		data     []cuda.Point
		target   cuda.Point
		expected string
	}{
		{
			k: 3,
			data: []cuda.Point{
				{Features: []float64{1, 2}, Label: "A"},
				{Features: []float64{2, 3}, Label: "A"},
				{Features: []float64{3, 4}, Label: "B"},
			},
			target:   cuda.Point{Features: []float64{2, 2}},
			expected: "A",
		},
	}

	for _, test := range tests {
		result, _ := cuda.KNN(test.k, test.data, test.target)
		if result != test.expected {
			t.Errorf("KNN(...) = %v, want %v", result, test.expected)
		}
	}
}

func TestFloatMonth(t *testing.T) {
	tests := []struct {
		month    string
		expected float64
	}{
		{"Jan", 1.0},
		{"Feb", 2.0},
	}

	for _, test := range tests {
		result := cuda.FloatMonth(test.month)
		if result != test.expected {
			t.Errorf("floatMonth(%v) = %v, want %v", test.month, result, test.expected)
		}
	}
}

func TestEuclideanDistance(t *testing.T) {
	// Define test cases
	tests := []struct {
		a, b     cuda.Point
		expected float64
	}{
		{
			a:        cuda.Point{Features: []float64{0, 0}},
			b:        cuda.Point{Features: []float64{3, 4}},
			expected: 5.0,
		},
	}

	for _, test := range tests {
		result := cuda.EuclideanDistance(test.a, test.b)
		if result != test.expected {
			t.Errorf("EuclideanDistance(%v, %v) = %v, want %v", test.a, test.b, result, test.expected)
		}
	}
}

func TestParseFloat(t *testing.T) {

	// Test cases for parseFloat

	var tests []struct {
		input string

		expected float64
	}

	for _, test := range tests {

		result := cuda.ParseFloat(test.input)

		if result != test.expected {

			t.Errorf("parseFloat(%v) = %v, want %v", test.input, result, test.expected)

		}

	}

}

func TestConvertGasolineDataToPoints(t *testing.T) {
	// Define test cases
	tests := []struct {
		data     []cuda.GasolineData
		expected []cuda.Point
	}{
		{
			data: []cuda.GasolineData{
				{Year: "2020", AverageGasolinePrices: "2.5", AverageAnnualCPIForGas: "1.2"},
				{Year: "2021", AverageGasolinePrices: "2.8", AverageAnnualCPIForGas: "1.4"},
				{Year: "2022", AverageGasolinePrices: "3.0", AverageAnnualCPIForGas: "1.5"},
			},
			expected: []cuda.Point{
				{Features: []float64{2020, 2.5, 1.2}, Label: "gas"},
				{Features: []float64{2021, 2.8, 1.4}, Label: "gas"},
				{Features: []float64{2022, 3.0, 1.5}, Label: "gas"},
			},
		},
	}

	for _, test := range tests {
		result := cuda.ConvertGasolineDataToPoints(test.data)

		// Check the length of the result
		if len(result) != len(test.expected) {
			t.Errorf("ConvertGasolineDataToPoints(...) produced incorrect length. Got %d, want %d", len(result), len(test.expected))
		}

		// Check each point in the result
		for i := range result {
			// Check label
			if result[i].Label != test.expected[i].Label {
				t.Errorf("ConvertGasolineDataToPoints(...) produced incorrect label. Got %s, want %s", result[i].Label, test.expected[i].Label)
			}

			// Check features
			for j := range result[i].Features {
				if !almostEqual(result[i].Features[j], test.expected[i].Features[j]) {
					t.Errorf("ConvertGasolineDataToPoints(...) produced incorrect feature. Got %v, want %v", result[i].Features, test.expected[i].Features)
				}
			}
		}
	}
}

func TestConvertBookDataToPoints(t *testing.T) {
	// Define test cases
	tests := []struct {
		data     cuda.BookData
		expected []cuda.Point
	}{
		{
			data: cuda.BookData{
				Domain: "books",
				Data: []cuda.Item{
					{Price: "15.5"},
					{Price: "20.0"},
					{Price: "18.3"},
				},
			},
			expected: []cuda.Point{
				{Features: []float64{15.5}, Label: "books"},
				{Features: []float64{20.0}, Label: "books"},
				{Features: []float64{18.3}, Label: "books"},
			},
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result := cuda.ConvertBookDataToPoints(test.data)

		// Check the length of the result
		if len(result) != len(test.expected) {
			t.Errorf("ConvertBookDataToPoints(...) produced incorrect length. Got %d, want %d", len(result), len(test.expected))
		}

		// Check each point in the result
		for i := range result {
			// Check label
			if result[i].Label != test.expected[i].Label {
				t.Errorf("ConvertBookDataToPoints(...) produced incorrect label. Got %s, want %s", result[i].Label, test.expected[i].Label)
			}

			// Check features
			for j := range result[i].Features {
				if !almostEqual(result[i].Features[j], test.expected[i].Features[j]) {
					t.Errorf("ConvertBookDataToPoints(...) produced incorrect feature. Got %v, want %v", result[i].Features, test.expected[i].Features)
				}
			}
		}
	}
}

// almostEqual checks if two float64 values are almost equal within a small tolerance
func almostEqual(a, b float64) bool {
	const tolerance = 1e-9
	return math.Abs(a-b) < tolerance
}
