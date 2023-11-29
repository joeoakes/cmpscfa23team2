package main

import (
	"encoding/json"
	"fmt"
	"github.com/jdkato/prose/v2"
	"io"
	"math"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var stopWords = map[string]bool{
	"and": true, "or": true, "the": true, "in": true,
	"of": true, "a": true, "is": true, "to": true,
	"with": true, "for": true, "s": true, "you": true,
	"required": true, "then": true, "so": true, "our": true,
	"your": true, "their": true, "'s": true, "her": true,
	"him": true, "its": true, "he": true, "be": true,
	"we": true, "as": true,
	// More stop words can be added here
}

// GenericTextData represents any text data with an associated category.
type GenericTextData struct {
	Title       *string `json:"title"`
	URL         *string `json:"url"`
	Description *string `json:"description"`
	Salary      *string `json:"salary"` // Changed from Price to Salary
	Category    string  `json:"domain"`
	Company     *string `json:"company"`  // Added Company
	Location    *string `json:"location"` // Added Location
}

// NaiveBayesClassifier struct to hold model data.
type NaiveBayesClassifier struct {
	wordFrequencies  map[string]map[string]int
	categoryCounts   map[string]int
	totalWords       int
	totalUniqueWords int
}

// Define categoryProb struct
type categoryProb struct {
	Category string
	Prob     float64
}

// JSONData structure to match your data format
type JSONData struct {
	Domain string `json:"domain"`
	URL    string `json:"url"`
	Data   []struct {
		Title       string `json:"title"`
		URL         string `json:"url"`
		Description string `json:"description"`
		Company     string `json:"company"`
		Location    string `json:"location"`
		Salary      string `json:"salary"`
	} `json:"data"`
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

func isStopWord(word string) bool {
	// stopWords is a map of common words to be excluded from the analysis.
	return stopWords[word]
}

//// preprocessText preprocesses the text by converting it to lower case and excluding stop words.
//func preprocessText(text string) ([]string, error) {
//	doc, err := prose.NewDocument(text)
//	if err != nil {
//		return nil, err
//	}
//
//	var processedWords []string
//	for _, token := range doc.Tokens() {
//		word := strings.ToLower(token.Text)
//		if !isStopWord(word) && token.Tag != "PUNCT" {
//			processedWords = append(processedWords, word)
//		}
//	}
//	return processedWords, nil
//}

// preprocessText preprocesses the text by converting it to lower case and excluding stop words and special characters.
func preprocessText(text string) ([]string, error) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		return nil, err
	}

	wordRegexp := regexp.MustCompile(`\b[a-zA-Z]{2,}\b`) // Match words with 2 or more letters
	var processedWords []string

	for _, token := range doc.Tokens() {
		word := strings.ToLower(token.Text)
		if wordRegexp.MatchString(word) && !isStopWord(word) {
			processedWords = append(processedWords, word)
		}
	}
	return processedWords, nil
}

