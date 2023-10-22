package DAL_TEST

import (
	"cmpscfa23team2/DAL"
	"github.com/google/uuid"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Setup: Initialize the database
	err := DAL.InitDB()
	if err != nil {
		panic("Failed to initialize the database: " + err.Error())
	}

	// Run all tests in the package
	code := m.Run()

	// Teardown: Close the database
	DAL.CloseDb()

	os.Exit(code)
}

func TestInsertOrUpdateStatusCode(t *testing.T) {
	err := DAL.InsertOrUpdateStatusCode("POS", "noth")
	if err != nil {
		DAL.InsertLog("400", "Failed to insert or update status code", "TestInsertOrUpdateStatusCode()")
		t.Fatalf("Failed to insert or update status code: %v", err)
	} else {
		DAL.InsertLog("200", "Successfully inserted or updated status code", "TestInsertOrUpdateStatusCode()")
	}
}

func TestFetchUserIDByName(t *testing.T) {
	_, err := DAL.FetchUserIDByName("Joesph Oakes")
	if err != nil {
		DAL.InsertLog("400", "Failed to fetch user ID", "TestFetchUserIDByName()")
		t.Fatalf("Failed to fetch user ID: %v", err)
	} else {
		DAL.InsertLog("200", "Successfully fetched user ID", "TestFetchUserIDByName()")
	}
}

func TestUpdateUser(t *testing.T) {
	err := DAL.UpdateUser("NewName", "jxo19", "ADM", "ADM", "password")
	if err != nil {
		DAL.InsertLog("400", "Failed to update user", "TestUpdateUser()")
		t.Fatalf("Failed to update user: %v", err)
	} else {
		DAL.InsertLog("200", "Successfully updated user", "TestUpdateUser()")
	}
}

func TestDeleteUser(t *testing.T) {
	err := DAL.DeleteUser("jxo19")
	if err != nil {
		DAL.InsertLog("400", "Failed to delete user", "TestDeleteUser()")
		t.Fatalf("Failed to delete user: %v", err)
	} else {
		DAL.InsertLog("200", "Successfully deleted user", "TestDeleteUser()")
	}
}

func TestWriteLog(t *testing.T) {
	uniqueLogID := uuid.New().String()
	currentTime := time.Now()

	err := DAL.WriteLog(uniqueLogID, "Pos", "Message logged successfully", "Engine1", currentTime)
	if err != nil {
		DAL.InsertLog("400", "Failed to write log", "TestWriteLog()")
		t.Fatalf("Failed to write log: %v", err)
	} else {
		DAL.InsertLog("200", "Successfully wrote log", "TestWriteLog()")
	}
}

func TestGetLog(t *testing.T) {
	logs, err := DAL.GetLog()
	if err != nil {
		DAL.InsertLog("400", "Failed to get logs", "TestGetLog()")
		t.Fatalf("Failed to get logs: %v", err)
	}
	for _, logItem := range logs {
		// Use reflect.DeepEqual to check for zero value
		DAL.InsertLog("200", "Successfully got logs", "TestGetLog()")
		if reflect.DeepEqual(logItem, DAL.Log{}) {
			t.Errorf("Log item is zero value")
		}
	}
}

func TestStoreLog(t *testing.T) {
	err := DAL.StoreLog("200", "Stored using procedure", "Engine1")
	if err != nil {
		DAL.InsertLog("400", "Failed to store log using stored procedure", "TestStoreLog()")
		t.Fatalf("Failed to store log using stored procedure: %v", err)
	} else {
		DAL.InsertLog("200", "Successfully stored log using stored procedure", "TestStoreLog()")
	}
}

func TestCreateUser(t *testing.T) {
	_, err := DAL.CreateUser("John", "john123", "ADM", "password", true)
	if err != nil {
		DAL.InsertLog("400", "Failed to create a new user", "TestCreateUser()")
		t.Fatalf("Failed to create a new user: %v", err)
	} else {
		DAL.InsertLog("200", "Successfully created a new user", "TestCreateUser()")
	}
}
