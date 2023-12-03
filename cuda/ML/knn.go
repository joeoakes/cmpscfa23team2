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

//// AirfareData represents the structure of each item in the JSON data for airfare
//type AirfareData struct {
//	Title          string            `json:"title"`
//	Year           string            `json:"year"`
//	Location       string            `json:"location"`
//	Features       []string          `json:"features"`
//	AdditionalInfo AirfareAdditional `json:"additional_info"`
//	Metadata       AirfareMetadata   `json:"metadata"`
//}
//
//// AirfareDataList represents a list of airfare data
//type AirfareDataList struct {
//	AirfareData []AirfareData `json:"airfare_data"`
//}
//
//// AirfareAdditional represents additional information for airfare data
//type AirfareAdditional struct {
//	Country    string         `json:"country"`
//	MonthsData []AirfareMonth `json:"months_data"`
//}
//
//// AirfareMonth represents each month's data for airfare
//type AirfareMonth struct {
//	Month string `json:"month"`
//	Rate  string `json:"rate"`
//	Year  string `json:"year"`
//}
//
//// AirfareMetadata represents metadata for airfare data
//type AirfareMetadata struct {
//	Source    string `json:"source"`
//	Timestamp string `json:"timestamp"`
//}

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

//// ConvertAirfareDataToPoints converts airfare data to points
//func ConvertAirfareDataToPoints(airfareData []AirfareData) []Point {
//	var data []Point
//	for _, entry := range airfareData {
//		features := []float64{
//			parseFloat(entry.Year),
//			parseFloat(entry.Location),
//		}
//
//		label := "airfare"
//		data = append(data, Point{Features: features, Label: label})
//	}
//	return data
//}

// ConvertGasolineDataToPoints converts gasoline data to points
func ConvertGasolineDataToPoints(gasolineData []GasolineData) []Point {
	var data []Point
	for _, entry := range gasolineData {
		features := []float64{
			parseFloat(entry.Year),
			parseFloat(entry.AverageGasolinePrices),
			parseFloat(entry.AverageAnnualCPIForGas),
		}

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
			// Add two more features (assuming they are available in your dataset)
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
func createScatterPlot(data []Point, target Point, predictedLabel, xLabel, yLabel string) error {
	p := plot.New()

	p.Title.Text = "KNN Scatter Plot"
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
	//// Load airfare data
	//airfareFilePath := "airfare_data_price.json"
	//airfareDataList := loadAirfareData(airfareFilePath)
	//airfareData := airfareDataList.AirfareData

	// Load gasoline data
	gasolineFilePath := "gasoline_data.json"
	gasolineData := loadGasolineData(gasolineFilePath)

	// Load book data
	booksFilePath := "books_data.json"
	bookData := loadBookData(booksFilePath)

	// Convert data to points
	//airfarePoints := ConvertAirfareDataToPoints(airfareData)
	gasolinePoints := ConvertGasolineDataToPoints(gasolineData)
	bookPoints := ConvertBookDataToPoints(bookData)

	// Example values for Year, Location, Average Gasoline Prices
	target := Point{
		// Use the same number of features as other cases
		Features: []float64{2023, 0, 0},
		// The label will be predicted
		Label: "",
	}

	// Check the number of features in the target
	targetFeaturesCount := len(target.Features)

	//// Check the number of features in airfarePoints
	//for _, point := range airfarePoints {
	//	if len(point.Features) != targetFeaturesCount {
	//		fmt.Println("Error: Features lengths mismatch in airfarePoints")
	//		return
	//	}
	//}

	// Check the number of features in gasolinePoints
	for _, point := range gasolinePoints {
		if len(point.Features) != targetFeaturesCount {
			fmt.Println("Error: Features lengths mismatch in gasolinePoints")
			return
		}
	}

	// Check the number of features in bookPoints
	for _, point := range bookPoints {
		if len(point.Features) != targetFeaturesCount {
			fmt.Println("Error: Features lengths mismatch in bookPoints")
			return
		}
	}

	// Perform KNN predictions for each domain
	k := 1 // Number of neighbors to consider
	//predictedLabelAirfare := KNN(k, airfarePoints, target)
	predictedLabelGasoline := KNN(k, gasolinePoints, target)
	predictedLabelBooks := KNN(k, bookPoints, target)

	//fmt.Printf("Predicted Label for Airfare: %s\n", predictedLabelAirfare)
	fmt.Printf("Predicted Label for Gasoline: %s\n", predictedLabelGasoline)
	fmt.Printf("Predicted Label for Books: %s\n", predictedLabelBooks)

	// Create scatter plots
	//err := createScatterPlot(airfarePoints, target, predictedLabelAirfare, "Year", "Location")
	//if err != nil {
	//	log.Fatal(err)
	//}

	err := createScatterPlot(gasolinePoints, target, predictedLabelGasoline, "Year", "Average Gasoline Prices")
	if err != nil {
		log.Fatal(err)
	}

	err = createScatterPlot(bookPoints, target, predictedLabelBooks, "Timestamp", "Price")
	if err != nil {
		log.Fatal(err)
	}
}

//// Function to load airfare data
//func loadAirfareData(filePath string) AirfareDataList {
//	file, err := os.Open(filePath)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//
//	var data AirfareDataList
//	decoder := json.NewDecoder(file)
//	err = decoder.Decode(&data)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return data
//}

func loadGasolineData(filePath string) []GasolineData {
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

func loadBookData(filePath string) BookData {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data BookData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
