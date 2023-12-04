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
	"we": true, "as": true, "on": true,
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
	wordFrequencies  map[string]map[string]float64
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
	Metadata struct {
		Source    string `json:"source"`
		Timestamp string `json:"timestamp"`
	} `json:"metadata"`
}

// NewNaiveBayesClassifier creates a new Naive Bayes Classifier.
func NewNaiveBayesClassifier() *NaiveBayesClassifier {
	return &NaiveBayesClassifier{
		wordFrequencies:  make(map[string]map[string]float64),
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

	wordRegexp := regexp.MustCompile(`\b[a-zA-Z+#]{2,}\b`) // Enhanced regex
	var processedWords []string

	for _, token := range doc.Tokens() {
		word := strings.ToLower(token.Text)
		if wordRegexp.MatchString(word) && !isStopWord(word) {
			processedWords = append(processedWords, word)
		}
	}
	return processedWords, nil
}

// Train with n-gram model for better context understanding
func (nbc *NaiveBayesClassifier) Train(data []GenericTextData) {
	uniqueWords := make(map[string]bool)
	nGramSize := 2 // Example: bi-grams

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
				nbc.wordFrequencies[category] = make(map[string]float64)
			}

			nGrams := createNGrams(words, nGramSize) // Function to create n-grams from words

			for _, nGram := range nGrams {
				if !uniqueWords[nGram] {
					uniqueWords[nGram] = true
					nbc.totalUniqueWords++
				}
				nbc.wordFrequencies[category][nGram]++
				nbc.totalWords++
			}

			nbc.categoryCounts[category]++
		}
	}
}

// CalculateTermFrequency calculates the frequency of each word in a document
func CalculateTermFrequency(text string) map[string]float64 {
	wordFreq := make(map[string]int)
	totalWords := 0
	words := strings.Fields(text)

	for _, word := range words {
		wordFreq[word]++
		totalWords++
	}

	tf := make(map[string]float64)
	for word, count := range wordFreq {
		tf[word] = float64(count) / float64(totalWords)
	}

	return tf
}

// CalculateIDF calculates the inverse document frequency for each word across all documents
func CalculateIDF(corpus []string) map[string]float64 {
	docCount := make(map[string]int)
	numDocs := len(corpus)

	for _, text := range corpus {
		wordSet := make(map[string]bool)
		words := strings.Fields(text)

		for _, word := range words {
			wordSet[word] = true
		}

		for word := range wordSet {
			docCount[word]++
		}
	}

	idf := make(map[string]float64)
	for word, count := range docCount {
		idf[word] = math.Log(float64(numDocs) / float64(count))
	}

	return idf
}

// TrainWithTFIDF updates the Naive Bayes Classifier with TF-IDF features
//
//	func (nbc *NaiveBayesClassifier) TrainWithTFIDF(data []GenericTextData) {
//		corpus := make([]string, len(data))
//		for i, item := range data {
//			corpus[i] = *item.Description
//		}
//
//		idf := CalculateIDF(corpus)
//
//		for _, item := range data {
//			category := item.Category
//			tf := CalculateTermFrequency(*item.Description)
//
//			for word, freq := range tf {
//				tfidf := freq * idf[word]
//
//				if nbc.wordFrequencies[category] == nil {
//					nbc.wordFrequencies[category] = make(map[string]float64)
//				}
//
//				nbc.wordFrequencies[category][word] += tfidf
//			}
//
//			nbc.categoryCounts[category]++
//		}
//	}
func TestModel(classifier *NaiveBayesClassifier, testData []GenericTextData) {
	correctPredictions := 0

	for _, data := range testData {
		if data.Description != nil {
			predictedCategories := classifier.PredictWithTFIDF(*data.Description)

			// Assuming the most probable category is the first in the sorted list
			if len(predictedCategories) > 0 && predictedCategories[0].Category == data.Category {
				correctPredictions++
			}
		}
	}

	accuracy := float64(correctPredictions) / float64(len(testData))
	fmt.Printf("Accuracy: %.2f%%\n", accuracy*100)
}

// PredictWithTFIDF predicts the category based on TF-IDF features
func (nbc *NaiveBayesClassifier) PredictWithTFIDF(text string) []categoryProb {
	tf := CalculateTermFrequency(text)
	idf := CalculateIDF([]string{text})

	categoryProbabilities := make(map[string]float64)
	for category := range nbc.categoryCounts {
		prob := math.Log(float64(nbc.categoryCounts[category]) / float64(len(nbc.categoryCounts)))

		for word, freq := range tf {
			tfidf := freq * idf[word]
			prob += math.Log((nbc.wordFrequencies[category][word] + 1.0) * tfidf)
		}

		categoryProbabilities[category] = math.Exp(prob)
	}

	var sortedCategories []categoryProb
	for category, prob := range categoryProbabilities {
		sortedCategories = append(sortedCategories, categoryProb{category, prob})
	}

	sort.Slice(sortedCategories, func(i, j int) bool {
		return sortedCategories[i].Prob > sortedCategories[j].Prob
	})

	return sortedCategories
}

