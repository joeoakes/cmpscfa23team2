package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// GenericTextData represents any text data with an associated category.
type GenericTextData struct {
	Description string
	Category    string
}

// NaiveBayesClassifier struct to hold model data.
type NaiveBayesClassifier struct {
	wordFrequencies  map[string]map[string]int
	categoryCounts   map[string]int
	totalWords       int
	totalUniqueWords int
}

// NewNaiveBayesClassifier creates a new Naive Bayes Classifier.
func NewNaiveBayesClassifier() *NaiveBayesClassifier {
	return &NaiveBayesClassifier{
		wordFrequencies:  make(map[string]map[string]int),
		categoryCounts:   make(map[string]int),
		totalWords:       0,
		totalUniqueWords: 0,
	}
}

// stopWords is a map of common words to be excluded from the analysis.
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

// Train takes generic text data and trains the classifier.
func (nbc *NaiveBayesClassifier) Train(data []GenericTextData) {
	uniqueWords := make(map[string]bool)

	for _, item := range data {
		words := preprocessText(item.Description)
		category := item.Category

		if nbc.wordFrequencies[category] == nil {
			nbc.wordFrequencies[category] = make(map[string]int)
		}

		for _, word := range words {
			if !uniqueWords[word] {
				uniqueWords[word] = true
				nbc.totalUniqueWords++
			}

			nbc.wordFrequencies[category][word]++
			nbc.totalWords++
		}

		nbc.categoryCounts[category]++
	}
}

// calculateProbability calculates the probability of a set of words belonging to a given category.
func (nbc *NaiveBayesClassifier) calculateProbability(words []string, category string) float64 {
	vocabSize := nbc.totalUniqueWords
	prob := math.Log(float64(nbc.categoryCounts[category]) / float64(len(nbc.categoryCounts)))

	for _, word := range words {
		wordFrequency := nbc.wordFrequencies[category][word]
		prob += math.Log(float64(wordFrequency+1) / float64(nbc.totalWords+vocabSize))
	}

	return prob
}

// PredictWithProbabilities predicts the category for a given description and returns probabilities.
func (nbc *NaiveBayesClassifier) PredictWithProbabilities(description string) (string, map[string]float64) {
	words := preprocessText(description)
	categoryProbabilities := make(map[string]float64)

	for category := range nbc.categoryCounts {
		prob := nbc.calculateProbability(words, category)
		categoryProbabilities[category] = math.Exp(prob) // Convert log probability back
	}

	// Get the category with the highest probability
	var topCategory string
	maxProb := 0.0
	for category, prob := range categoryProbabilities {
		if prob > maxProb {
			maxProb = prob
			topCategory = category
		}
	}

	return topCategory, categoryProbabilities
}

// displayTopCategories displays top 'n' categories based on probability.
func displayTopCategories(probs map[string]float64, n int) {
	type categoryProb struct {
		Category string
		Prob     float64
	}
	var sortedProbs []categoryProb
	for category, prob := range probs {
		sortedProbs = append(sortedProbs, categoryProb{category, prob})
	}

	sort.Slice(sortedProbs, func(i, j int) bool {
		return sortedProbs[i].Prob > sortedProbs[j].Prob
	})

	for i := 0; i < n && i < len(sortedProbs); i++ {
		fmt.Printf("%d. %s: %e\n", i+1, sortedProbs[i].Category, sortedProbs[i].Prob)
	}
}

func main() {
	data := []GenericTextData{
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
	displayTopCategories(probs, 3)

	// Test Case 1 - this doesnt work right
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

//func main() {
//	// Diverse data for different domains
//	data := []GenericTextData{
//		// Weather Domain
//		{"Sunny and warm conditions expected", "Weather"},
//		{"Heavy rainfall and thunderstorms in the area", "Weather"},
//		{"Cold temperatures and snow showers", "Weather"},
//
//		// E-commerce Domain
//		{"Latest smartphones with advanced features", "E-commerce"},
//		{"Fashionable clothing and accessories on sale", "E-commerce"},
//		{"Home appliances and electronics at discount prices", "E-commerce"},
//
//		// Healthcare Domain
//		{"Advancements in medical research and drug development", "Healthcare"},
//		{"Nutritional advice and diet planning", "Healthcare"},
//		{"Pediatric care and children's health services", "Healthcare"},
//
//		// Technology Domain
//		{"Innovations in artificial intelligence and machine learning", "Technology"},
//		{"Cloud computing services and solutions", "Technology"},
//		{"Cybersecurity threats and data protection measures", "Technology"},
//
//		// Education Domain
//		{"Online learning platforms and e-learning tools", "Education"},
//		{"Scholarship opportunities and academic programs", "Education"},
//		{"Educational policy reforms and teaching methods", "Education"},
//	}
//
//	classifier := NewNaiveBayesClassifier()
//	classifier.Train(data)
//
//	// Test Cases for different domains
//	testDescriptions := []string{
//		"Cloudy with a chance of rain",
//		"Summer fashion trends and styles",
//		"Breakthrough in cancer treatment",
//		"Developing secure web applications",
//		"Remote teaching and virtual classrooms",
//	}
//
//	for _, description := range testDescriptions {
//		predictedCategory, probs := classifier.PredictWithProbabilities(description)
//		fmt.Printf("Description: '%s'\nPredicted Category: %s\n", description, predictedCategory)
//		fmt.Println("Top Categories with Probabilities:")
//		displayTopCategories(probs, 3)
//		fmt.Println()
//	}
//}
