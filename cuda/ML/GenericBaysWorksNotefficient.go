package main

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/jdkato/prose/v2"
//	"io"
//	"math"
//	"math/rand"
//	"os"
//	"regexp"
//	"sort"
//	"strings"
//	"time"
//)
//
//var stopWords = map[string]bool{
//	"and": true, "or": true, "the": true, "in": true,
//	"of": true, "a": true, "is": true, "to": true,
//	"with": true, "for": true, "s": true, "you": true,
//	"required": true, "then": true, "so": true, "our": true,
//	"your": true,
//	// More stop words can be added here
//}
//
//// GenericTextData represents any text data with an associated category.
//type GenericTextData struct {
//	Title       *string `json:"title"`
//	URL         *string `json:"url"`
//	Description *string `json:"description"`
//	Salary      *string `json:"salary"` // Changed from Price to Salary
//	Category    string  `json:"domain"`
//	Company     *string `json:"company"`  // Added Company
//	Location    *string `json:"location"` // Added Location
//}
//
//// NaiveBayesClassifier struct to hold model data.
//type NaiveBayesClassifier struct {
//	wordFrequencies  map[string]map[string]int
//	categoryCounts   map[string]int
//	totalWords       int
//	totalUniqueWords int
//}
//
//// Define categoryProb struct
//type categoryProb struct {
//	Category string
//	Prob     float64
//}
//
//// JSONData structure to match your data format
//type JSONData struct {
//	Domain   string            `json:"domain"`
//	URL      string            `json:"url"`
//	Data     []GenericTextData `json:"data"`
//	Metadata struct {
//		Source    string `json:"source"`
//		Timestamp string `json:"timestamp"`
//	} `json:"metadata"`
//}
//
//// NewNaiveBayesClassifier creates a new Naive Bayes Classifier.
//func NewNaiveBayesClassifier() *NaiveBayesClassifier {
//	return &NaiveBayesClassifier{
//		wordFrequencies:  make(map[string]map[string]int),
//		categoryCounts:   make(map[string]int),
//		totalWords:       0,
//		totalUniqueWords: 0,
//	}
//}
//
//func isStopWord(word string) bool {
//	// stopWords is a map of common words to be excluded from the analysis.
//	fmt.Println(stopWords)
//	return stopWords[word]
//}
//
////// preprocessText preprocesses the text by converting it to lower case and excluding stop words.
////func preprocessText(text string) ([]string, error) {
////	doc, err := prose.NewDocument(text)
////	if err != nil {
////		return nil, err
////	}
////
////	var processedWords []string
////	for _, token := range doc.Tokens() {
////		word := strings.ToLower(token.Text)
////		if !isStopWord(word) && token.Tag != "PUNCT" {
////			processedWords = append(processedWords, word)
////		}
////	}
////	return processedWords, nil
////}
//
//// preprocessText preprocesses the text by converting it to lower case and excluding stop words and special characters.
//func preprocessText(text string) ([]string, error) {
//	doc, err := prose.NewDocument(text)
//	if err != nil {
//		return nil, err
//	}
//
//	wordRegexp := regexp.MustCompile(`\b\w+\b`) // Regular expression to match word characters
//	var processedWords []string
//
//	for _, token := range doc.Tokens() {
//		word := strings.ToLower(token.Text)
//		if wordRegexp.MatchString(word) && !isStopWord(word) {
//			processedWords = append(processedWords, word)
//		}
//	}
//	return processedWords, nil
//}
//
//// Train takes generic text data and trains the classifier.
//func (nbc *NaiveBayesClassifier) Train(data []GenericTextData) {
//	uniqueWords := make(map[string]bool)
//
//	for _, item := range data {
//		text := *item.Title + " " + *item.Description
//		words, err := preprocessText(text)
//		if err != nil {
//			fmt.Println("Error preprocessing text:", err)
//			continue
//		}
//		category := item.Category
//
//		if nbc.wordFrequencies[category] == nil {
//			nbc.wordFrequencies[category] = make(map[string]int)
//		}
//
//		for _, word := range words {
//			if !uniqueWords[word] {
//				uniqueWords[word] = true
//				nbc.totalUniqueWords++
//			}
//
//			nbc.wordFrequencies[category][word]++
//			nbc.totalWords++
//		}
//
//		nbc.categoryCounts[category]++
//	}
//}
//
//// calculateProbability calculates the probability of a set of words belonging to a given category.
//func (nbc *NaiveBayesClassifier) calculateProbability(words []string, category string) float64 {
//	vocabSize := nbc.totalUniqueWords
//	prob := math.Log(float64(nbc.categoryCounts[category]) / float64(len(nbc.categoryCounts)))
//
//	for _, word := range words {
//		wordFrequency := nbc.wordFrequencies[category][word]
//		prob += math.Log(float64(wordFrequency+1) / float64(nbc.totalWords+vocabSize))
//	}
//
//	return prob
//}
//
//// PredictWithProbabilities updated to return a sorted list of category probabilities
//func (nbc *NaiveBayesClassifier) PredictWithProbabilities(skills []string) []categoryProb {
//	if len(skills) == 0 {
//		return nil
//	}
//
//	var combinedSkills string
//	for _, skill := range skills {
//		processedSkill, err := preprocessText(skill)
//		if err != nil {
//			fmt.Println("Error preprocessing skill:", err)
//			continue
//		}
//		combinedSkills += strings.Join(processedSkill, " ") + " "
//	}
//
//	words, err := preprocessText(combinedSkills)
//	if err != nil {
//		fmt.Println("Error preprocessing skills:", err)
//		return nil
//	}
//
//	categoryProbabilities := make(map[string]float64)
//	for category := range nbc.categoryCounts {
//		prob := nbc.calculateProbability(words, category)
//		categoryProbabilities[category] += math.Exp(prob) // Sum probabilities for each category
//	}
//
//	// Convert to a slice and sort by probability
//	var sortedCategories []categoryProb
//	for category, prob := range categoryProbabilities {
//		sortedCategories = append(sortedCategories, categoryProb{category, prob})
//	}
//	sort.Slice(sortedCategories, func(i, j int) bool {
//		return sortedCategories[i].Prob > sortedCategories[j].Prob
//	})
//
//	return sortedCategories
//}
//
//// LoadDataFromJSON function updated to match the JSON structure
//func LoadDataFromJSON(filename string) ([]GenericTextData, error) {
//	file, err := os.Open(filename)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	var jsonData []JSONData
//	byteValue, err := io.ReadAll(file)
//	if err != nil {
//		return nil, err
//	}
//
//	err = json.Unmarshal(byteValue, &jsonData)
//	if err != nil {
//		return nil, err
//	}
//
//	var data []GenericTextData
//	for _, jData := range jsonData {
//		for _, item := range jData.Data {
//			data = append(data, item)
//		}
//	}
//
//	return data, nil
//}
//
//// LoadDataFromMultipleJSONFiles loads data from multiple JSON files.
//func LoadDataFromMultipleJSONFiles(filenames []string) ([]GenericTextData, error) {
//	var combinedData []GenericTextData
//
//	for _, filename := range filenames {
//		fileData, err := LoadDataFromJSON(filename)
//		if err != nil {
//			return nil, fmt.Errorf("error loading file %s: %v", filename, err)
//		}
//		combinedData = append(combinedData, fileData...)
//	}
//
//	return combinedData, nil
//}
//
//// Function to extract and sort the most frequent words from the dataset
//func getMostFrequentWords(data []GenericTextData, topN int) []string {
//	wordFrequency := make(map[string]int)
//
//	for _, item := range data {
//		if item.Description != nil {
//			words, err := preprocessText(*item.Description)
//			if err != nil {
//				fmt.Println("Error preprocessing text:", err)
//				continue
//			}
//			for _, word := range words {
//				if _, isStopWord := stopWords[word]; !isStopWord {
//					wordFrequency[word]++
//				}
//			}
//		}
//	}
//	// Convert to slice and sort
//	type wordCountPair struct {
//		Word  string
//		Count int
//	}
//	var sortedWords []wordCountPair
//	for word, count := range wordFrequency {
//		sortedWords = append(sortedWords, wordCountPair{word, count})
//	}
//	sort.Slice(sortedWords, func(i, j int) bool {
//		return sortedWords[i].Count > sortedWords[j].Count
//	})
//
//	// Extract top N words
//	var topWords []string
//	for i := 0; i < topN && i < len(sortedWords); i++ {
//		topWords = append(topWords, sortedWords[i].Word)
//	}
//	return topWords
//}
//
//func main() {
//	startTime := time.Now()
//
//	jsonFiles := []string{
//		"C:\\Users\\mathe\\GoLandProjects\\cmpscfa23team2\\crab\\output\\combined_jobs.json",
//	}
//
//	// Load and combine data from all JSON files
//	combinedData, err := LoadDataFromMultipleJSONFiles(jsonFiles)
//	if err != nil {
//		fmt.Println("Error loading data:", err)
//		os.Exit(1)
//	}
//
//	// Shuffle and split data
//	// Creating a new random source for reproducible sequences
//	src := rand.NewSource(time.Now().UnixNano())
//	rnd := rand.New(src)
//
//	// Shuffle and split data using the new random source
//	rnd.Shuffle(len(combinedData), func(i, j int) { combinedData[i], combinedData[j] = combinedData[j], combinedData[i] })
//	trainDataSize := int(float64(len(combinedData)) * 0.8)
//	trainData := combinedData[:trainDataSize]
//
//	// Train classifier
//	classifier := NewNaiveBayesClassifier()
//	classifier.Train(trainData)
//
//	//testSkillsSets := [][]string{
//	//	{"cybersecurity", "encryption", "network security"}, // Cybersecurity skills
//	//	{"patient care", "nursing", "medical diagnosis"},    // Healthcare skills
//	//	{"financial analysis", "business development"},      // Business skills
//	//	// Add more skill sets as needed
//	//}
//
//	// Use the getMostFrequentWords function to generate dynamic skill sets
//	topSkills := getMostFrequentWords(trainData, 5) // Get top 5 frequent skills from training data
//
//	// Test the classifier with dynamically generated skills
//	for _, skill := range topSkills {
//		testSkills := []string{skill}
//		sortedCategories := classifier.PredictWithProbabilities(testSkills)
//
//		for _, categoryProb := range sortedCategories {
//			fmt.Printf("Job titles relevant to '%s' skill in the '%s' category:\n", skill, categoryProb.Category)
//			count := 0
//			shownTitles := make(map[string]bool) // Map to track shown job titles
//
//			for _, job := range combinedData {
//				if job.Category == categoryProb.Category && job.Title != nil {
//					title := *job.Title
//					if !shownTitles[title] {
//						fmt.Printf("Title: %s\n", title)
//						shownTitles[title] = true // Mark this title as shown
//						count++
//						if count >= 3 { // Limit to top 3 unique job titles
//							break
//						}
//					}
//				}
//			}
//			fmt.Println()
//		}
//	}
//	elapsedTime := time.Since(startTime)
//	fmt.Printf("Execution time: %s\n", elapsedTime)
//}
