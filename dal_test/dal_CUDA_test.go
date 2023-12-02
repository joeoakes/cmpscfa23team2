package dal_test

import (
	"cmpscfa23team2/dal"
	"testing"
)

//func TestEngineIDExists(t *testing.T) {
//	engineID := "test_engine_id2"
//	// Insert a sample engine ID for testing
//	err := dal.InsertSampleEngine(engineID, "Test Engine", "test engine description")
//	if err != nil {
//		dal.InsertLog("400", "Failed to insert sample engine for testing", "TestEngineIDExists()")
//		t.Fatalf("Failed to insert sample engine for testing: %v", err)
//	} else {
//		dal.InsertLog("200", "Successfully inserted sample engine for testing", "TestEngineIDExists()")
//	}
//
//	exists, err := dal.EngineIDExists(engineID)
//	if err != nil {
//		dal.InsertLog("400", "Error checking engine ID", "TestEngineIDExists()")
//		t.Fatalf("Error checking engine ID: %v", err)
//	} else {
//		dal.InsertLog("200", "Successfully checked engine ID", "TestEngineIDExists()")
//	}
//
//	if !exists {
//		t.Error("Expected engine ID to exist, but it does not.")
//	}
//}

//	func TestInsertPrediction(t *testing.T) {
//		engineID := "test_engine23"
//		// Insert a sample engine ID for testing
//		err := dal.InsertSampleEngine(engineID, "test_engine2", "Test Engine Description")
//		if err != nil {
//			dal.InsertLog("400", "Failed to insert sample engine for testing", "TestInsertPrediction()")
//			t.Fatalf("Failed to insert sample engine for testing: %v", err)
//		} else {
//			dal.InsertLog("200", "Successfully inserted sample engine for testing", "TestInsertPrediction()")
//		}
//
//		mlPrediction := dal.PerformMLPrediction("test_data")
//		predictionJSON, err := dal.ConvertPredictionToJSON(mlPrediction)
//		if err != nil {
//			dal.InsertLog("400", "Failed to convert prediction to JSON", "TestInsertPrediction()")
//			t.Fatalf("Failed to convert prediction to JSON: %v", err)
//		} else {
//			dal.InsertLog("200", "Successfully converted prediction to JSON", "TestInsertPrediction()")
//		}
//
//		err = dal.InsertPrediction(engineID, predictionJSON)
//		if err != nil {
//			dal.InsertLog("400", "Failed to insert prediction", "TestInsertPrediction()")
//			t.Fatalf("Failed to insert prediction: %v", err)
//		} else {
//			dal.InsertLog("200", "Successfully inserted prediction", "TestInsertPrediction()")
//		}
//	}
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
