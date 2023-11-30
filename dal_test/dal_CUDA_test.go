package dal

import (
	"cmpscfa23team2/dal"
	"testing"
)

func TestEngineIDExists(t *testing.T) {

	engineID := "test_engine_id"
	// Insert a sample engine ID for testing
	err := dal.InsertSampleEngine(engineID, "Test Engine", "test engine description")
	if err != nil {
		t.Fatalf("Failed to insert sample engine for testing: %v", err)
	}

	exists, err := dal.EngineIDExists(engineID)
	if err != nil {
		t.Fatalf("Error checking engine ID: %v", err)
	}

	if !exists {
		t.Error("Expected engine ID to exist, but it does not.")
	}
}

// func TestInsertPrediction(t *testing.T) {
//
//	engineID := "test_engine2"
//	// Insert a sample engine ID for testing
//	err := dal.InsertSampleEngine(engineID, "test_engine2", "Test Engine Description")
//	if err != nil {
//		t.Fatalf("Failed to insert sample engine for testing: %v", err)
//	}
//	mlPrediction := dal.PerformMLPrediction("test_data")
//	predictionJSON, err := dal.ConvertPredictionToJSON(mlPrediction)
//	if err != nil {
//		t.Fatalf("Failed to convert prediction to JSON: %v", err)
//	}
//
//	err = dal.InsertPrediction(engineID, predictionJSON)
//	if err != nil {
//		t.Fatalf("Failed to insert prediction: %v", err)
//	}
//
// }

func TestInsertPrediction(t *testing.T) {
	// Mock file name and text-based prediction data
	mockFileName := "sampleTextFile.txt"
	mockTextPrediction := "This is a sample text-based prediction."

	// Call the InsertPrediction function with text-based prediction and file name
	err := dal.InsertPrediction(mockFileName, mockTextPrediction)
	if err != nil {
		t.Errorf("InsertPrediction with text failed: %v", err)
	}

	// Mock image file name and file path prediction data
	mockImageFileName := "predictedImage.png"
	mockImageFilePath := "/path/to/predicted/image.png"

	// Call the InsertPrediction function with image file path and file name
	err = dal.InsertPrediction(mockImageFileName, mockImageFilePath)
	if err != nil {
		t.Errorf("InsertPrediction with image file path failed: %v", err)
	}

	// Additional validation checks can be added here if necessary
	// For example, querying the database to ensure data was inserted correctly
}
