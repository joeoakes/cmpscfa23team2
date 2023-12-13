package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
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

// SearchJobByTitle searches for a job by its title and prints its details.
func SearchJobByTitle(data []JobData, title string) {
	for _, job := range data {
		if job.Title == title {
			printJobDetails(job)
			return
		}
	}
	fmt.Println("Job title not found:", title)
}

// printJobDetails prints the details of a job.
func printJobDetails(job JobData) {
	fmt.Printf("Title: %s\nURL: %s\nCompany: %s\nLocation: %s\nSalary: %s\nDescription: %s\n\n",
		job.Title, job.URL, job.Company, job.Location, job.Salary, job.Description)
}

func main1() {
	classifier := NewNaiveBayesClassifier()

	// File paths
	filePaths := []string{
		"C:\\Users\\Public\\GoLandProjects\\JustAFork\\crab\\output\\Law_jobs.json",
		"C:\\Users\\Public\\GoLandProjects\\JustAFork\\crab\\output\\Business_jobs.json",
		"C:\\Users\\Public\\GoLandProjects\\JustAFork\\crab\\output\\SoftwareEng_jobs.json",
	}

	// Mapping of file names to their respective domains
	domainMapping := map[string]string{
		"Law_jobs.json":         "Law",
		"Business_jobs.json":    "Business",
		"SoftwareEng_jobs.json": "Tech",
	}

	for _, filePath := range filePaths {
		fmt.Printf("Processing file: %s\n", filePath)

		data, err := LoadDataFromJSON(filePath)
		if err != nil {
			fmt.Printf("Error loading data from %s: %v\n", filePath, err)
			continue
		}

		fileName := filepath.Base(filePath)
		domain, ok := domainMapping[fileName]
		if !ok {
			fmt.Printf("Domain not found for file %s\n", fileName)
			continue
		}
		fmt.Printf("Predicting best matching job for domain '%s'...\n", domain)
		classifier.PredictBestMatchingJob(domain, data)
		fmt.Println("--------------------------------------------------\n\n")

		// Example job title to search for
		jobTitle := "Software Release DevOps Engineer" // replace with the title you want to search for
		fmt.Printf("Searching for job title: %s\n", jobTitle)
		SearchJobByTitle(data, jobTitle)
	}
}
