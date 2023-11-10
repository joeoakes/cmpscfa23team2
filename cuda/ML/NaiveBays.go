package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// JobListing represents a job listing with its description and industry.
type JobListing struct {
	Description string
	Industry    string
}

// NaiveBayesClassifier struct to hold model data.
type NaiveBayesClassifier struct {
	wordFrequencies  map[string]map[string]int
	industryCounts   map[string]int
	totalWords       int
	totalUniqueWords int // New field to track unique words
}

// NewNaiveBayesClassifier creates a new Naive Bayes Classifier.
func NewNaiveBayesClassifier() *NaiveBayesClassifier {
	return &NaiveBayesClassifier{
		wordFrequencies: make(map[string]map[string]int),
		industryCounts:  make(map[string]int),
		totalWords:      0,
	}
}

var stopWords = map[string]bool{
	"and": true, "or": true, "the": true, "in": true,
	"of": true, "a": true, "is": true, "to": true,
	// More stop words can be added here
}

// preprocessText preprocesses the text by converting it to lower case and removing stop words.
func preprocessText(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	var processedWords []string
	for _, word := range words {
		if _, ok := stopWords[word]; !ok {
			processedWords = append(processedWords, word)
		}
	}
	return processedWords
}

func (nbc *NaiveBayesClassifier) Train(data []JobListing) {
	uniqueWords := make(map[string]bool)

	for _, listing := range data {
		words := preprocessText(listing.Description)
		industry := listing.Industry

		if nbc.wordFrequencies[industry] == nil {
			nbc.wordFrequencies[industry] = make(map[string]int)
		}

		for _, word := range words {
			if !uniqueWords[word] {
				uniqueWords[word] = true
				nbc.totalUniqueWords++
			}

			nbc.wordFrequencies[industry][word]++
			nbc.totalWords++
		}

		nbc.industryCounts[industry]++
	}
}

// calculateProbability calculates the probability of a set of words belonging to a given industry.
func (nbc *NaiveBayesClassifier) calculateProbability(words []string, industry string) float64 {
	vocabSize := nbc.totalUniqueWords
	prob := math.Log(float64(nbc.industryCounts[industry]) / float64(len(nbc.industryCounts)))

	for _, word := range words {
		wordFrequency := nbc.wordFrequencies[industry][word]
		prob += math.Log(float64(wordFrequency+1) / float64(nbc.totalWords+vocabSize))
	}

	return prob
}

// PredictWithProbabilities predicts the industry for a given job description and returns probabilities.
func (nbc *NaiveBayesClassifier) PredictWithProbabilities(description string) (string, map[string]float64) {
	words := preprocessText(description)
	industryProbabilities := make(map[string]float64)

	for industry := range nbc.industryCounts {
		prob := nbc.calculateProbability(words, industry)
		industryProbabilities[industry] = math.Exp(prob) // Convert log probability back
	}

	// Get the industry with the highest probability
	var topIndustry string
	maxProb := 0.0
	for industry, prob := range industryProbabilities {
		if prob > maxProb {
			maxProb = prob
			topIndustry = industry
		}
	}

	return topIndustry, industryProbabilities
}

// displayTopIndustries displays top 'n' industries based on probability.
func displayTopIndustries(probs map[string]float64, n int) {
	type industryProb struct {
		Industry string
		Prob     float64
	}
	var sortedProbs []industryProb
	for industry, prob := range probs {
		sortedProbs = append(sortedProbs, industryProb{industry, prob})
	}

	sort.Slice(sortedProbs, func(i, j int) bool {
		return sortedProbs[i].Prob > sortedProbs[j].Prob
	})

	for i := 0; i < n && i < len(sortedProbs); i++ {
		fmt.Printf("%d. %s: %e\n", i+1, sortedProbs[i].Industry, sortedProbs[i].Prob)
	}
}

