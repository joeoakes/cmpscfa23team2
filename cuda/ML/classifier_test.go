package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestTrain(t *testing.T) {
	classifier := NewNaiveBayesClassifier()
	data := []GenericTextData{
		{Description: strPtr("python programming"), Category: "Technology"},
		{Description: strPtr("artificial intelligence"), Category: "Technology"},
	}

	classifier.Train(data)

	techWords := classifier.wordFrequencies["Technology"]
	if techWords["python"] != 1 {
		t.Errorf("Expected 'python' to have frequency 1, got %d", techWords["python"])
	}
	if techWords["artificial"] != 1 {
		t.Errorf("Expected 'artificial' to have frequency 1, got %d", techWords["artificial"])
	}
}

func TestPreprocessText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "standard case",
			input:    "The quick brown fox",
			expected: []string{"quick", "brown", "fox"},
		},
		{
			name:     "with punctuation",
			input:    "Hello, world!",
			expected: []string{"hello", "world"},
		},
		{
			name:     "with numbers",
			input:    "123 go",
			expected: []string{"123", "go"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{}, // Corrected expectation
		},
		{
			name:     "only stop words",
			input:    "and or the in of a is to",
			expected: []string{}, // Corrected expectation
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := preprocessText(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("preprocessText(%v) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestPredictWithProbabilities(t *testing.T) {
	classifier := NewNaiveBayesClassifier()
	data := []GenericTextData{
		{Description: strPtr("python and data analysis"), Category: "Technology"},
		{Description: strPtr("modern art and design"), Category: "Art"},
	}

	classifier.Train(data)

	desc := "python machine learning"
	predictedCategory, _ := classifier.PredictWithProbabilities(&desc)
	if predictedCategory != "Technology" {
		t.Errorf("Predicted category for '%s' is %s, want %s", desc, predictedCategory, "Technology")
	}

	desc = "abstract paintings"
	predictedCategory, _ = classifier.PredictWithProbabilities(&desc)
	if predictedCategory != "Art" {
		t.Errorf("Predicted category for '%s' is %s, want %s", desc, predictedCategory, "Art")
	}

	desc = "" // Test with empty description
	predictedCategory, _ = classifier.PredictWithProbabilities(&desc)
	if predictedCategory != "" {
		t.Errorf("Predicted category for empty description is %s, want empty string", predictedCategory)
	}
}
func TestLoadDataFromJSON(t *testing.T) {
	// Create a temporary JSON file with test data
	content := `{"items": [{"domain": "Technology", "data": {"description": "Machine learning and data science"}}]}`
	tmpfile, err := ioutil.TempFile("", "test*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test LoadDataFromJSON
	data, err := LoadDataFromJSON(tmpfile.Name())
	if err != nil {
		t.Errorf("LoadDataFromJSON() error = %v", err)
		return
	}
	if len(data) != 1 || *data[0].Description != "Machine learning and data science" {
		t.Errorf("LoadDataFromJSON() = %v, want %v", data, content)
	}
}

// Additional test functions can be added as needed
