package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// AirfareData represents the structure of each item in the JSON data for airfare
type AirfareData struct {
	Title          string            `json:"title"`
	Year           string            `json:"year"`
	Location       string            `json:"location"`
	Features       []string          `json:"features"`
	AdditionalInfo AirfareAdditional `json:"additional_info"`
	Metadata       AirfareMetadata   `json:"metadata"`
}

// AirfareDataList represents a list of airfare data
type AirfareDataList struct {
	AirfareData []AirfareData `json:"airfare_data_price"`
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

// GasolineData represents the structure of each item in the JSON data for gas
type GasolineData struct {
	Year                     string `json:"year"`
	AverageGasolinePrices    string `json:"average_gasoline_prices"`
	AverageAnnualCPIForGas   string `json:"average_annual_cpi_for_gas"`
	GasPricesAdjustedForInfl string `json:"gas_prices_adjusted_for_inflation"`
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

type BookData struct {
	Domain string `json:"domain"`
	Data   []Item `json:"data"`
}

// Point represents a data point in 2D space
type Point struct {
	Features []float64
	Label    string
}

// EuclideanDistance computes the Euclidean distance between two points
func EuclideanDistance(a, b Point) float64 {
	if len(a.Features) != len(b.Features) {
		fmt.Println("Error: Features lengths mismatch")
		return 0.0
	}
	sum := 0.0
	for i := range a.Features {
		diff := a.Features[i] - b.Features[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

// ByDistance implements sort.Interface for []Point based on the distance to a fixed point
type ByDistance struct {
	Points   []Point
	Target   Point
	DistFunc func(Point, Point) float64
}

func (a ByDistance) Len() int      { return len(a.Points) }
func (a ByDistance) Swap(i, j int) { a.Points[i], a.Points[j] = a.Points[j], a.Points[i] }
func (a ByDistance) Less(i, j int) bool {
	return a.DistFunc(a.Target, a.Points[i]) < a.DistFunc(a.Target, a.Points[j])
}

// KNN performs the k-nearest neighbor classification
func KNN(k int, data []Point, target Point) string {
	sort.Sort(ByDistance{Points: data, Target: target, DistFunc: EuclideanDistance})

	// Ensure there are enough points in the data slice
	if len(data) < k {
		fmt.Println("Error: Not enough data points for k-nearest neighbors")
		return ""
	}

	labelVotes := make(map[string]int)

	// Access the first k elements of the sorted data slice
	for _, p := range data[:k] {
		labelVotes[p.Label]++
	}

	maxVotes := 0
	var predictedLabel string
	for label, votes := range labelVotes {
		if votes > maxVotes {
			maxVotes = votes
			predictedLabel = label
		}
	}

	return predictedLabel
}

// ConvertAirfareDataToPoints converts airfare data to points
func ConvertAirfareDataToPoints(airfareData AirfareData) []Point {
	var data []Point
	for _, monthData := range airfareData.AdditionalInfo.MonthsData {
		// Assuming the structure of AirfareMonth, adjust the field names accordingly
		features := []float64{
			parseFloat(monthData.Year),
			parseFloat(monthData.Rate),
			// Add more features as needed
		}

		label := "airfare"
		data = append(data, Point{Features: features, Label: label})
	}
	return data
}

// ConvertGasolineDataToPoints converts gasoline data to points
func ConvertGasolineDataToPoints(gasolineData []GasolineData) []Point {
	var data []Point
	for _, entry := range gasolineData {
		// Check for missing values and handle them appropriately
		year := parseFloat(entry.Year)
		avgGas := parseFloat(entry.AverageGasolinePrices)
		avgCPI := parseFloat(entry.AverageAnnualCPIForGas)

		// Skip the entry if any of the required values is missing
		if math.IsNaN(year) || math.IsNaN(avgGas) || math.IsNaN(avgCPI) {
			continue
		}

		features := []float64{year, avgGas, avgCPI}
		label := "gas"
		data = append(data, Point{Features: features, Label: label})
	}
	return data
}

// ConvertBookDataToPoints converts book data to points
func ConvertBookDataToPoints(bookData BookData) []Point {
	var data []Point
	for _, entry := range bookData.Data {
		// Exclude timestamp from features
		features := []float64{
			parseFloat(entry.Price),
			0.0, // Placeholder for the second feature
			0.0, // Placeholder for the third feature
		}

		label := "books"
		data = append(data, Point{Features: features, Label: label})
	}
	return data
}

// parseFloat parses a string to float64, handling special cases
func parseFloat(s string) float64 {
	// Remove currency symbols, if any
	s = strings.TrimPrefix(s, "£")

	// Parse float
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// Log an error message and return 0.0 for non-numeric values
		fmt.Printf("Error parsing float: %v\n", err)
		return 0.0
	}
	return value
}

var colorMap = map[string]color.Color{
	"airfare": color.RGBA{R: 0, G: 0, B: 255, A: 255}, // Blue
	"gas":     color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Red
	"books":   color.RGBA{R: 0, G: 255, B: 0, A: 255}, // Green
}

// createScatterPlot function with corrections
func createScatterPlot(data []Point, target Point, predictedLabel string) error {
	p := plot.New()

	p.Title.Text = "KNN Scatter Plot"

	// Set suitable x and y labels based on the selected dataset
	var xLabel, yLabel string
	switch predictedLabel {
	case "gas":
		xLabel = "Year"
		yLabel = "Average Gasoline Prices"

	case "books":
		xLabel = "Book Index"
		yLabel = "Price" // Change this to the appropriate label for books

	case "airfare":
		xLabel = "Month"
		yLabel = "Average Airfare Prices" // Change this to the appropriate label for airfare

	default:
		log.Fatal("Invalid dataset choice")
	}
	// set the minimum value of the y-axis to 0
	p.Y.Min = 0.0

	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	// Filter data based on the predicted label
	var filteredData []Point
	for _, point := range data {
		if point.Label == predictedLabel {
			filteredData = append(filteredData, point)
		}
	}

	// Plot the filtered data
	pts := make(plotter.XYs, len(filteredData))
	for i, point := range filteredData {
		pts[i].X = point.Features[0]
		pts[i].Y = point.Features[1]
	}

	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		return err
	}

	scatter.Color = colorMap[predictedLabel]
	p.Add(scatter)

	// Plot the target point with the predicted label
	pts = make(plotter.XYs, 1)
	pts[0].X = target.Features[0]
	pts[0].Y = target.Features[1]

	fmt.Printf("Target Point: X=%v, Y=%v\n", pts[0].X, pts[0].Y)

	scatter, err = plotter.NewScatter(pts)
	if err != nil {
		return err
	}

	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	scatter.GlyphStyle.Radius = vg.Points(3) // Set the radius to a suitable value
	scatter.Color = colorMap[predictedLabel] // Use the predicted label's color
	p.Add(scatter)

	// Save the plot to a PNG file
	filename := fmt.Sprintf("%s_scatter_plot.png", strings.ToLower(predictedLabel))
	if err := p.Save(6*vg.Inch, 4*vg.Inch, filename); err != nil {
		return err
	}

	fmt.Printf("Scatter plot created and saved as '%s'\n", filename)
	return nil
}

func main() {
	// Get user input for selecting the dataset
	var selectedDataset string
	fmt.Print("Choose a dataset (gas, books, airfare): ")
	fmt.Scanln(&selectedDataset)

	// Load the selected dataset
	var allPoints []Point
	switch selectedDataset {
	case "gas":
		gasolineDataFile, err := os.ReadFile("gasoline_data.json")
		if err != nil {
			log.Fatal(err)
		}
		var gasolineData []GasolineData
		err = json.Unmarshal(gasolineDataFile, &gasolineData)
		if err != nil {
			log.Fatal(err)
		}
		allPoints = ConvertGasolineDataToPoints(gasolineData)

	case "books":
		bookDataFile, err := os.ReadFile("books_data.json")
		if err != nil {
			log.Fatal(err)
		}
		var bookData BookData
		err = json.Unmarshal(bookDataFile, &bookData)
		if err != nil {
			log.Fatal(err)
		}
		allPoints = ConvertBookDataToPoints(bookData)

	case "airfare":
		airfareDataFile, err := os.ReadFile("airfare_data_price.json")
		if err != nil {
			log.Fatal(err)
		}
		var airfareDataList AirfareDataList
		err = json.Unmarshal(airfareDataFile, &airfareDataList)
		if err != nil {
			log.Fatal(err)
		}
		for _, airfareData := range airfareDataList.AirfareData {
			allPoints = append(allPoints, ConvertAirfareDataToPoints(airfareData)...)
		}

	default:
		log.Fatal("Invalid dataset choice")
	}

	// Example usage of KNN
	target := Point{
		Features: make([]float64, len(allPoints[0].Features)),
	}

	switch selectedDataset {
	case "gas":
		target.Features[2] = 2023

	case "books":
		target.Features[0] = 100

	case "airfare":
		target.Features[0] = 2023

	default:
		log.Fatal("Invalid dataset choice")
	}

	k := 5
	predictedLabel := KNN(k, allPoints, target)

	// Calculate the dynamic middle value for the second feature based on the predicted label
	for i := range target.Features {
		var sumFeature float64
		count := 0
		for _, p := range allPoints[:k] {
			if p.Label == predictedLabel {
				sumFeature += p.Features[i]
				count++
			}
		}
		if count > 0 {
			target.Features[i] = sumFeature / float64(count)
		}
	}

	fmt.Printf("Predicted Label: %s\n", predictedLabel)
	fmt.Printf("Predicted Features: %v\n", target.Features)

	// Create and save scatter plot
	err := createScatterPlot(allPoints, target, predictedLabel)
	if err != nil {
		log.Fatal(err)
	}

}
