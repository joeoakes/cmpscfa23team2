package DAL_TEST

import (
	_ "github.com/go-sql-driver/mysql"
)

//func TestGetUserRole(t *testing.T) {
//	// Set up a test database and defer its cleanup
//	testDB = setupTestDB()
//	defer testDB.Close()
//
//	// Test the GetUserRole function
//	userID := "testuser"
//	expectedRole := "admin"
//
//	// Insert a test user and role into the test database
//	// This is a simplified example; in a real test, you would use your database setup.
//	_, err := testDB.Exec("INSERT INTO users (id, role) VALUES (?, ?)", userID, expectedRole)
//	if err != nil {
//		t.Fatalf("Failed to set up the test: %v", err)
//	}
//
//	// Call the GetUserRole function
//	role, err := GetUserRole(testDB, userID)
//	if err != nil {
//		t.Fatalf("GetUserRole failed: %v", err)
//	}
//
//	// Check if the retrieved role matches the expected role
//	if role != expectedRole {
//		t.Errorf("GetUserRole returned unexpected role. Expected: %s, Got: %s", expectedRole, role)
//	}
//}
//
//func TestIsUserActive(t *testing.T) {
//	// Set up a test database and defer its cleanup
//	testDB = setupTestDB()
//	defer testDB.Close()
//
//	// Test the IsUserActive function
//	userID := "testuser"
//
//	// Insert a test user and mark them as active in the test database
//	// This is a simplified example; in a real test, you would use your database setup.
//	_, err := testDB.Exec("INSERT INTO users (id, active) VALUES (?, true)", userID)
//	if err != nil {
//		t.Fatalf("Failed to set up the test: %v", err)
//	}
//
//	// Call the IsUserActive function
//	isActive, err := IsUserActive(testDB, userID)
//	if err != nil {
//		t.Fatalf("IsUserActive failed: %v", err)
//	}
//
//	// Check if the user is active
//	if !isActive {
//		t.Errorf("IsUserActive returned false for an active user")
//