func main() {
	data := []JobListing{
		// Agriculture
		{"Agricultural technology and farm management", "Agriculture"},
		{"Sustainable farming practices", "Agriculture"},

		// Aerospace
		{"Aerospace engineering and avionics", "Aerospace"},
		{"Spacecraft design and systems engineering", "Aerospace"},

		// Art
		{"Photography and visual arts", "Art"},
		{"Contemporary art curation and management", "Art"},

		// Automotive
		{"Automotive engineering and design", "Automotive"},
		{"Electric vehicle technology", "Automotive"},

		// Biotech
		{"Biotechnology research and genetic engineering", "Biotech"},
		{"Clinical trial management in biotech", "Biotech"},

		// Culinary
		{"Culinary arts and restaurant management", "Culinary"},
		{"Gourmet cuisine and fine dining", "Culinary"},

		// Design
		{"Graphic design and digital media creation", "Design"},
		{"Interior design and space planning", "Design"},

		// Education
		{"Educational curriculum design and teaching", "Education"},
		{"Educational technology and e-learning", "Education"},

		// Energy
		{"Renewable energy and green technologies", "Energy"},
		{"Oil and gas exploration and production", "Energy"},

		// Engineering
		{"Civil engineering and infrastructure development", "Engineering"},
		{"Mechanical engineering in manufacturing", "Engineering"},

		// Environmental
		{"Environmental policy and sustainability", "Environmental"},
		{"Conservation biology and ecosystem management", "Environmental"},

		// Event Management
		{"Event planning and management", "Event Management"},
		{"Corporate event coordination", "Event Management"},

		// Fashion
		{"Fashion design and apparel merchandising", "Fashion"},
		{"Trend forecasting in fashion", "Fashion"},

		// Finance
		{"Managing financial assets", "Finance"},
		{"Investment banking and asset management", "Finance"},

		// Fitness
		{"Fitness training and sports coaching", "Fitness"},
		{"Health and wellness coaching", "Fitness"},

		// Healthcare
		{"Nursing care and patient management", "Healthcare"},
		{"Clinical psychology and mental health counseling", "Healthcare"},

		// Hospitality
		{"Hospitality services and hotel management", "Hospitality"},
		{"Travel and tourism management", "Hospitality"},

		// HR
		{"Human resources policies and employee relations", "HR"},
		{"Talent acquisition and recruitment", "HR"},

		// Journalism
		{"Journalism and news reporting", "Journalism"},
		{"Investigative journalism and documentary production", "Journalism"},

		// Legal
		{"Legal compliance and corporate law", "Legal"},
		{"Intellectual property law and patent strategy", "Legal"},

		// Marketing
		{"Marketing strategy and brand development", "Marketing"},
		{"Digital marketing and social media strategy", "Marketing"},

		// Media
		{"Film production and media arts", "Media"},
		{"Broadcast journalism and television production", "Media"},

		// Operations
		{"Supply chain logistics and management", "Operations"},
		{"Business operations and process improvement", "Operations"},

		// Pharmaceutical
		{"Pharmaceutical research and development", "Pharmaceutical"},
		{"Regulatory affairs in pharmaceuticals", "Pharmaceutical"},

		// Politics
		{"International diplomacy and policy analysis", "Politics"},
		{"Political campaign management and strategy", "Politics"},

		// Public Relations
		{"Public relations and corporate communications", "Public Relations"},
		{"Crisis communication and reputation management", "Public Relations"},

		// Real Estate
		{"Real estate development and property management", "Real Estate"},
		{"Commercial real estate investment and brokerage", "Real Estate"},

		// Retail
		{"Retail sales management and merchandising", "Retail"},
		{"E-commerce and digital retail strategy", "Retail"},

		// Technology
		{"Software development in Go", "Technology"},
		{"Developing advanced machine learning models", "Technology"},

		// Veterinary
		{"Veterinary services and animal care", "Veterinary"},
		{"Animal nutrition and wellness", "Veterinary"},
	}

	classifier := NewNaiveBayesClassifier()
	classifier.Train(data)

	newJobListing := "Programming in Python and Go"
	predictedIndustry, probs := classifier.PredictWithProbabilities(newJobListing)

	fmt.Printf("Predicted Industry: %s\n", predictedIndustry)
	fmt.Println("Top Industries with Probabilities:")
	displayTopIndustries(probs, 3)

	// Test Case 1
	testJobDescription1 := "Designing user interfaces and user experiences"
	predictedIndustry1, _ := classifier.PredictWithProbabilities(testJobDescription1)
	fmt.Printf("Test 1 - Job Description: '%s'\nPredicted Industry: %s\n\n", testJobDescription1, predictedIndustry1)

	// Test Case 2
	testJobDescription2 := "Developing policies for environmental conservation"
	predictedIndustry2, _ := classifier.PredictWithProbabilities(testJobDescription2)
	fmt.Printf("Test 2 - Job Description: '%s'\nPredicted Industry: %s\n\n", testJobDescription2, predictedIndustry2)

	// Test Case 3
	testJobDescription3 := "Managing investment portfolios and financial risks"
	predictedIndustry3, _ := classifier.PredictWithProbabilities(testJobDescription3)
	fmt.Printf("Test 3 - Job Description: '%s'\nPredicted Industry: %s\n\n", testJobDescription3, predictedIndustry3)

	// Test Case 4
	testJobDescription4 := "Culinary techniques and kitchen management"
	predictedIndustry4, _ := classifier.PredictWithProbabilities(testJobDescription4)
	fmt.Printf("Test 4 - Job Description: '%s'\nPredicted Industry: %s\n", testJobDescription4, predictedIndustry4)
}