// Train takes generic text data and trains the classifier.
func (nbc *NaiveBayesClassifier) Train(data []GenericTextData) {
	uniqueWords := make(map[string]bool)

	for _, item := range data {
		if item.Description != nil && item.Title != nil {
			text := *item.Title + " " + *item.Description
			words, err := preprocessText(text)
			if err != nil {
				fmt.Println("Error preprocessing text:", err)
				continue
			}
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

// PredictWithProbabilities updated to return a sorted list of category probabilities
func (nbc *NaiveBayesClassifier) PredictWithProbabilities(skills []string) []categoryProb {
	if len(skills) == 0 {
		return nil
	}

	var combinedSkills string
	for _, skill := range skills {
		processedSkill, err := preprocessText(skill)
		if err != nil {
			fmt.Println("Error preprocessing skill:", err)
			continue
		}
		combinedSkills += strings.Join(processedSkill, " ") + " "
	}

	words, err := preprocessText(combinedSkills)
	if err != nil {
		fmt.Println("Error preprocessing skills:", err)
		return nil
	}

	categoryProbabilities := make(map[string]float64)
	for category := range nbc.categoryCounts {
		prob := nbc.calculateProbability(words, category)
		categoryProbabilities[category] += math.Exp(prob) // Sum probabilities for each category
	}

	// Convert to a slice and sort by probability
	var sortedCategories []categoryProb
	for category, prob := range categoryProbabilities {
		sortedCategories = append(sortedCategories, categoryProb{category, prob})
	}
	sort.Slice(sortedCategories, func(i, j int) bool {
		return sortedCategories[i].Prob > sortedCategories[j].Prob
	})

	return sortedCategories
}

// LoadDataFromJSON function updated to extract only title, description, and category
func LoadDataFromJSON(filename string, dataChan chan<- []GenericTextData, errChan chan<- error) {
	var jsonData JSONData
	var data []GenericTextData

	file, err := os.Open(filename)
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		errChan <- err
		return
	}

	err = json.Unmarshal(byteValue, &jsonData)
	if err != nil {
		errChan <- err
		return
	}

	for _, item := range jsonData.Data {
		data = append(data, GenericTextData{
			Title:       &item.Title,
			Description: &item.Description,
			Category:    jsonData.Domain,
		})
	}
	fmt.Println("Loading data from:", filename) // Debugging statement
	dataChan <- data                            // Send the data to the channel

}

// LoadDataFromMultipleJSONFiles loads data from multiple JSON files.
func LoadDataFromMultipleJSONFiles(filenames []string) ([]GenericTextData, error) {
	dataChan := make(chan []GenericTextData, len(filenames))
	errChan := make(chan error, len(filenames))

	for _, filename := range filenames {
		go LoadDataFromJSON(filename, dataChan, errChan) // Start a goroutine
	}

	var combinedData []GenericTextData
	for i := 0; i < len(filenames); i++ {
		select {
		case data := <-dataChan:
			combinedData = append(combinedData, data...)
		case err := <-errChan:
			return nil, err
		}
	}
	close(dataChan)
	close(errChan)
	fmt.Println("Data loading complete. Number of entries loaded:", len(combinedData)) // Debugging statement

	return combinedData, nil
}

// Function to extract and sort the most frequent words from the dataset
func getMostFrequentWords(data []GenericTextData, topN int) []string {
	wordFrequency := make(map[string]int)

	for _, item := range data {
		if item.Description != nil {
			words, err := preprocessText(*item.Description)
			if err != nil {
				fmt.Println("Error preprocessing text:", err)
				continue
			}
			for _, word := range words {
				if _, isStopWord := stopWords[word]; !isStopWord {
					wordFrequency[word]++
				}
			}
		}
	}
	// Convert to slice and sort
	type wordCountPair struct {
		Word  string
		Count int
	}
	var sortedWords []wordCountPair
	for word, count := range wordFrequency {
		sortedWords = append(sortedWords, wordCountPair{word, count})
	}
	sort.Slice(sortedWords, func(i, j int) bool {
		return sortedWords[i].Count > sortedWords[j].Count
	})

	// Extract top N words
	var topWords []string
	for i := 0; i < topN && i < len(sortedWords); i++ {
		topWords = append(topWords, sortedWords[i].Word)
	}
	return topWords
}

func main() {
	startTime := time.Now()

	jsonFiles := []string{
		"C:\\Users\\Public\\GoLandProjects\\PredictAi\\crab\\indeed_jobs.json",
		"C:\\Users\\Public\\GoLandProjects\\PredictAi\\crab\\indeed_jobs2.json",
		"C:\\Users\\Public\\GoLandProjects\\PredictAi\\crab\\indeed_jobs3.json",
	}

	// Load and combine data from all JSON files
	combinedData, err := LoadDataFromMultipleJSONFiles(jsonFiles)
	if err != nil {
		fmt.Println("Error loading data:", err)
		os.Exit(1)
	}

	// Debugging: Print the size of the combined data
	fmt.Println("Total data size before shuffle and split:", len(combinedData))

	// Shuffle and split data
	// Creating a new random source for reproducible sequences
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	// Shuffle and split data using the new random source
	rnd.Shuffle(len(combinedData), func(i, j int) { combinedData[i], combinedData[j] = combinedData[j], combinedData[i] })
	trainDataSize := int(float64(len(combinedData)) * 0.7)
	trainData := combinedData[:trainDataSize]

	// Debugging: Print the size of the training data
	fmt.Println("Training data size after split:", len(trainData))

	// Train classifier
	classifier := NewNaiveBayesClassifier()
	classifier.Train(trainData)

	//testSkillsSets := [][]string{
	//	{"cybersecurity", "encryption", "network security"}, // Cybersecurity skills
	//	{"patient care", "nursing", "medical diagnosis"},    // Healthcare skills
	//	{"financial analysis", "business development"},      // Business skills
	//	// Add more skill sets as needed
	//}

	// Use the getMostFrequentWords function to generate dynamic skill sets
	topSkills := getMostFrequentWords(trainData, 5) // Get top 5 frequent skills from training data

	// Test the classifier with dynamically generated skills
	for _, skill := range topSkills {
		if _, isStopWord := stopWords[skill]; !isStopWord {
			testSkills := []string{skill}
			sortedCategories := classifier.PredictWithProbabilities(testSkills)

			for _, categoryProb := range sortedCategories {
				fmt.Printf("Job titles relevant to '%s' skill in the '%s' category:\n", skill, categoryProb.Category)
				count := 0

				for _, job := range combinedData {
					if job.Category == categoryProb.Category && job.Title != nil {
						title := *job.Title
						if strings.Contains(strings.ToLower(*job.Description), skill) {
							fmt.Printf("Title: %s\n", title)
							count++
							if count >= 3 { // Limit to top 3 unique job titles
								break
							}
						}
					}
				}
				fmt.Println()
			}
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Execution time: %s\n", elapsedTime)
}
