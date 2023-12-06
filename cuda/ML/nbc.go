package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

// JobData represents a single job entry.
type JobData struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Salary      string `json:"salary"`
	Company     string `json:"company"`
	Location    string `json:"location"`
}

// JobDataContainer represents the structure of your JSON file.
type JobDataContainer struct {
	Domain   string    `json:"domain"`
	URL      string    `json:"url"`
	Data     []JobData `json:"data"`
	Metadata struct {
		Source    string `json:"source"`
		Timestamp string `json:"timestamp"`
	} `json:"metadata"`
}

// SkillData represents the demand for a skill in a category.
type SkillData struct {
	Skill   string
	Demand  int
	Matches []JobData
}

// NaiveBayesClassifier struct to hold model data.
type NaiveBayesClassifier struct {
	wordFrequencies  map[string]map[string]int
	categoryCounts   map[string]int
	totalWords       int
	totalUniqueWords int
	skillSets        map[string][]string
}

// NewNaiveBayesClassifier creates a new Naive Bayes Classifier.
func NewNaiveBayesClassifier() *NaiveBayesClassifier {
	return &NaiveBayesClassifier{
		wordFrequencies:  make(map[string]map[string]int),
		categoryCounts:   make(map[string]int),
		totalWords:       0,
		totalUniqueWords: 0,
		skillSets: map[string][]string{
			"Tech":     {"design", "Python", "Java", "React"},
			"Business": {"Management", "Finance", "marketing", "management", "Microsoft Office"},
			"Law":      {"school of law", "Litigation"},
		},
	}
}

// Train takes generic text data and trains the classifier.
func (nbc *NaiveBayesClassifier) Train(data []JobData, domain string) {
	uniqueWords := make(map[string]bool)

	for _, item := range data {
		text := item.Title + " " + item.Description
		words := strings.Fields(text)

		if nbc.wordFrequencies[domain] == nil {
			nbc.wordFrequencies[domain] = make(map[string]int)
		}

		for _, word := range words {
			word = strings.ToLower(word)
			if !uniqueWords[word] {
				uniqueWords[word] = true
				nbc.totalUniqueWords++
			}

			nbc.wordFrequencies[domain][word]++
			nbc.totalWords++
		}

		nbc.categoryCounts[domain]++
	}
}

// PredictWithProbabilities predicts the most likely category for a set of skills.
func (nbc *NaiveBayesClassifier) PredictWithProbabilities(skills []string) []string {
	// Assuming implementation of probability calculation
	// Returns a list of categories sorted by their probability score
	return []string{"Category1", "Category2"} // Placeholder
}

// LoadDataFromJSON updated to match the new GenericTextData structure.
func LoadDataFromJSON(filename string) ([]JobData, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var container JobDataContainer
	err = json.Unmarshal(file, &container)
	if err != nil {
		return nil, err
	}

	return container.Data, nil
}

type jobMatch struct {
	job           JobData
	count         int
	matchedSkills []string
}

// PredictBestMatchingJob predicts the best matching job based on skills.
func (nbc *NaiveBayesClassifier) PredictBestMatchingJob(domain string, data []JobData) {
	skills := nbc.skillSets[domain]
	skillDemand := make(map[string]SkillData)

	// Initialize SkillData
	for _, skill := range skills {
		skillDemand[skill] = SkillData{Skill: skill, Demand: 0, Matches: make([]JobData, 0)}
	}

	// Count demand and find matching jobs
	jobSkills := make(map[string][]string) // Map to store matched skills for each job
	for _, job := range data {
		for _, skill := range skills {
			if strings.Contains(job.Description, skill) || strings.Contains(job.Title, skill) {
				jobSkills[job.URL] = append(jobSkills[job.URL], skill)

				// Update demand
				skillData := skillDemand[skill]
				skillData.Demand++
				skillDemand[skill] = skillData
			}
		}
	}

	// Find top 3 jobs with most matching skills
	var topJobs []jobMatch
	for _, job := range data {
		if matchedSkills, exists := jobSkills[job.URL]; exists {
			topJobs = append(topJobs, jobMatch{job, len(matchedSkills), matchedSkills})
		}
	}
	sort.Slice(topJobs, func(i, j int) bool {
		return topJobs[i].count > topJobs[j].count
	})
	if len(topJobs) > 3 {
		topJobs = topJobs[:3]
	}

	// Display skill demand and top jobs
	for _, skill := range skills {
		fmt.Printf("Skill: %s, Demand: %d\n", skill, skillDemand[skill].Demand)
	}
	fmt.Printf("\nTop Jobs for '%s' Domain:\n", domain)
	for _, match := range topJobs {
		fmt.Printf("Job Title: %s, URL: %s, Company: %s, Location: %s, Salary: %s\nDescription: %s\nMatching Skills: %d [%s]\n\n",
			match.job.Title, match.job.URL, match.job.Company, match.job.Location, match.job.Salary, match.job.Description, match.count, strings.Join(match.matchedSkills, ", "))
	}
}

