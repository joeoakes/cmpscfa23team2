package main

import (
	"fmt"
	"math"
	"sort"
)

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

func main2() {
	// Example dataset
	data := []Point{
		{Features: []float64{1, 2}, Label: "A"},
		{Features: []float64{3, 1}, Label: "B"},
		{Features: []float64{2, 4}, Label: "C"},
		{Features: []float64{5, 3}, Label: "D"},
		// Add more data points if needed
	}

	// Target point to classify
	target := Point{Features: []float64{5, 1}}

	// Perform kNN classification
	k := 1 // Number of neighbors
	label := KNN(k, data, target)
	fmt.Printf("The label predicted for the target is '%s'\n", label)
}
