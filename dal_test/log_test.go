package dal_test

import (
	dal "cmpscfa23team2/dal"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestInsertOrUpdateStatusCode(t *testing.T) {
	err := dal.InsertOrUpdateStatusCode("POS", "noth")
	if err != nil {
		dal.InsertLog("400", "Failed to insert or update status code", "TestInsertOrUpdateStatusCode()")
		t.Fatalf("Failed to insert or update status code: %v", err)
	} else {
		dal.InsertLog("200", "Successfully inserted or updated status code", "TestInsertOrUpdateStatusCode()")
	}
}

func TestWriteLog(t *testing.T) {
	uniqueLogID := uuid.New().String()
	currentTime := time.Now()

	err := dal.WriteLog(uniqueLogID, "Pos", "Message logged successfully", "Engine1", currentTime)
	if err != nil {
		dal.InsertLog("400", "Failed to write log", "TestWriteLog()")
		t.Fatalf("Failed to write log: %v", err)
	} else {
		dal.InsertLog("200", "Successfully wrote log", "TestWriteLog()")
	}
}

func TestGetLog(t *testing.T) {
	logs, err := dal.GetLog()
	if err != nil {
		dal.InsertLog("400", "Failed to get logs", "TestGetLog()")
		t.Fatalf("Failed to get logs: %v", err)
	}
	for _, logItem := range logs {
		// Use reflect.DeepEqual to check for zero value
		dal.InsertLog("200", "Successfully got logs", "TestGetLog()")
		if reflect.DeepEqual(logItem, dal.Log{}) {
			t.Errorf("Log item is zero value")
		}
	}
}

func TestStoreLog(t *testing.T) {
	err := dal.StoreLog("200", "Stored using procedure", "Engine1")
	if err != nil {
		dal.InsertLog("400", "Failed to store log using stored procedure", "TestStoreLog()")
		t.Fatalf("Failed to store log using stored procedure: %v", err)
	} else {
		dal.InsertLog("200", "Successfully stored log using stored procedure", "TestStoreLog()")
	}
}
func TestGetSuccess(t *testing.T) {
	logs, err := dal.GetSuccess()
	if err != nil {
		dal.InsertLog("400", "Failed to get success logs", "TestGetSuccess()")
		t.Fatalf("Failed to get success logs: %v", err)
	}
	for _, logItem := range logs {
		if reflect.DeepEqual(logItem, dal.Log{}) {
			t.Errorf("Log item is zero value")
		}
	}
	dal.InsertLog("200", "Successfully got success logs", "TestGetSuccess()")
}
