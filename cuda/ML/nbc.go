package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

// DomainResult represents the result for a domain with top job matches and skill demand.
type DomainResult struct {
	Domain      string               `json:"domain"`
	SkillDemand map[string]SkillData `json:"skill_demand"`
	TopJobs     []jobMatch           `json:"top_jobs"`
}

// NewNaiveBayesClassifier creates a new Naive Bayes Classifier.
//
//	func NewNaiveBayesClassifier() *NaiveBayesClassifier {
//		return &NaiveBayesClassifier{
//			wordFrequencies:  make(map[string]map[string]int),
//			categoryCounts:   make(map[string]int),
//			totalWords:       0,
//			totalUniqueWords: 0,
//			skillSets: map[string][]string{
//				"Tech":     {"Python", "Java", "React"},
//				"Business": {"Management", "Finance", "marketing", "management", "Microsoft Office"},
//				"Law":      {"school of law", "Litigation"},
//			},
//		}
//	}
func NewNaiveBayesClassifier() *NaiveBayesClassifier {
	return &NaiveBayesClassifier{
		wordFrequencies:  make(map[string]map[string]int),
		categoryCounts:   make(map[string]int),
		totalWords:       0,
		totalUniqueWords: 0,
		skillSets: map[string][]string{
			"SoftwareEng": {"Software", "Java", "React", "C++", "JavaScript", "DevOps", "Cloud", "AWS", "Backend", "Frontend", "Full Stack", "Angular", "Node.js", "SQL", "NoSQL", "Git", "Linux", "Embedded", "API", "Microservices"},
			"Business":    {"Management", "Finance", "Marketing", "Sales", "Microsoft Office"},
			"Law":         {"Law", "Litigation", "Legal", "Contract", "Compliance"},
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
//func (nbc *NaiveBayesClassifier) PredictWithProbabilities(skills []string) []string {
//	// Assuming implementation of probability calculation
//	// Returns a list of categories sorted by their probability score
//	return []string{"Category1", "Category2"} // Placeholder
//}

// LoadDataFromJSON updated to match the new GenericTextData structure.
func LoadDataFromJSON(filename string) (JobDataContainer, error) {
	var container JobDataContainer
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return container, err
	}

	err = json.Unmarshal(file, &container)
	if err != nil {
		return container, err
	}

	return container, nil
}

type jobMatch struct {
	Job           JobData  `json:"job"`
	Count         int      `json:"count"`
	MatchedSkills []string `json:"matched_skills"`
}

func (nbc *NaiveBayesClassifier) PredictBestMatchingJob(domain string, data []JobData) []string {
	skills := nbc.skillSets[domain]
	jobSkills := make(map[string][]string)

	for _, job := range data {
		jobDesc := strings.ToLower(job.Description)
		jobTitle := strings.ToLower(job.Title)

		for _, skill := range skills {
			skillLower := strings.ToLower(skill)
			if strings.Contains(jobDesc, skillLower) || strings.Contains(jobTitle, skillLower) {
				jobSkills[job.URL] = append(jobSkills[job.URL], skill)
			}
		}
	}

	if len(jobSkills) == 0 {
		fmt.Println("No jobs matched for the domain:", domain)
		return []string{}
	}

	// Sorting and picking top jobs
	var topJobs []jobMatch
	for url, matchedSkills := range jobSkills {
		for _, job := range data {
			if job.URL == url {
				topJobs = append(topJobs, jobMatch{Job: job, Count: len(matchedSkills), MatchedSkills: matchedSkills})
				break
			}
		}
	}
	sort.Slice(topJobs, func(i, j int) bool {
		return topJobs[i].Count > topJobs[j].Count
	})

	var topJobTitles []string
	for i, match := range topJobs {
		if i == 3 {
			break
		}
		topJobTitles = append(topJobTitles, match.Job.Title)
	}

	return topJobTitles
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

func main() {
	// Define relative paths to the JSON files - later
	//basePath := "../../output/" // Adjust this path according to your directory structure
	//filePaths := []string{
	//	filepath.Join(basePath, "Law_jobs.json"),
	//	filepath.Join(basePath, "Business_jobs.json"),
	//	filepath.Join(basePath, "SoftwareEng_jobs.json"),
	//}

	classifier := NewNaiveBayesClassifier()
	filePaths := []string{
		"C:\\Users\\Public\\GoLandProjects\\JustAFork\\crab\\output\\Law_jobs.json",
		"C:\\Users\\Public\\GoLandProjects\\JustAFork\\crab\\output\\SoftwareEng_jobs.json",
		"C:\\Users\\Public\\GoLandProjects\\JustAFork\\crab\\output\\Business_jobs.json",
	}

	//domainMapping := map[string]string{
	//	"Law_jobs.json":         "Law",
	//	"Business_jobs.json":    "Business",
	//	"SoftwareEng_jobs.json": "Tech",
	//}

	// Create Nbc_output directory if it doesn't exist
	//outputDir := filepath.Join(basePath, "Nbc_output")
	//if _, err := os.Stat(outputDir); os.IsNotExist(err) {
	//	os.Mkdir(outputDir, os.ModePerm)
	//}

	outputDir := "Nbc_output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	for _, filePath := range filePaths {
		container, err := LoadDataFromJSON(filePath)
		if err != nil {
			fmt.Printf("Error loading data from %s: %v\n", filePath, err)
			continue
		}

		domain := filepath.Base(filePath)
		domain = strings.Split(domain, "_")[0] // Extracting domain name from filename

		topJobTitles := classifier.PredictBestMatchingJob(domain, container.Data)
		var topJobsData []JobData
		for _, title := range topJobTitles {
			for _, job := range container.Data {
				if job.Title == title {
					topJobsData = append(topJobsData, job)
					break
				}
			}
		}

		result := JobDataContainer{
			Domain:   domain,
			URL:      container.URL, // Use the URL from the loaded JSON data
			Data:     topJobsData,
			Metadata: container.Metadata,
		}

		resultJSON, _ := json.MarshalIndent(result, "", "    ")
		outputFilename := filepath.Join(outputDir, domain+"_top_jobs.json")
		_ = ioutil.WriteFile(outputFilename, resultJSON, 0644)

		fmt.Printf("Top 3 jobs for '%s' domain written to %s\n", domain, outputFilename)
		fmt.Println("--------------------------------------------------\n\n")
	}
}
