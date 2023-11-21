package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"
)

// GenericTextData represents any text data with an associated category.
type GenericTextData struct {
	Title       *string `json:"title"`
	URL         *string `json:"url"`
	Description *string `json:"description"`
	Price       *string `json:"price"`
	Category    string  `json:"domain"`
	Source      *string `json:"source"`
	Timestamp   *string `json:"timestamp"`
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

// preprocessText preprocesses the text by converting it to lower case, removing punctuation, and excluding stop words.
func preprocessText(text string) []string {
	// Regular expression to match word characters
	wordRegexp := regexp.MustCompile("\\b\\w+\\b")

	words := wordRegexp.FindAllString(strings.ToLower(text), -1)
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
		// Check if Description is not nil before processing
		if item.Description != nil {
			words := preprocessText(*item.Description)
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
func (nbc *NaiveBayesClassifier) PredictWithProbabilities(description *string) (string, map[string]float64) {
	// Check if description is nil or empty before proceeding
	if description == nil || *description == "" {
		return "", nil
	}

	words := preprocessText(*description) // Dereference the pointer to get the string value
	if len(words) == 0 {
		return "", nil // Return empty category if no words after preprocessing
	}

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

// JSON structure to match your data format
type JSONData struct {
	Items []struct {
		Domain string `json:"domain"`
		Data   struct {
			Title       *string `json:"title"`
			URL         *string `json:"url"`
			Description *string `json:"description"`
			Price       *string `json:"price"`
			Metadata    *struct {
				Source    *string `json:"source"`
				Timestamp *string `json:"timestamp"`
			} `json:"metadata"`
		} `json:"data"`
	} `json:"items"`
}

// LoadDataFromJSON updated to handle your JSON structure
func LoadDataFromJSON(filename string) ([]GenericTextData, error) {
	var jsonData JSONData
	var data []GenericTextData

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &jsonData)
	if err != nil {
		return nil, err
	}

	for _, item := range jsonData.Items {
		var source, timestamp *string
		if item.Data.Metadata != nil {
			source = item.Data.Metadata.Source
			timestamp = item.Data.Metadata.Timestamp
		}
		data = append(data, GenericTextData{
			Title:       item.Data.Title,
			URL:         item.Data.URL,
			Description: item.Data.Description,
			Price:       item.Data.Price,
			Category:    item.Domain,
			Source:      source,
			Timestamp:   timestamp,
		})
	}

	return data, nil
}

// Helper function to convert string to *string
func strPtr(s string) *string {
	return &s
}

// Function to return the slice of GenericTextData
func getTrainingData() []GenericTextData {
	return []GenericTextData{
		// Agriculture
		{Description: strPtr("Agricultural technology and farm management"), Category: "Agriculture"},
		{Description: strPtr("Sustainable farming practices"), Category: "Agriculture"},

		// Aerospace
		{Description: strPtr("Aerospace engineering and avionics"), Category: "Aerospace"},
		{Description: strPtr("Spacecraft design and systems engineering"), Category: "Aerospace"},

		// Art
		{Description: strPtr("Photography and visual arts"), Category: "Art"},
		{Description: strPtr("Contemporary art curation and management"), Category: "Art"},

		// Automotive
		{Description: strPtr("Automotive engineering and design"), Category: "Automotive"},
		{Description: strPtr("Electric vehicle technology"), Category: "Automotive"},

		// Biotech
		{Description: strPtr("Biotechnology research and genetic engineering"), Category: "Biotech"},
		{Description: strPtr("Clinical trial management in biotech"), Category: "Biotech"},

		// Culinary
		{Description: strPtr("Culinary arts and restaurant management"), Category: "Culinary"},
		{Description: strPtr("Gourmet cuisine and fine dining"), Category: "Culinary"},

		// Design
		{Description: strPtr("Graphic design and digital media creation"), Category: "Design"},
		{Description: strPtr("Interior design and space planning"), Category: "Design"},

		// Education
		{Description: strPtr("Educational curriculum design and teaching"), Category: "Education"},
		{Description: strPtr("Educational technology and e-learning"), Category: "Education"},

		// Energy
		{Description: strPtr("Renewable energy and green technologies"), Category: "Energy"},
		{Description: strPtr("Oil and gas exploration and production"), Category: "Energy"},

		// Engineering
		{Description: strPtr("Civil engineering and infrastructure development"), Category: "Engineering"},
		{Description: strPtr("Mechanical engineering in manufacturing"), Category: "Engineering"},

		// Environmental
		{Description: strPtr("Environmental policy and sustainability"), Category: "Environmental"},
		{Description: strPtr("Conservation biology and ecosystem management"), Category: "Environmental"},

		// Event Management
		{Description: strPtr("Event planning and management"), Category: "Event Management"},
		{Description: strPtr("Corporate event coordination"), Category: "Event Management"},

		// Fashion
		{Description: strPtr("Fashion design and apparel merchandising"), Category: "Fashion"},
		{Description: strPtr("Trend forecasting in fashion"), Category: "Fashion"},

		// Finance
		{Description: strPtr("Managing financial assets"), Category: "Finance"},
		{Description: strPtr("Investment banking and asset management"), Category: "Finance"},

		// Fitness
		{Description: strPtr("Fitness training and sports coaching"), Category: "Fitness"},
		{Description: strPtr("Health and wellness coaching"), Category: "Fitness"},

		// Healthcare
		{Description: strPtr("Nursing care and patient management"), Category: "Healthcare"},
		{Description: strPtr("Clinical psychology and mental health counseling"), Category: "Healthcare"},

		// Hospitality
		{Description: strPtr("Hospitality services and hotel management"), Category: "Hospitality"},
		{Description: strPtr("Travel and tourism management"), Category: "Hospitality"},

		// HR
		{Description: strPtr("Human resources policies and employee relations"), Category: "HR"},
		{Description: strPtr("Talent acquisition and recruitment"), Category: "HR"},

		// Journalism
		{Description: strPtr("Journalism and news reporting"), Category: "Journalism"},
		{Description: strPtr("Investigative journalism and documentary production"), Category: "Journalism"},

		// Legal
		{Description: strPtr("Legal compliance and corporate law"), Category: "Legal"},
		{Description: strPtr("Intellectual property law and patent strategy"), Category: "Legal"},

		// Marketing
		{Description: strPtr("Marketing strategy and brand development"), Category: "Marketing"},
		{Description: strPtr("Digital marketing and social media strategy"), Category: "Marketing"},

		// Media
		{Description: strPtr("Film production and media arts"), Category: "Media"},
		{Description: strPtr("Broadcast journalism and television production"), Category: "Media"},

		// Operations
		{Description: strPtr("Supply chain logistics and management"), Category: "Operations"},
		{Description: strPtr("Business operations and process improvement"), Category: "Operations"},

		// Pharmaceutical
		{Description: strPtr("Pharmaceutical research and development"), Category: "Pharmaceutical"},
		{Description: strPtr("Regulatory affairs in pharmaceuticals"), Category: "Pharmaceutical"},

		// Politics
		{Description: strPtr("International diplomacy and policy analysis"), Category: "Politics"},
		{Description: strPtr("Political campaign management and strategy"), Category: "Politics"},

		// Public Relations
		{Description: strPtr("Public relations and corporate communications"), Category: "Public Relations"},
		{Description: strPtr("Crisis communication and reputation management"), Category: "Public Relations"},

		// Real Estate
		{Description: strPtr("Real estate development and property management"), Category: "Real Estate"},
		{Description: strPtr("Commercial real estate investment and brokerage"), Category: "Real Estate"},

		// Retail
		{Description: strPtr("Retail sales management and merchandising"), Category: "Retail"},
		{Description: strPtr("E-commerce and digital retail strategy"), Category: "Retail"},

		// Technology
		{Description: strPtr("Software development in Go"), Category: "Technology"},
		{Description: strPtr("Developing advanced machine learning models"), Category: "Technology"},

		// Veterinary
		{Description: strPtr("Veterinary services and animal care"), Category: "Veterinary"},
		{Description: strPtr("Animal nutrition and wellness"), Category: "Veterinary"},
	}
}

// prints non-empty fields from the GenericTextData
func getScrapedData() {
	jsonFile := "C:\\Users\\Public\\GoLandProjects\\PredictAi\\cuda\\ML\\scrapedData.json"
	//fmt.Println(jsonFile)

	jsonData, err := LoadDataFromJSON(jsonFile)
	if err != nil {
		fmt.Println("Error loading data from JSON:", err)
		os.Exit(1)
	}

	for _, item := range jsonData {
		fmt.Println("Category:", item.Category)
		if item.Title != nil && *item.Title != "" {
			fmt.Println("Title:", *item.Title)
		}
		if item.URL != nil && *item.URL != "" {
			fmt.Println("URL:", *item.URL)
		}
		if item.Description != nil && *item.Description != "" {
			fmt.Println("Description:", *item.Description)
		}
		if item.Price != nil && *item.Price != "" {
			fmt.Println("Price:", *item.Price)
		}
		if item.Source != nil && *item.Source != "" {
			fmt.Println("Source:", *item.Source)
		}
		if item.Timestamp != nil && *item.Timestamp != "" {
			fmt.Println("Timestamp:", *item.Timestamp)
		}
		fmt.Println()
	}
}
func main() {
	getScrapedData()
	fmt.Println()

	data := getTrainingData()

	classifier := NewNaiveBayesClassifier()
	classifier.Train(data)

	newJobListing := "Programming in Python and Go"
	predictedIndustry, probs := classifier.PredictWithProbabilities(&newJobListing)

	fmt.Printf("Predicted Industry: %s\n", predictedIndustry)
	fmt.Println("Top Industries with Probabilities:")
	displayTopCategories(probs, 3)

	// Test Case 1 - this doesn't work right
	testJobDescription1 := "Designing user interfaces and user experiences"
	predictedIndustry1, _ := classifier.PredictWithProbabilities(&testJobDescription1)
	fmt.Printf("Test 1 - Job Description: '%s'\nPredicted Industry: %s\n\n", testJobDescription1, predictedIndustry1)

	// Test Case 2
	testJobDescription2 := "Developing policies for environmental conservation"
	predictedIndustry2, _ := classifier.PredictWithProbabilities(&testJobDescription2)
	fmt.Printf("Test 2 - Job Description: '%s'\nPredicted Industry: %s\n\n", testJobDescription2, predictedIndustry2)

	// Test Case 3
	testJobDescription3 := "Managing investment portfolios and financial risks"
	predictedIndustry3, _ := classifier.PredictWithProbabilities(&testJobDescription3)
	fmt.Printf("Test 3 - Job Description: '%s'\nPredicted Industry: %s\n\n", testJobDescription3, predictedIndustry3)

	// Test Case 4
	testJobDescription4 := "Culinary techniques and kitchen management"
	predictedIndustry4, _ := classifier.PredictWithProbabilities(&testJobDescription4)
	fmt.Printf("Test 4 - Job Description: '%s'\nPredicted Industry: %s\n", testJobDescription4, predictedIndustry4)
}
