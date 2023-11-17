package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

// WeatherData Hard coded weather data
type WeatherData struct {
	Location    string
	Date        string
	Temperature float64
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

// reads data from a csv file. Change to however you want it to be or have it do
func ReadDataFromFile(filename string) ([]Point, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var data []Point
	for _, line := range lines {
		var features []float64
		for _, value := range line[:len(line)-1] {
			feature, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			features = append(features, feature)
		}
		label := line[len(line)-1]

		data = append(data, Point{Features: features, Label: label})
	}
	return data, nil
}

// write the data to csv file
func WriteDatatoFile(filename string, data []Point) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, point := range data {
		var record []string
		for _, feature := range point.Features {
			record = append(record, strconv.FormatFloat(feature, 'f', -1, 64))

		}
		record = append(record, point.Label)
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}
	return nil
}

// added this here just in case we need it to read a csv file

func main() {
	//filename := "your_data.csv"
	//data, err := ReadDataFromFile(filename)
	//if err != nil {
	//	fmt.Printf("Error reading data from file: %v\n", err)
	//	return
	//}
	// calling function to read the csv file
	// but unsure of whether it will be used in the future
	// so leave it there - Binh

	// KNN with hard coded weather data
	data := []Point{
		{Features: []float64{1, 2}, Label: "Sunny"},
		{Features: []float64{3, 1}, Label: "Rainy"},
		{Features: []float64{2, 4}, Label: "Cloudy"},
		{Features: []float64{5, 3}, Label: "Sunny"},
	}
	err := WriteDatatoFile("dal/data.csv", data)

	// Target point to classify (representing new weather data)
	target := Point{Features: []float64{5, 1}}
	Data, err := ReadDataFromFile("dal/data.csv")
	if err != nil {
		fmt.Printf("Error reading data from csv file: %v\n", err)
		return
	}
	storedData := append(Data)

	// Perform KNN classification
	k := 1 // Number of neighbors
	label := KNN(k, storedData, target)
	fmt.Printf("The label predicted for the target is '%s'\n", label)

	err = WriteDatatoFile("dal/data_updated.csv", storedData)
	if err != nil {
		fmt.Printf("Error writing data to csv file: %v\n", err)
		return

	}
	fmt.Printf("Updated data written to dal/data_updated.csv")

}
