package main

import (
	"encoding/json"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

type GasolineData struct {
	Year                          string `json:"year"`
	AverageGasolinePrices         string `json:"average_gasoline_prices"`
	AverageAnnualCPIForGas        string `json:"average_annual_cpi_for_gas"`
	GasPricesAdjustedForInflation string `json:"gas_prices_adjusted_for_inflation"`
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

	labelVotes := make(map[string]int)

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

// NumFeatures creating a fixed length for feature slices so there is no mismatch in EuclideanDistance
const NumFeatures = 10

func ConvertGasolineDataToPoints(gasolineData []GasolineData) []Point {
	var data []Point
	for _, entry := range gasolineData {
		features := []float64{
			parseFloat(entry.AverageGasolinePrices),
			parseFloat(entry.AverageAnnualCPIForGas),
		}

		for len(features) < NumFeatures {
			features = append(features, 0.0)
		}
		features = features[:NumFeatures]

		label := "gas-prices"
		data = append(data, Point{Features: features, Label: label})
	}
	return data
}

func parseFloat(s string) float64 {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Printf("Error parsing float: %v\n", err)
		return 0.0
	}
	return value
}

var colorMap = map[string]color.Color{
	"gas-prices": color.RGBA{R: 0, G: 0, B: 255, A: 255}, // Blue
	"target":     color.RGBA{R: 255, G: 165, A: 255},     // Orange
}

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
	for _, point := range filteredData {
		pts := make(plotter.XYs, 1)
		pts[0].X = point.Features[0]
		pts[0].Y = point.Features[1]

		s, err := plotter.NewScatter(pts)
		if err != nil {
			log.Panic(err)
		}
		s.GlyphStyle.Color = colorMap[point.Label]
		p.Add(s)
	}

	// Plot the target point
	targetPts := make(plotter.XYs, 1)
	targetPts[0].X = target.Features[0]
	targetPts[0].Y = target.Features[1]
	targetScatter, err := plotter.NewScatter(targetPts)
	if err != nil {
		log.Panic(err)
	}
	targetScatter.GlyphStyle.Shape = draw.PyramidGlyph{}
	targetScatter.GlyphStyle.Color = colorMap["target"]
	p.Add(targetScatter)

	// Add target point to the legend
	p.Legend.Add("target", targetScatter)

	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 6*vg.Inch, "knn_scatter_plot.png"); err != nil {
		log.Panic(err)
	}

	return nil
}

func calculateCentroid(points []Point) Point {
	if len(points) == 0 {
		return Point{}
	}

	numFeatures := len(points[0].Features)
	meanFeatures := make([]float64, numFeatures)
	for _, point := range points {
		for i, feature := range point.Features {
			meanFeatures[i] += feature
		}
	}

	for i := range meanFeatures {
		meanFeatures[i] /= float64(len(points))
	}

	return Point{Features: meanFeatures}
}

func main() {
	// Read data from JSON file
	file, err := os.ReadFile("gasoline_data.json")
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	// Unmarshal JSON data
	var gasolineData []GasolineData
	err = json.Unmarshal(file, &gasolineData)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	// Convert gasoline data to points
	data := ConvertGasolineDataToPoints(gasolineData)

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Choose a random target item
	targetIndex := rand.Intn(len(gasolineData))
	targetItem := gasolineData[targetIndex]
	target := ConvertGasolineDataToPoints([]GasolineData{targetItem})[0]

	// Number of nearest neighbors for KNN
	k := 1

	// Perform KNN on the data
	predictedLabel := KNN(k, data, target)
	fmt.Printf("Label predicted for target is '%s'\n", predictedLabel)

	// Create and save the scatter plot
	err = createScatterPlot(data, target, predictedLabel, "Average Gasoline Prices", "Average Annual CPI for Gas")
	if err != nil {
		fmt.Println("Error creating scatter plot:", err)
		return
	}
}
