package ML

import (
	"cmpscfa23team2/cuda/ML"
	"fmt"
	"reflect"
	"testing"
)

// TestNewNaiveBayesClassifier tests the NewNaiveBayesClassifier function.
func TestNewNaiveBayesClassifier(t *testing.T) {
	expectedClassifier := &ML.NaiveBayesClassifier{
		WordFrequencies:  make(map[string]map[string]int),
		CategoryCounts:   make(map[string]int),
		TotalWords:       0,
		TotalUniqueWords: 0,
		SkillSets: map[string][]string{
			"Tech":     {"Software", "Java", "React", "C++", "JavaScript", "DevOps", "Cloud", "AWS", "Backend", "Frontend", "Full Stack", "Angular", "Node.js", "SQL", "NoSQL", "Git", "Linux", "Embedded", "API", "Microservices"},
			"Business": {"Management", "Finance", "Marketing", "Sales", "Microsoft Office"},
			"Law":      {"Law", "Litigation", "Legal", "Contract", "Compliance"},
		},
	}

	actualClassifier := ML.NewNaiveBayesClassifier()

	if !reflect.DeepEqual(expectedClassifier, actualClassifier) {
		t.Errorf("Expected %+v, got %+v", expectedClassifier, actualClassifier)
	}

}

func TestNaiveBayesClassifier_Train(t *testing.T) {
	// Create an instance of NaiveBayesClassifier
	classifier := ML.NewNaiveBayesClassifier()

	// Create a sample dataset for training
	sampleData := []ML.JobData{
		// Add job data here for training
		{Title: "Software", Description: "Java React"},
		// Add more job data as needed
	}

	// Train the classifier
	classifier.Train(sampleData, "Tech")

	// Add assertions to check if the classifier has been trained correctly
	// For example, you can check if word frequencies, category counts, etc., are as expected

	// Example assertion (modify based on your implementation)
	expectedWordFrequencies := map[string]map[string]int{
		"Tech": {"software": 1, "java": 1, "react": 1},
	}

	// Check if the WordFrequencies map is not nil
	if classifier.WordFrequencies == nil {
		t.Error("WordFrequencies map is nil after training")
	}

	// Check if the WordFrequencies map for "Tech" is as expected
	if !reflect.DeepEqual(classifier.WordFrequencies["Tech"], expectedWordFrequencies["Tech"]) {
		t.Errorf("Unexpected word frequencies after training. Got %v, expected %v", classifier.WordFrequencies["Tech"], expectedWordFrequencies["Tech"])
	}
}

func TestLoadDataFromJSON(t *testing.T) {
	// Create a sample JSON file for testing
	// Ensure that the content of the JSON file is known for testing

	// Load data from the sample JSON file
	container, err := ML.LoadDataFromJSON("C:\\Users\\User\\GolandProjects\\cmpscfa23team2\\Nbc_output\\SoftwareEng_top_jobs.json")
	if err != nil {
		t.Fatalf("Error loading data from JSON: %v", err)
	}
	fmt.Printf("Loaded domain: %s\n", container.Domain)

	// Add assertions to check if the loaded data is as expected
	// For example, you can check the values of container.Domain, container.Data, etc.

	// Example assertion (modify based on your implementation)
	//expectedDomain := "job-market"
	//if container.Domain != expectedDomain {
	//	t.Errorf("Unexpected domain after loading JSON. Got %s, expected %s", container.Domain, expectedDomain)
	//}
}

func TestNaiveBayesClassifier_PredictBestMatchingJob(t *testing.T) {
	// Create an instance of NaiveBayesClassifier
	classifier := ML.NewNaiveBayesClassifier()

	// Create a sample dataset for training
	//var data = ML.JobData{
	//	Title:       "SoftwareEng",
	//	URL:         "SampleURL.com",
	//	Description: "Wut",
	//	Salary:      "1,000,000",
	//	Company:     "Google",
	//	Location:    "Florida",
	//}
	var sampleData []ML.JobData

	// Train the classifier
	classifier.Train(sampleData, "SoftwareEng")

	// Create a sample dataset for predicting
	var predictData []ML.JobData

	// Predict the best matching jobs
	result := classifier.PredictBestMatchingJob("SoftwareEng", predictData)

	// Add assertions to check if the result is as expected
	// For example, you can check if the top job titles match your expectations

	// Example assertion (modify based on your implementation)
	expectedTitles := []string{"Software Engineer", "Full Stack Developer", "DevOps Engineer"}
	if !reflect.DeepEqual(result, expectedTitles) {
		t.Errorf("Unexpected result after predicting best matching jobs. Got %v, expected %v", result, expectedTitles)
	}
}

func TestSearchJobByTitle(t *testing.T) {
	// Create a sample dataset for searching
	var sampleData []ML.JobData

	// Search for a specific job title
	ML.SearchJobByTitle(sampleData, "Software Engineer")

	// Add assertions or checks to verify if the job details are printed as expected
	// You may need to capture the standard output for more advanced testing
}