func main() {
	classifier := NewNaiveBayesClassifier()

	// File paths
	filePaths := []string{
		"C:\\Users\\mathe\\GoLandProjects\\cmpscfa23team2\\crab\\output\\Law_jobs.json",
		"C:\\Users\\mathe\\GoLandProjects\\cmpscfa23team2\\crab\\output\\Business_jobs.json",
		"C:\\Users\\mathe\\GoLandProjects\\cmpscfa23team2\\crab\\output\\Software Engineer_jobs.json",
	}

	// Mapping of file paths to their respective domains
	domainMapping := map[string]string{
		"Law_jobs.json":               "Law",
		"Business_jobs.json":          "Business",
		"Software Engineer_jobs.json": "Tech",
	}

	for _, filePath := range filePaths {
		data, err := LoadDataFromJSON(filePath)
		if err != nil {
			fmt.Println("Error loading data from:", filePath, err)
			continue
		}

		// Extracting domain from file name
		segments := strings.Split(filePath, "\\")
		fileName := segments[len(segments)-1]
		domain := domainMapping[fileName]

		// Predict best matching job for each domain
		classifier.PredictBestMatchingJob(domain, data)
	}
}

//
//import (
//	"encoding/json"
//	"fmt"
//	"io"
//	"math"
//	"os"
//	"regexp"
//	"sort"
//	"strings"
//	"time"
//)
//
//// List of common stop words to be excluded from analysis.
//var stopWords = map[string]bool{
//	"and": true, "or": true, "the": true, "in": true,
//	"of": true, "a": true, "is": true, "to": true,
//	"with": true, "for": true, "s": true, "you": true,
//	// More stop words can be added here...
//}
//
//// GenericTextData represents any text data with an associated category.
//type GenericTextData struct {
//	Title       *string `json:"title"`
//	Description *string `json:"description"`
//	Category    string  `json:"domain"`
//}
//
//// NaiveBayesClassifier struct to hold model data.
//type NaiveBayesClassifier struct {
//	wordFrequencies  map[string]map[string]float64
//	categoryCounts   map[string]int
//	totalWords       int
//	totalUniqueWords int
//}
//
//// categoryProb struct for category probabilities.
//type categoryProb struct {
//	Category string
//	Prob     float64
//}
//
//// JSONData structure to match your data format.
//type JSONData struct {
//	Domain string `json:"domain"`
//	Data   []struct {
//		Title       string `json:"title"`
//		Description string `json:"description"`
//	} `json:"data"`
//}
//
//// NewNaiveBayesClassifier creates a new Naive Bayes Classifier.
//func NewNaiveBayesClassifier() *NaiveBayesClassifier {
//	return &NaiveBayesClassifier{
//		wordFrequencies: make(map[string]map[string]float64),
//		categoryCounts:  make(map[string]int),
//	}
//}
//
//// Checks if a word is a stop word.
//func isStopWord(word string) bool {
//	return stopWords[word]
//}
//
//// preprocessText preprocesses the text by converting it to lower case, excluding stop words, and special characters.
//func preprocessText(text string) ([]string, error) {
//	wordRegexp := regexp.MustCompile(`\b[a-zA-Z+#]{4,}\b`) // Only words with 4 or more characters
//	var processedWords []string
//	matches := wordRegexp.FindAllString(text, -1)
//	for _, word := range matches {
//		word = strings.ToLower(word)
//		if !isStopWord(word) {
//			processedWords = append(processedWords, word)
//		}
//	}
//	return processedWords, nil
//}
//
//// CalculateIDF calculates the inverse document frequency for each word across all documents
//func CalculateIDF(corpus []string) map[string]float64 {
//	docCount := make(map[string]int)
//	numDocs := len(corpus)
//
//	for _, text := range corpus {
//		wordSet := make(map[string]bool)
//		words := strings.Fields(text)
//
//		for _, word := range words {
//			wordSet[word] = true
//		}
//
//		for word := range wordSet {
//			docCount[word]++
//		}
//	}
//
//	idf := make(map[string]float64)
//	for word, count := range docCount {
//		idf[word] = math.Log(float64(numDocs) / float64(count))
//	}
//
//	return idf
//}
//
//// CalculateTermFrequency calculates the frequency of each word in a document
//func CalculateTermFrequency(text string) map[string]float64 {
//	wordFreq := make(map[string]int)
//	totalWords := 0
//	words := strings.Fields(text)
//
//	for _, word := range words {
//		wordFreq[word]++
//		totalWords++
//	}
//
//	tf := make(map[string]float64)
//	for word, count := range wordFreq {
//		tf[word] = float64(count) / float64(totalWords)
//	}
//
//	return tf
//}
//
//// TrainWithTFIDF updates the Naive Bayes Classifier with TF-IDF features
//func (nbc *NaiveBayesClassifier) TrainWithTFIDF(data []GenericTextData) {
//	corpus := make([]string, len(data))
//	for i, item := range data {
//		corpus[i] = *item.Description
//	}
//
//	idf := CalculateIDF(corpus)
//
//	for _, item := range data {
//		category := item.Category
//		tf := CalculateTermFrequency(*item.Description)
//
//		for word, freq := range tf {
//			tfidf := freq * idf[word]
//
//			if nbc.wordFrequencies[category] == nil {
//				nbc.wordFrequencies[category] = make(map[string]float64)
//			}
//
//			nbc.wordFrequencies[category][word] += tfidf
//		}
//
//		nbc.categoryCounts[category]++
//	}
//}
//
//// PredictWithTFIDF predicts the category based on TF-IDF features
//func (nbc *NaiveBayesClassifier) PredictWithTFIDF(text string) []categoryProb {
//	tf := CalculateTermFrequency(text)
//	idf := CalculateIDF([]string{text})
//
//	categoryProbabilities := make(map[string]float64)
//	for category := range nbc.categoryCounts {
//		prob := math.Log(float64(nbc.categoryCounts[category]) / float64(len(nbc.categoryCounts)))
//
//		for word, freq := range tf {
//			tfidf := freq * idf[word]
//			prob += math.Log((nbc.wordFrequencies[category][word] + 1.0) * tfidf)
//		}
//
//		categoryProbabilities[category] = math.Exp(prob)
//	}
//
//	var sortedCategories []categoryProb
//	for category, prob := range categoryProbabilities {
//		sortedCategories = append(sortedCategories, categoryProb{category, prob})
//	}
//
//	sort.Slice(sortedCategories, func(i, j int) bool {
//		return sortedCategories[i].Prob > sortedCategories[j].Prob
//	})
//
//	return sortedCategories
//}
//
//// LoadDataFromJSON function updated to extract only title, description, and category.
//func LoadDataFromJSON(filename string) ([]GenericTextData, error) {
//	var jsonData []JSONData
//	var data []GenericTextData
//
//	file, err := os.Open(filename)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
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
//	for _, jsonItem := range jsonData {
//		for _, item := range jsonItem.Data {
//			data = append(data, GenericTextData{
//				Title:       &item.Title,
//				Description: &item.Description,
//				Category:    jsonItem.Domain,
//			})
//		}
//	}
//
//	return data, nil
//}
//
//func main() {
//	startTime := time.Now()
//
//	jsonFile := "C:\\Users\\mathe\\GolandProjects\\cmpscfa23team2\\crab\\output\\combined_jobs.json"
//
//	// Load data from JSON file.
//	combinedData, err := LoadDataFromJSON(jsonFile)
//	if err != nil {
//		fmt.Println("Error loading data:", err)
//		os.Exit(1)
//	}
//
//	// Split data into training and testing.
//	trainDataSize := int(float64(len(combinedData)) * 0.7)
//	trainData := combinedData[:trainDataSize]
//	testData := combinedData[trainDataSize:]
//
//	// Train classifier.
//	classifier := NewNaiveBayesClassifier()
//	classifier.TrainWithTFIDF(trainData)
//
//	// Test classifier.
//	correctPredictions := 0
//	for _, data := range testData {
//		if data.Description != nil {
//			predictedCategories := classifier.PredictWithTFIDF(*data.Description)
//			// Assuming the most probable category is the first in the sorted list
//			if len(predictedCategories) > 0 && predictedCategories[0].Category == data.Category {
//				correctPredictions++
//			}
//		}
//	}
//	accuracy := float64(correctPredictions) / float64(len(testData))
//	fmt.Printf("Accuracy: %.2f%%\n", accuracy*100)
//
//	elapsedTime := time.Since(startTime)
//	fmt.Printf("Execution time: %s\n", elapsedTime)
//}
