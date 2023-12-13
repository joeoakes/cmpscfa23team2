package dal

// Dal_cuda

// Function to check if the engine_id exists in scraper_engine table
//
// This function checks if a given engine ID exists in a databse table and returns a boolean indicating existence or an error.
//func EngineIDExists(engineID string) (bool, error) {
//  var exists bool
//  query := "SELECT EXISTS(SELECT 1 FROM scraper_engine WHERE engine_id=?)"
//  err := DB.QueryRow(query, engineID).Scan(&exists)
//  if err != nil {
//     InsertLog("400", "Error checking engine ID existence: "+err.Error(), "EngineIDExists()")
//     return false, err
//  } else {
//     InsertLog("200", "Successfully checked if engine ID exists.", "EngineIDExists()")
//     log.Println("Successfully checked if engine ID exists.")
//  }
//  return exists, nil
//}

// Function to insert a new prediction
// The function InsertPrediction, that checks the existence of an engineID, logs the result, and inserts predictionInfo into a database table if the engineID exists, handling errors along the way.

//func InsertPrediction(engineID string, predictionInfo string) error {
//  exists, err := EngineIDExists(engineID)
//  if err != nil {
//     InsertLog("400", "Error checking engine ID: "+err.Error(), "InsertPrediction()")
//     return fmt.Errorf("Error checking engine ID: %v", err)
//  } else {
//     InsertLog("200", "Successfully checked if engine ID exists.", "InsertPrediction()")
//     log.Println("Successfully checked if engine ID exists.")
//  }
//  if !exists {
//     InsertLog("400", "engine_id does not exist", "InsertPrediction()")
//     return fmt.Errorf("engine_id %s does not exist", engineID)
//  } else {
//     InsertLog("200", "Engine ID exists.", "InsertPrediction()")
//     log.Println("Engine ID exists.")
//  }
//
//  query := "INSERT INTO predictions (prediction_id, input_data, prediction_info) VALUES (?, ?, ?)"
//  _, err := DB.Exec(query, newUUID, fileName, predictionInfo)
//  if err != nil {
//     InsertLog("400", "Error storing prediction: "+err.Error(), "InsertPrediction()")
//     return fmt.Errorf("Error storing prediction: %v", err)
//
//  } else {
//     InsertLog("200", "Successfully inserted prediction.", "InsertPrediction()")
//     log.Println("Successfully inserted prediction.")
//  }
//  return nil
//}

// Function to insert a sample engine ID into scraper_engine table
//
// Function inserts a sample engine's information into a database table, logs success, and returns any encountered errors.
//
//	func InsertSampleEngine(engineID, engineName, engineDescription string) error {
//	   query := "INSERT INTO scraper_engine (engine_id, engine_name, engine_description) VALUES (?, ?, ?)"
//	   _, err := DB.Exec(query, engineID, engineName, engineDescription)
//	   if err != nil {
//	      InsertLog("400", "Error inserting sample engine: "+err.Error(), "InsertSampleEngine()")
//	      return fmt.Errorf("Error inserting sample engine: %v", err)
//	   } else {
//	      InsertLog("200", "Successfully inserted sample engine.", "InsertSampleEngine()")
//	      log.Println("Successfully inserted sample engine.")
//	   }
//	   return nil
//	}

//func InsertPrediction(algorithm, queryIdentifier, fileName, predictionInfo string) error {
//	// Generate a new UUID for the prediction
//	newUUID := uuid.New().String()
//
//	var query string
//	switch algorithm {
//	case "KNN":
//		query = "INSERT INTO knn_predictions (prediction_id, query_identifier, input_data, prediction_info) VALUES (?, ?, ?, ?)"
//	case "LinearRegression":
//		query = "INSERT INTO linear_regression_predictions (prediction_id, query_identifier, input_data, prediction_info) VALUES (?, ?, ?, ?)"
//	case "NaiveBayes":
//		query = "INSERT INTO naive_bayes_predictions (prediction_id, query_identifier, input_data, prediction_info) VALUES (?, ?, ?, ?)"
//	default:
//		return fmt.Errorf("Unrecognized algorithm: %v", algorithm)
//	}
//
//	_, err := DB.Exec(query, newUUID, queryIdentifier, fileName, predictionInfo)
//	if err != nil {
//		return fmt.Errorf("Error storing prediction for %v: %v", algorithm, err)
//	}
//
//	log.Printf("Successfully inserted prediction with ID %s for %v algorithm.", newUUID, algorithm)
//	return nil
//}

