package dal_test

import (
	"cmpscfa23team2/dal"
	"testing"
)

func TestEngineIDExists(t *testing.T) {
	engineID := "test_engine_id2"
	// Insert a sample engine ID for testing
	err := dal.InsertSampleEngine(engineID, "Test Engine", "test engine description")
	if err != nil {
		dal.InsertLog("400", "Failed to insert sample engine for testing", "TestEngineIDExists()")
		t.Fatalf("Failed to insert sample engine for testing: %v", err)
	} else {
		dal.InsertLog("200", "Successfully inserted sample engine for testing", "TestEngineIDExists()")
	}

	exists, err := dal.EngineIDExists(engineID)
	if err != nil {
		dal.InsertLog("400", "Error checking engine ID", "TestEngineIDExists()")
		t.Fatalf("Error checking engine ID: %v", err)
	} else {
		dal.InsertLog("200", "Successfully checked engine ID", "TestEngineIDExists()")
	}

	if !exists {
		t.Error("Expected engine ID to exist, but it does not.")
	}
}

func TestInsertPrediction(t *testing.T) {
	engineID := "test_engine23"
	// Insert a sample engine ID for testing
	err := dal.InsertSampleEngine(engineID, "test_engine2", "Test Engine Description")
	if err != nil {
		dal.InsertLog("400", "Failed to insert sample engine for testing", "TestInsertPrediction()")
		t.Fatalf("Failed to insert sample engine for testing: %v", err)
	} else {
		dal.InsertLog("200", "Successfully inserted sample engine for testing", "TestInsertPrediction()")
	}

	mlPrediction := dal.PerformMLPrediction("test_data")
	predictionJSON, err := dal.ConvertPredictionToJSON(mlPrediction)
	if err != nil {
		dal.InsertLog("400", "Failed to convert prediction to JSON", "TestInsertPrediction()")
		t.Fatalf("Failed to convert prediction to JSON: %v", err)
	} else {
		dal.InsertLog("200", "Successfully converted prediction to JSON", "TestInsertPrediction()")
	}

	err = dal.InsertPrediction(engineID, predictionJSON)
	if err != nil {
		dal.InsertLog("400", "Failed to insert prediction", "TestInsertPrediction()")
		t.Fatalf("Failed to insert prediction: %v", err)
	} else {
		dal.InsertLog("200", "Successfully inserted prediction", "TestInsertPrediction()")
	}
}
