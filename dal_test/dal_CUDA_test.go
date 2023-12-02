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

// TestInsertPrediction tests inserting various types of predictions for different algorithms.
func TestInsertPrediction(t *testing.T) {
	// Sample query identifier
	queryIdentifier := "sample_query_id"

	// Test cases for KNN algorithm, keep the same path format if you want to add pictures
	if err := dal.InsertPrediction("KNN", "Real Estate Query 1 Philadelphia", "realEstateQuery1Pic.png", "LinearRegression/realEstateQuery1Pic.png"); err != nil {
		t.Errorf("Failed to insert text prediction for KNN: %v", err)
	}
	if err := dal.InsertPrediction("KNN", "Real Estate Query 2 New York", "newyork.txt", "Text prediction for New york"); err != nil {
		t.Errorf("Failed to insert image path prediction for KNN: %v", err)
	}

	// Test cases for Linear Regression algorithm
	if err := dal.InsertPrediction("LinearRegression", queryIdentifier, "text_prediction_lr.txt", "Sample text prediction for Linear Regression"); err != nil {
		t.Errorf("Failed to insert text prediction for Linear Regression: %v", err)
	}
	if err := dal.InsertPrediction("LinearRegression", queryIdentifier, "image_lr.png", "/path/to/image_lr.png"); err != nil {
		t.Errorf("Failed to insert image path prediction for Linear Regression: %v", err)
	}

	// Test cases for Naive Bayes algorithm
	if err := dal.InsertPrediction("NaiveBayes", queryIdentifier, "text_prediction_nb.txt", "Sample text prediction for Naive Bayes"); err != nil {
		t.Errorf("Failed to insert text prediction for Naive Bayes: %v", err)
	}
	if err := dal.InsertPrediction("NaiveBayes", queryIdentifier, "image_nb.png", "/path/to/image_nb.png"); err != nil {
		t.Errorf("Failed to insert image path prediction for Naive Bayes: %v", err)
	}

	// Test case for an invalid algorithm
	if err := dal.InsertPrediction("InvalidAlgorithm", queryIdentifier, "invalid_data.txt", "This should fail"); err == nil {
		t.Error("Expected an error for an unrecognized algorithm, but got none")
	}
}
