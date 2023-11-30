package dal

import (
	"cmpscfa23team2/dal"
	"database/sql"
	"errors"
	"testing"
)

// test creating users
func TestCreateUser(t *testing.T) {
	user, err := dal.CreateUser("johnpork", "jp514", "DEV", "resister", true)
	if err != nil {
		dal.InsertLog("400", "Failed to create user", "TestCreateUser()")
		t.Errorf("Expected no error, but got an error: %v", err)
	} else {
		dal.InsertLog("200", "Successfully created user", "TestCreateUser()")
	}

	if user == "" {
		t.Errorf("Expected a user, but got an empty string.")
	}
}

// test updating user
func TestUpdateUser(t *testing.T) {
	err := dal.UpdateUser("a0s901xcamkap1985", "johnpork", "jp513", "DEV", "Resister")
	if err != nil {
		dal.InsertLog("400", "Failed to update user", "TestUpdateUser()")
		t.Errorf("Couldn't update user: %v", err)
	} else {
		dal.InsertLog("200", "Successfully updated user", "TestUpdateUser()")
	}
}

// test delete user
func TestDeleteUser(t *testing.T) {
	err := dal.DeleteUser("da53655b-7c53-11ee-aa3b-6c2b59772aba")
	if err != nil {
		dal.InsertLog("400", "Failed to delete user", "TestDeleteUser()")
		t.Errorf("Couldn't delete user: %v", err)
	} else {
		dal.InsertLog("200", "Successfully deleted user", "TestDeleteUser()")
	}
}

// test getting a user by login
func TestGetUserByLogin(t *testing.T) {
	validUserLogin := "jxo19"
	// test for valid user login
	user, err := dal.GetUserByLogin(validUserLogin)
	if err != nil {
		dal.InsertLog("400", "Failed to get user by login", "TestGetUserByLogin()")
		t.Errorf("Expected no error, but got an error : %v", err)
	} else {
		dal.InsertLog("200", "Successfully got user by login", "TestGetUserByLogin()")
	}

	if user == nil {
		t.Errorf("Expected a user, but got nil.")
	}

	// test for invalid user login
	invalidUserLogin := "nonexist"
	user, err = dal.GetUserByLogin(invalidUserLogin)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, but got error: %v", err)
		dal.InsertLog("400", "Failed to get user by login", "TestGetUserByLogin()")
	} else {
		dal.InsertLog("200", "Successfully got user by login", "TestGetUserByLogin()")
	}

	if user != nil {
		t.Errorf("Expected nil user for nonexistent user login.")
	}
}

// test getting users by ID.
// User id changes every time you run the SQL scripts, so make sure to change ID
func TestGetUserByID(t *testing.T) {
	validUserID := "07f70456-8f2e-11ee-ae02-30d042e80ac3"
	user, err := dal.GetUserByID(validUserID)
	if err != nil {
		dal.InsertLog("400", "Failed to get user by ID", "TestGetUserByID()")
		t.Errorf("Expected no error, but got an error: %v", err)
	} else {
		dal.InsertLog("200", "Successfully got user by ID", "TestGetUserByID()")
	}

	if user == nil {
		t.Errorf("Expected a user, but got nil.")
	}

	invalidUserID := "298s0aois-s13-sa"
	user, err = dal.GetUserByID(invalidUserID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, but got error: %v", err)
		dal.InsertLog("400", "Failed to get user by ID", "TestGetUserByID()")
	} else {
		dal.InsertLog("200", "Successfully got user by ID", "TestGetUserByID()")
	}

	if user != nil {
		t.Errorf("Expected nil user for nonexistent user ID")
	}
}

// test getting users by role
func TestGetUsersByRole(t *testing.T) {
	validRole := "DEV"
	users, err := dal.GetUsersByRole(validRole)
	if err != nil {
		dal.InsertLog("400", "Failed to get users by role", "TestGetUsersByRole()")
		t.Errorf("Expected no error, but got an error: %v", err)
	} else {
		dal.InsertLog("200", "Successfully got users by role", "TestGetUsersByRole()")
	}

	if len(users) == 0 {
		t.Errorf("Expected at least one user, but got none.")
	}
}

// test getting all users
func TestGetAllUsers(t *testing.T) {
	users, err := dal.GetAllUsers()
	if err != nil {
		dal.InsertLog("400", "Failed to get all users", "TestGetAllUsers()")
		t.Errorf("Expected no error, but got an error: %v", err)
	} else {
		dal.InsertLog("200", "Successfully got all users", "TestGetAllUsers()")
	}

	if len(users) == 0 {
		t.Errorf("Expected at least one user, but got none.")
	}
}

// test fetching a user's ID by name
func TestFetchUserIDByName(t *testing.T) {
	validusername := "Joshua Ferrell"
	userID, err := dal.FetchUserIDByName(validusername)
	if err != nil {
		dal.InsertLog("400", "Failed to fetch user ID by name", "TestFetchUserIDByName()")
		t.Errorf("Expected no error, but got an error: %v", err)
	} else {
		dal.InsertLog("200", "Successfully fetched user ID by name", "TestFetchUserIDByName()")
	}

	if userID == "" {
		t.Errorf("Expected a user ID, but got an empty string.")
	}
}