// Helper function to create n-grams
func createNGrams(words []string, size int) []string {
	var nGrams []string
	for i := 0; i < len(words)-(size-1); i++ {
		nGram := ""
		for j := 0; j < size; j++ {
			if j > 0 {
				nGram += " "
			}
			nGram += words[i+j]
		}
		nGrams = append(nGrams, nGram)
	}
	return nGrams
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

// PredictWithProbabilities improved to handle unseen words and weighted probabilities
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
		prob := math.Log(float64(nbc.categoryCounts[category]) / float64(len(nbc.categoryCounts)))

		for _, word := range words {
			wordFrequency := nbc.wordFrequencies[category][word]
			weight := 1.0 // Assign different weights based on word type or importance
			smoothedProb := math.Log((float64(wordFrequency) + weight) / float64(nbc.totalWords+nbc.totalUniqueWords))
			prob += smoothedProb
		}

		categoryProbabilities[category] = math.Exp(prob) // Sum probabilities for each category
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
	var jsonData []JSONData
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

	for _, jsonItem := range jsonData {
		for _, item := range jsonItem.Data {
			data = append(data, GenericTextData{
				Title:       &item.Title,
				Description: &item.Description,
				Category:    jsonItem.Domain,
				Company:     &item.Company,
				Location:    &item.Location,
				Salary:      &item.Salary,
				URL:         &item.URL,
			})
		}
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
func getMostFrequentWords(data []GenericTextData, topN int) []string {
	wordFrequency := make(map[string]int)
	wordSet := make(map[string]bool) // Set to keep track of unique words

	for _, item := range data {
		if item.Description != nil {
			words, err := preprocessText(*item.Description)
			if err != nil {
				fmt.Println("Error preprocessing text:", err)
				continue
			}
			for _, word := range words {
				if _, isStopWord := stopWords[word]; !isStopWord && !wordSet[word] {
					wordFrequency[word]++
					wordSet[word] = true // Mark word as seen
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

// TrainWithTFIDF updates the Naive Bayes Classifier with TF-IDF features
func (nbc *NaiveBayesClassifier) TrainWithTFIDF(data []GenericTextData) {
	corpus := make([]string, len(data))
	for i, item := range data {
		corpus[i] = *item.Description
	}

	idf := CalculateIDF(corpus)

	for _, item := range data {
		category := item.Category
		tf := CalculateTermFrequency(*item.Description)

		for word, freq := range tf {
			tfidf := freq * idf[word]

			if nbc.wordFrequencies[category] == nil {
				nbc.wordFrequencies[category] = make(map[string]float64)
			}

			nbc.wordFrequencies[category][word] += tfidf
		}

		nbc.categoryCounts[category]++
	}
}

// // Function to extract and sort the most frequent words from the dataset
//
//	func getMostFrequentWords(data []GenericTextData, topN int) []string {
//		wordFrequency := make(map[string]int)
//
//		for _, item := range data {
//			if item.Description != nil {
//				words, err := preprocessText(*item.Description)
//				if err != nil {
//					fmt.Println("Error preprocessing text:", err)
//					continue
//				}
//				for _, word := range words {
//					if _, isStopWord := stopWords[word]; !isStopWord {
//						wordFrequency[word]++
//					}
//				}
//			}
//		}
//		// Convert to slice and sort
//		type wordCountPair struct {
//			Word  string
//			Count int
//		}
//		var sortedWords []wordCountPair
//		for word, count := range wordFrequency {
//			fmt.Printf("Top skill: %s, Frequency: %d\n", word, count) // Debug top skills
//			sortedWords = append(sortedWords, wordCountPair{word, count})
//		}
//		sort.Slice(sortedWords, func(i, j int) bool {
//			return sortedWords[i].Count > sortedWords[j].Count
//		})
//
//		// Extract top N words
//		var topWords []string
//		for i := 0; i < topN && i < len(sortedWords); i++ {
//			topWords = append(topWords, sortedWords[i].Word)
//		}
//		return topWords
//	}
func main2() {
	startTime := time.Now()

	jsonFiles := []string{
		"C:\\Users\\mathe\\GolandProjects\\cmpscfa23team2\\crab\\output\\combined_jobs.json",
	}

	// Load data
	combinedData, err := LoadDataFromMultipleJSONFiles(jsonFiles)
	if err != nil {
		fmt.Println("Error loading data:", err)
		os.Exit(1)
	}

	// Split data into training and testing
	trainDataSize := int(float64(len(combinedData)) * 0.7)
	trainData := combinedData[:trainDataSize]
	testData := combinedData[trainDataSize:]

	// Train classifier
	classifier := NewNaiveBayesClassifier()
	classifier.TrainWithTFIDF(trainData)

	// Test classifier
	TestModel(classifier, testData)

	elapsedTime := time.Since(startTime)
	fmt.Printf("Execution time: %s\n", elapsedTime)
}
func extractSkillsAndQualifications(description string) []string {
	// Define a regular expression to identify skills and qualifications
	skillRegexp := regexp.MustCompile(`[a-zA-Z+#]+`) // Modify this regex according to your needs

	matches := skillRegexp.FindAllString(description, -1)

	var skills []string
	for _, match := range matches {
		if !isStopWord(match) {
			skills = append(skills, match)
		}
	}
	return skills
}

// ExtractSkillsAndQualifications function to extract specific skills
func ExtractSkillsAndQualifications(description string) []string {
	skillRegexp := regexp.MustCompile(`[a-zA-Z+#]+`)
	matches := skillRegexp.FindAllString(description, -1)

	var skills []string
	for _, match := range matches {
		if !isStopWord(match) {
			skills = append(skills, match)
		}
	}
	return skills
}

// Function to calculate the demand for skills in each category
func CalculateSkillDemand(data []GenericTextData) map[string]map[string]int {
	skillDemand := make(map[string]map[string]int)

	for _, job := range data {
		if job.Description != nil {
			skills := ExtractSkillsAndQualifications(*job.Description)
			if skillDemand[job.Category] == nil {
				skillDemand[job.Category] = make(map[string]int)
			}
			for _, skill := range skills {
				skillDemand[job.Category][skill]++
			}
		}
	}
	return skillDemand
}

func main() {
	startTime := time.Now()

	jsonFiles := []string{"C:\\Users\\mathe\\GolandProjects\\cmpscfa23team2\\crab\\output\\combined_jobs.json"} // Update the path to your JSON file

	// Load and combine data from all JSON files
	combinedData, err := LoadDataFromMultipleJSONFiles(jsonFiles)
	if err != nil {
		fmt.Println("Error loading data:", err)
		os.Exit(1)
	}

	// Calculate the demand for skills in each category
	skillDemand := CalculateSkillDemand(combinedData)

	// Print the demand for top 5 skills in each category
	for category, skills := range skillDemand {
		fmt.Printf("Category: %s\n", category)
		// Sort and pick top 5 skills
		topSkills := getTopNSkills(skills, 5)
		for skill, count := range topSkills {
			fmt.Printf("Skill: %s, Demand: %d\n", skill, count)
		}
		fmt.Println()
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Execution time: %s\n", elapsedTime)
}

// getTopNSkills returns the top N skills based on demand
func getTopNSkills(skillMap map[string]int, n int) map[string]int {
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range skillMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	topSkills := make(map[string]int)
	for i, kv := range ss {
		if i < n {
			topSkills[kv.Key] = kv.Value
		}
	}

	return topSkills
}

func main3() {
	startTime := time.Now()

	jsonFiles := []string{
		"C:\\Users\\mathe\\GolandProjects\\cmpscfa23team2\\crab\\output\\combined_jobs.json",
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
	// Assuming trainDataSize is the size of your training data
	trainData := combinedData[trainDataSize:]

	// Debugging: Print the size of the training data
	fmt.Println("Training data size after split:", len(trainData))

	// Train classifier
	classifier := NewNaiveBayesClassifier()
	classifier.TrainWithTFIDF(trainData)

	//testSkillsSets := [][]string{
	//	{"cybersecurity", "encryption", "network security"}, // Cybersecurity skills
	//	{"patient care", "nursing", "medical diagnosis"},    // Healthcare skills
	//	{"financial analysis", "business development"},      // Business skills
	//	// Add more skill sets as needed
	//}

	// Use the getMostFrequentWords function to generate dynamic skill sets
	topSkills := getMostFrequentWords(trainData, 3) // Get top 5 frequent skills from training data
	for _, job := range combinedData {
		if job.Description != nil {
			skills := extractSkillsAndQualifications(*job.Description)
			fmt.Printf("Job Title: %s, Extracted Skills: %v\n", *job.Title, skills)
		}
	}
	// Test the classifier with dynamically generated skills
	for _, skill := range topSkills {
		if _, isStopWord := stopWords[skill]; !isStopWord {
			testSkills := []string{skill}
			sortedCategories := classifier.PredictWithProbabilities(testSkills)

			for _, categoryProb := range sortedCategories {
				fmt.Printf("Job titles relevant to '%s' skill in the '%s' category:\n", skill, categoryProb.Category)
				count := 0
				uniqueTitles := make(map[string]bool) // Map to keep track of unique titles

				for _, job := range combinedData {
					if job.Category == categoryProb.Category && job.Title != nil {
						title := *job.Title
						if strings.Contains(strings.ToLower(*job.Description), skill) && !uniqueTitles[title] {
							fmt.Printf("Title: %s\n", title)
							uniqueTitles[title] = true
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
