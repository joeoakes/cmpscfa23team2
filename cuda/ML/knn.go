package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"sort"
	"time"
)

// WeatherData Hard coded weather data
type WeatherData struct {
	Location    string
	Date        string
	Temperature float64
}
type Item struct {
	Domain string `json:"domain"`
	Data   struct {
		Title       string   `json:"title"`
		URL         string   `json:"url"`
		Description string   `json:"description"`
		Price       string   `json:"price"`
		Location    string   `json:"location"`
		Features    []string `json:"features"`
		Reviews     []struct {
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

func main1() {
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
	//target := Point{Features: []float64{79, 12}} // replace with relevant features for target
	k := 1
	label := KNN(k, data, target)
	fmt.Printf("Label predicted for target is '%s'\n", label)

}
