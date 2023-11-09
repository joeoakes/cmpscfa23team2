package main

import (
	"math"
	"strings"
)

// JobListing represents a job listing with its description and industry.
type JobListing struct {
	Description string
	Industry    string
}

// NaiveBayesClassifier struct to hold model data.
type NaiveBayesClassifier struct {
	wordFrequencies map[string]map[string]int
	industryCounts  map[string]int
	totalWords      int
}

// NewNaiveBayesClassifier creates a new Naive Bayes Classifier.
func NewNaiveBayesClassifier() *NaiveBayesClassifier {
	return &NaiveBayesClassifier{
		wordFrequencies: make(map[string]map[string]int),
		industryCounts:  make(map[string]int),
		totalWords:      0,
	}
}

// Train trains the Naive Bayes Classifier with a slice of JobListings.
func (nbc *NaiveBayesClassifier) Train(data []JobListing) {
	for _, listing := range data {
		words := strings.Fields(listing.Description)
		industry := listing.Industry

		if nbc.wordFrequencies[industry] == nil {
			nbc.wordFrequencies[industry] = make(map[string]int)
		}

		for _, word := range words {
			nbc.wordFrequencies[industry][word]++
			nbc.totalWords++
		}

		nbc.industryCounts[industry]++
	}
}

// calculateProbability calculates the probability of a set of words belonging to a given industry.
func (nbc *NaiveBayesClassifier) calculateProbability(words []string, industry string) float64 {
	vocabSize := len(nbc.wordFrequencies) // Assuming each industry's map contains unique words only
	prob := math.Log(float64(nbc.industryCounts[industry]) / float64(len(nbc.industryCounts)))

	for _, word := range words {
		wordFrequency := nbc.wordFrequencies[industry][word]
		prob += math.Log(float64(wordFrequency+1) / float64(nbc.totalWords+vocabSize))
	}

	return prob
}

// Predict predicts the industry for a given job listing description.
func (nbc *NaiveBayesClassifier) Predict(description string) string {
	words := strings.Fields(description)
	highestProb := math.Inf(-1) // Set to negative infinity
	predictedIndustry := ""

	for industry := range nbc.industryCounts {
		prob := nbc.calculateProbability(words, industry)
		if prob > highestProb {
			highestProb = prob
			predictedIndustry = industry
		}
	}

	return predictedIndustry
}

func main() {
	// Example usage
	data := []JobListing{
		{"Software development in Go", "Technology"},
		{"Managing financial assets", "Finance"},
		// Add more job listings for training
	}

	classifier := NewNaiveBayesClassifier()
	classifier.Train(data)

	newJobListing := "Programming in Python and Go"
	prediction := classifier.Predict(newJobListing)
	println("Predicted Industry:", prediction)
}
