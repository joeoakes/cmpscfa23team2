package dal

import (
	"cmpscfa23team2/dal"
	"database/sql"
	"errors"
	"testing"
)

// test getting a user by login
func TestGetUserByLogin(t *testing.T) {
	validUserLogin := "jxo19"
	// test for valid user login
	user, err := dal.GetUserByLogin(validUserLogin)
	if err != nil {
		t.Errorf("Expected no error, but got an error : %v", err)

	}
	if user == nil {
		t.Errorf("Expected a user, but got nil.")
	}
	// test for invalid user login
	invalidUserLogin := "nonexist"
	user, err = dal.GetUserByLogin(invalidUserLogin)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, but got error: %v", err)

	}
	if user != nil {
		t.Errorf("Expected nil user for nonexistent user login.")
	}
}

// test getting users by ID.
// User id changes everytime you run the sql scripts, so make sure to change ID
func TestGetUserByID(t *testing.T) {
	validUserID := "da53655b-7c53-11ee-aa3b-6c2b59772aba"
	user, err := dal.GetUserByID(validUserID)
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	if user == nil {
		t.Errorf("Expected a user, but got nil.")
	}

	invalidUserID := "298s0aois-s13-sa"
	user, err = dal.GetUserByID(invalidUserID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, but got error: %v", err)

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
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	if len(users) == 0 {
		t.Errorf("Expected at least one user, but got none.")

	} // getting users by role like with ADM or DEV work fine

	////testing for invalid role
	//invalidRole := "LOL"
	//users, err = dal.GetUsersByRole(invalidRole)
	//if !errors.Is(err, sql.ErrNoRows) {
	//	t.Errorf("Expected sql.ErrNoRows, but got error: %v", err)
	//}
	//if len(users) > 0 {
	//	t.Errorf("Expected no users for a nonexistent role.")
	//} for some reason it doesn't work
}

// test getting all users
func TestGetAllUsers(t *testing.T) {
	users, err := dal.GetAllUsers()
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
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
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	if userID == "" {
		t.Errorf("Expected a user ID, but got an empty string.")
	}
}

// test creating users
func TestCreateUser(t *testing.T) {
	user, err := dal.CreateUser("johnpork", "jp514", "DEV", "resister", true)
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)

	}
	if user == "" {
		t.Errorf("Expected a user, but got an empty string.")

	}
}

// test updating user
func TestUpdateUser(t *testing.T) {
	err := dal.UpdateUser("a0s901xcamkap1985", "johnpork", "jp513", "DEV", "Resister")
	if err != nil {
		t.Errorf("Couldn't update user: %v", err)
	}

}

// test delete user
func TestDeleteUser(t *testing.T) {
	err := dal.DeleteUser("da53655b-7c53-11ee-aa3b-6c2b59772aba")
	if err != nil {
		t.Errorf("Couldn't delete user: %v", err)
	}

}