//func FetchPredictionData(query, domain string) (PredictionData, error) {
//	var (
//		data     PredictionData
//		queryStr string
//		err      error
//	)
//
//	switch domain {
//	case "E-commerce (Price Prediction)":
//		queryStr = "SELECT prediction_info FROM linear_regression_predictions WHERE query_identifier = ?"
//	case "Gas Prices (Industry Trend Analysis)":
//		queryStr = "SELECT prediction_info FROM linear_regression_predictions WHERE query_identifier = ?"
//	case "RealEstate":
//		queryStr = "SELECT prediction_info FROM knn_predictions WHERE query_identifier = ?"
//	case "Job Market (Industry Trend Analysis)":
//		queryStr = "SELECT prediction_info FROM naive_bayes_predictions WHERE query_identifier = ?"
//	default:
//		return PredictionData{}, fmt.Errorf("unrecognized domain: %s", domain)
//	}
//
//	err = DB.QueryRow(queryStr, query).Scan(&data.PredictionInfo)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return PredictionData{}, fmt.Errorf("no prediction data found for query: %s", query)
//		}
//		return PredictionData{}, err
//	}
//	// Construct the image path
//	imagePath := fmt.Sprintf("/static/Assets/MachineLearning/LinearRegression/%s_scatter_plot.png", query)
//
//	// Log the generated image path
//	log.Printf("Generated image path: %s\n", imagePath)
//	// Include the image path in the response
//	data.ImagePath = imagePath
//
//	return data, nil
//}

// test cases for above:

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
//func TestInsertPrediction(t *testing.T) {
//	// Sample query identifier
//	queryIdentifier := "sample_query_id"
//
//	// Test cases for KNN algorithm, keep the same path format if you want to add pictures
//	if err := dal.InsertPrediction("KNN", "Real Estate Query 1 Philadelphia", "realEstateQuery1Pic.png", "LinearRegression/realEstateQuery1Pic.png"); err != nil {
//		t.Errorf("Failed to insert text prediction for KNN: %v", err)
//	}
//	if err := dal.InsertPrediction("KNN", "Real Estate Query 2 New York", "newyork.txt", "Text prediction for New york"); err != nil {
//		t.Errorf("Failed to insert image path prediction for KNN: %v", err)
//	}
//
//	// Test cases for Linear Regression algorithm
//	if err := dal.InsertPrediction("LinearRegression", queryIdentifier, "text_prediction_lr.txt", "Sample text prediction for Linear Regression"); err != nil {
//		t.Errorf("Failed to insert text prediction for Linear Regression: %v", err)
//	}
//	if err := dal.InsertPrediction("LinearRegression", queryIdentifier, "image_lr.png", "/path/to/image_lr.png"); err != nil {
//		t.Errorf("Failed to insert image path prediction for Linear Regression: %v", err)
//	}
//
//	// Test cases for Naive Bayes algorithm
//	if err := dal.InsertPrediction("NaiveBayes", queryIdentifier, "text_prediction_nb.txt", "Sample text prediction for Naive Bayes"); err != nil {
//		t.Errorf("Failed to insert text prediction for Naive Bayes: %v", err)
//	}
//	if err := dal.InsertPrediction("NaiveBayes", queryIdentifier, "image_nb.png", "/path/to/image_nb.png"); err != nil {
//		t.Errorf("Failed to insert image path prediction for Naive Bayes: %v", err)
//	}
//
//	// Test case for an invalid algorithm
//	if err := dal.InsertPrediction("InvalidAlgorithm", queryIdentifier, "invalid_data.txt", "This should fail"); err == nil {
//		t.Error("Expected an error for an unrecognized algorithm, but got none")
//	}
//}
