package main

import (
	"encoding/json"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"io/ioutil"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Item struct {
	Domain string `json:"domain"`
	Data   struct {
		Title       string `json:"title"`
		URL         string `json:"url"`
		Description string `json:"description"`
		Price       string `json:"price"`
		Location    string `json:"location"`

		Features []string `json:"features"`
		Reviews  []struct {
			User    string `json:"user"`
			Rating  int    `json:"rating"`
			Comment string `json:"comment"`
		}
		Images         []string          `json:"images"`
		AdditionalInfo map[string]string `json:"additional_info"`
		Metadata       struct {
			Source    string `json:"source"`
			Timestamp string `json:"timestamp"`
		} `json:"metadata"`
	} `json:"data"`
}

// Point represents a data point in 2D space
type Point struct {
	Features []float64
	Label    string
}

// EuclideanDistance computes the Euclidean distance between two points
func EuclideanDistance(a, b Point) float64 {
	// error handling for trying to access an element of a slice that doesn't exist
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
	// additional error checking
	if a.DistFunc == nil {
		fmt.Println("Error: DistFunc is nil")
		return false
	}

	return a.DistFunc(a.Target, a.Points[i]) < a.DistFunc(a.Target, a.Points[j])
}

// KNN performs the k-nearest neighbor classification
func KNN(k int, data []Point, target Point) string {
	// Sort the points by distance from target
	sort.Sort(ByDistance{Points: data, Target: target, DistFunc: EuclideanDistance})

	// Map to keep track of the frequency of each label among the k closest points
	labelVotes := make(map[string]int)

	// Tally votes from the k nearest neighbors
	for _, p := range data[:k] {
		labelVotes[p.Label]++
	}

	// Find the label with the most votes
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

func ConvertItemsToPoints(items []Item) []Point {
	var data []Point
	for _, item := range items {
		var features []float64

		switch item.Domain {
		case "e-commerce":
			features = []float64{float64(len(item.Data.Description)), float64(len(item.Data.Features))}
		case "real-estate":
			features = []float64{float64(len(item.Data.Description)), float64(countSubstring(item.Data.Features, "Bedrooms"))}
		case "job-market":
			features = []float64{float64(len(item.Data.Description)), float64(len(item.Data.Features))}
		}

		// pad or truncate features to ensure fixed length
		for len(features) < NumFeatures {
			features = append(features, 0.0)
		}
		features = features[:NumFeatures]

		label := item.Domain
		data = append(data, Point{Features: features, Label: label})

	}
	return data
}
func countSubstring(slice []string, substring string) int {
	count := 0
	for _, s := range slice {
		if s == substring {
			count++
		}
	}
	return count
}

// createScatterPlot creates and saves a scatter plot with highlighted target and predicted points.
func createScatterPlot(data []Point, target Point, predictedLabel, filename, xLabel, yLabel string) error {
	p := plot.New()

	p.Title.Text = "KNN Scatter Plot"
	p.Title.TextStyle.Font.Size = 16 // Set title font size

	p.X.Label.Text = xLabel
	p.X.Label.TextStyle.Font.Size = 16 // Set X-axis label font size
	p.X.Tick.Label.Font.Size = 16      // Set X-axis tick label font size

	p.Y.Label.Text = yLabel
	p.Y.Label.TextStyle.Font.Size = 16 // Set Y-axis label font size
	p.Y.Tick.Label.Font.Size = 16      // Set Y-axis tick label font size

	// Create a map to store colors for each domain
	colorMap := map[string]color.RGBA{
		"e-commerce":  color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Red
		"real-estate": color.RGBA{R: 0, G: 255, B: 0, A: 255}, // Green
		"job-market":  color.RGBA{R: 0, G: 0, B: 255, A: 255}, // Blue
		"target":      color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Target color (Red)
		"predicted":   color.RGBA{R: 0, G: 0, B: 255, A: 255}, // Predicted color (Blue)
	}

	// Create points for the scatter plot
	for _, point := range data {
		// Create scatter plot points with different colors for different domains
		s, err := plotter.NewScatter(plotter.XYs{{X: point.Features[0], Y: point.Features[1]}})
		if err != nil {
			fmt.Printf("Error creating scatter plot: %v\n", err)
			return err
		}

		s.GlyphStyle.Color = colorMap[point.Label]
		s.GlyphStyle.Radius = vg.Points(10) // Adjust the radius as needed

		p.Add(s)
	}

	// Highlight the target and predicted points
	targetScatter, err := plotter.NewScatter(plotter.XYs{{X: target.Features[0], Y: target.Features[1]}})
	if err != nil {
		fmt.Printf("Error creating target scatter plot: %v\n", err)
		return err
	}
	targetScatter.GlyphStyle.Color = colorMap["target"]
	p.Add(targetScatter)

	predictedScatter, err := plotter.NewScatter(plotter.XYs{{X: target.Features[0], Y: target.Features[1]}})
	if err != nil {
		fmt.Printf("Error creating predicted scatter plot: %v\n", err)
		return err
	}
	predictedScatter.GlyphStyle.Color = colorMap["predicted"]
	p.Add(predictedScatter)

	// Save the plot to a file
	if err := p.Save(16*vg.Inch, 8*vg.Inch, filename); err != nil {
		fmt.Printf("Error saving scatter plot: %v\n", err)
		return err
	}

	return nil
}

// main function
func main() {
	// reading data from JSON file
	file, err := ioutil.ReadFile("crab/template.json")
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	// turn JSON data into items
	var jsonData map[string][]Item
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}
	items, ok := jsonData["items"]
	if !ok {
		fmt.Printf("Error: 'items' field not found in JSON.\n")
		return
	}

	// convert items to points
	data := ConvertItemsToPoints(items)
	// seed random number generator to make it truly random
	rand.Seed(time.Now().UnixNano())

	// randomly choose a target point from the dataset
	targetIndex := rand.Intn(len(items))
	targetItem := items[targetIndex]
	target := ConvertItemsToPoints([]Item{targetItem})[0]

	// target point to classify
	k := 1
	label := KNN(k, data, target)
	fmt.Printf("Label predicted for target is '%s'\n", label)

	// creating and saving the scatter plot
	err = createScatterPlot(data, target, label, "scatter_plot.png", "Feature 1", "Feature 2")
	if err != nil {
		fmt.Println("Error creating scatter plot:", err)
		return
	}
}
