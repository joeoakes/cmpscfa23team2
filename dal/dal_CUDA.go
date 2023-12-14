package dal

// Import required packages
import (
	"database/sql"
	"encoding/json"                    // For JSON handling
	"fmt"                              // For formatted I/O
	_ "github.com/go-sql-driver/mysql" // Import mysql driver
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"io/ioutil"
	"log" // For logging
	"os"
	"reflect"
	"testing"
)

// Prediction struct models the data structure of a prediction in the database
//
// This code defines a struct named "Prediction" with fields for PredictionID, EngineID, InputData, PredictionInfo, and PredictionTime.
type Prediction struct {
	PredictionID   string
	EngineID       string
	InputData      string
	PredictionInfo string
	PredictionTime string
}

// PredictionData represents the structure of the prediction data
type PredictionData struct {
	PredictionInfo string    `json:"prediction_info"`
	InputData      string    `json:"input_data"`
	ImagePath      string    `json:"image_path"`
	Skills         string    `json:"skills"`
	JobListings    []JobData `json:"job_listings"`
	SpecificJob    *JobData  `json:"specific_job,omitempty"` // Add this line
}

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

func InsertPrediction(algorithm, queryIdentifier, fileName, predictionInfo, skills string) error {
	// Generate a new UUID for the prediction
	newUUID := uuid.New().String()

	var query string
	switch algorithm {
	case "KNN":
		query = "INSERT INTO knn_predictions (prediction_id, query_identifier, input_data, prediction_info) VALUES (?, ?, ?, ?)"
	case "LinearRegression":
		query = "INSERT INTO linear_regression_predictions (prediction_id, query_identifier, input_data, prediction_info) VALUES (?, ?, ?, ?)"
	case "NaiveBayes":
		query = "INSERT INTO naive_bayes_predictions (prediction_id, query_identifier, input_data, prediction_info) VALUES (?, ?, ?, ?)"
	default:
		return fmt.Errorf("Unrecognized algorithm: %v", algorithm)
	}

	_, err := DB.Exec(query, newUUID, queryIdentifier, skills, predictionInfo)
	if err != nil {
		return fmt.Errorf("Error storing prediction for %v: %v", algorithm, err)
	}

	log.Printf("Successfully inserted prediction with ID %s for %v algorithm.", newUUID, algorithm)
	return nil
}

// Simulated ML model prediction function
//
// It definesa function that simulates an ML model prediction with a 2-second delay
// and logs a success message before returning a prediction result as a formatted string.
//func PerformMLPrediction(inputData string) string {
//	// Simulate some delay for ML model prediction
//	time.Sleep(2 * time.Second)
//	log.Println("Successfully performed ML prediction.")
//	return fmt.Sprintf("Prediction result for %s", inputData)
//}

// Convert prediction result to JSON
//
// defines a function that converts a given prediction result string into a JSON format, logging a success message and returning the JSON string or an error.
func ConvertPredictionToJSON(predictionResult string) (string, error) {
	predictionMap := map[string]string{"result": predictionResult}
	predictionJSON, err := json.Marshal(predictionMap)
	if err != nil {
		InsertLog("400", "Error converting prediction to JSON: "+err.Error(), "ConvertPredictionToJSON()")
		return "", err
	} else {
		InsertLog("200", "Successfully converted prediction to JSON.", "ConvertPredictionToJSON()")
		log.Println("Successfully converted prediction to JSON.")
	}
	return string(predictionJSON), nil
}

// LoadDataFromJSON updated to match the new GenericTextData structure.
//func LoadDataFromJSON(filename string) ([]JobData, error) {
//	file, err := ioutil.ReadFile(filename)
//	if err != nil {
//		return nil, err
//	}
//
//	var container JobDataContainer
//	err = json.Unmarshal(file, &container)
//	if err != nil {
//		return nil, err
//	}
//
//	return container.Data, nil
//}

// Updated LoadDataFromJSON function
func LoadDataFromJSON(filename string, specificJobTitle string) ([]JobData, *JobData, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	var container JobDataContainer
	err = json.Unmarshal(file, &container)
	if err != nil {
		return nil, nil, err
	}

	var specificJob *JobData
	for _, job := range container.Data {
		if job.Title == specificJobTitle {
			specificJob = &job
			break
		}
	}

	return container.Data, specificJob, nil
}

type jobMatch struct {
	job           JobData
	count         int
	matchedSkills []string
}

func constructImagePath(queryIdentifier, domain string) string {
	// This is just an example. You need to modify it based on your actual file structure and requirements.
	basePath := "/static/images/"
	return basePath + domain + "/" + queryIdentifier + ".png"
}

// FetchPredictionData function
//
//	func FetchPredictionData(queryIdentifier, domain string) (PredictionData, error) {
//		var data PredictionData
//		var predictionPath string
//
//		queryStr := "SELECT input_data, prediction_info FROM naive_bayes_predictions WHERE query_identifier = ?"
//		err := DB.QueryRow(queryStr, queryIdentifier).Scan(&data.Skills, &predictionPath)
//		if err != nil {
//			if err == sql.ErrNoRows {
//				return PredictionData{}, fmt.Errorf("no prediction data found for query identifier: %s", queryIdentifier)
//			}
//			return PredictionData{}, err
//		}
//
//		// Read JSON file from the prediction path
//		file, err := ioutil.ReadFile(predictionPath)
//		if err != nil {
//			return PredictionData{}, fmt.Errorf("error reading JSON file: %s", err)
//		}
//
//		var container JobDataContainer
//		if err := json.Unmarshal(file, &container); err != nil {
//			return PredictionData{}, fmt.Errorf("error parsing JSON data: %s", err)
//		}
//
//		// Update data struct
//		data.JobListings = container.Data
//
//		// Example: Search for a specific job title (this can be dynamic based on user input)
//		specificJobTitle := "Software Release DevOps Engineer" // Replace with dynamic input as needed
//		data.SpecificJob = SearchJobByTitle(container.Data, specificJobTitle)
//
//		// Construct the image path (if applicable)
//		//data.ImagePath = constructImagePath(queryIdentifier, domain) // Implement this function as needed
//
//		return data, nil
//	}
//
//func FetchPredictionData(queryIdentifier, domain string) (PredictionData, error) {
//	var data PredictionData
//	var predictionPath, jobTitle string
//
//	// Fetch job title and JSON path from the database
//	queryStr := "SELECT input_data, prediction_info FROM naive_bayes_predictions WHERE query_identifier = ?"
//	err := DB.QueryRow(queryStr, queryIdentifier).Scan(&jobTitle, &predictionPath)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return PredictionData{}, fmt.Errorf("no prediction data found for query identifier: %s", queryIdentifier)
//		}
//		return PredictionData{}, err
//	}
//
//	// Read and parse the JSON file
//	file, err := ioutil.ReadFile(predictionPath)
//	if err != nil {
//		return PredictionData{}, fmt.Errorf("error reading JSON file: %s", err)
//	}
//
//	var container JobDataContainer
//	if err := json.Unmarshal(file, &container); err != nil {
//		return PredictionData{}, fmt.Errorf("error parsing JSON data: %s", err)
//	}
//
//	// Filter job data by the title fetched from the database
//	var specificJob *JobData
//	for _, job := range container.Data {
//		if job.Title == jobTitle {
//			specificJob = &job
//			break
//		}
//	}
//
//	// Update data struct
//	data.JobListings = container.Data
//	data.SpecificJob = specificJob // Set the specific job data if a match is found
//
//	// Include any additional logic as needed, such as constructing image paths
//
//	return data, nil
//}

// FetchPredictionData fetches prediction data based on the domain and query identifier
func FetchPredictionData(queryIdentifier, domain string) (PredictionData, error) {
	var data PredictionData
	var queryStr string
	var err error

	switch domain {
	case "E-commerce (Price Prediction)", "Gas Prices (Industry Trend Analysis)":
		queryStr = "SELECT prediction_info FROM linear_regression_predictions WHERE query_identifier = ?"
		err = DB.QueryRow(queryStr, queryIdentifier).Scan(&data.PredictionInfo)
		if err != nil {
			return handleDBError(err, queryIdentifier)
		}
		data.ImagePath = fmt.Sprintf("/static/Assets/MachineLearning/LinearRegression/%s_scatter_plot.png", queryIdentifier)

	case "RealEstate":
		queryStr = "SELECT prediction_info FROM knn_predictions WHERE query_identifier = ?"
		err = DB.QueryRow(queryStr, queryIdentifier).Scan(&data.PredictionInfo)
		if err != nil {
			return handleDBError(err, queryIdentifier)
		}
		data.ImagePath = fmt.Sprintf("/static/Assets/MachineLearning/KNN/%s_scatter_plot.png", queryIdentifier)

	case "Job Market (Industry Trend Analysis)":
		var predictionPath, jobTitle string
		queryStr = "SELECT input_data, prediction_info FROM naive_bayes_predictions WHERE query_identifier = ?"
		err = DB.QueryRow(queryStr, queryIdentifier).Scan(&jobTitle, &predictionPath)
		if err != nil {
			return handleDBError(err, queryIdentifier)
		}

		if _, err := os.Stat(predictionPath); os.IsNotExist(err) {
			return PredictionData{}, fmt.Errorf("JSON file not found at path: %s", predictionPath)
		}

		file, err := ioutil.ReadFile(predictionPath)
		if err != nil {
			return PredictionData{}, fmt.Errorf("error reading JSON file: %s", err)
		}

		var container JobDataContainer
		if err := json.Unmarshal(file, &container); err != nil {
			return PredictionData{}, fmt.Errorf("error parsing JSON data: %s", err)
		}

		data.JobListings = container.Data
		data.SpecificJob = SearchJobByTitle(container.Data, jobTitle)

	default:
		return PredictionData{}, fmt.Errorf("unrecognized domain: %s", domain)
	}

	return data, nil
}

func handleDBError(err error, queryIdentifier string) (PredictionData, error) {
	if err == sql.ErrNoRows {
		return PredictionData{}, fmt.Errorf("no prediction data found for query identifier: %s", queryIdentifier)
	}
	return PredictionData{}, err
}

// SearchJobByTitle searches for a job by its title and returns its details.
func SearchJobByTitle(data []JobData, title string) *JobData {
	for _, job := range data {
		if job.Title == title {
			return &job
		}
	}
	log.Println("Job title not found:", title)
	return nil
}

func TestSearchJobByTitle(t *testing.T) {
	jobs := []JobData{
		{Title: "Software Engineer", URL: "url1", Company: "Company1"},
		{Title: "Software Release DevOps Engineer", URL: "url2", Company: "Company2"},
		// Add more jobs as needed
	}

	tests := []struct {
		title    string
		expected *JobData
	}{
		{"Software Release DevOps Engineer", &jobs[1]},
		{"Non Existent Job", nil},
	}

	for _, test := range tests {
		result := SearchJobByTitle(jobs, test.title)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("SearchJobByTitle(%s) = %v; expected %v", test.title, result, test.expected)
		}
	}
}

func formatJobData(job *JobData) string {
	if job == nil {
		return "nil"
	}
	return fmt.Sprintf("\nTitle: %s\nURL: %s\nCompany: %s\nLocation: %s\nSalary: %s\nDescription: %s\n",
		job.Title, job.URL, job.Company, job.Location, job.Salary, job.Description)
}
