package dal_test

import (
	"cmpscfa23team2/dal"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secretpassword"
	hashed, err := dal.HashPassword(password)
	if err != nil {
		t.Errorf("Error occured: %v", err)
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		t.Errorf("Hashed password does not match the original password: %v", err)
	}
}

func TestComparePassword(t *testing.T) {
	password := "password123" //use any password you like
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := dal.ComparePassword(hashedPassword, password)
	if err != nil {
		t.Errorf("Passwords should match, but they don't: %v", err)

	}
	anotherPassword := "anotherPassword"
	err = dal.ComparePassword(hashedPassword, anotherPassword)
	if err == nil {
		t.Errorf("Passwords should not match, but they do.")
	}

}

func TestGenerateToken(t *testing.T) {
	// replace with test values
	userID := "testUserID"
	token, err := dal.GenerateToken(userID)
	if err != nil {
		t.Errorf("Token generation failed: %v", err)
	}
	valid, err := dal.ValidateToken(token)
	if err != nil {
		t.Errorf("Token validation failed: %v", err)

	}
	if !valid {
		t.Errorf("Token is not valid.")

	}
}

//func TestValidateToken(t *testing.T) {
//	// replace with actual token string
//	tokenString := "SECRETKEY123"
//	valid, err := DAL.ValidateToken(tokenString)
//	if err != nil {
//		t.Errorf("token validation failed: %v", err)
//
//	}
//	if !valid {
//		t.Errorf("token is not valid.")
//	}
//}

func TestValidateToken(t *testing.T) {
	// if it is a valid token
	validToken, _ := dal.GenerateToken("testUserID")
	isValid, err := dal.ValidateToken(validToken)
	assert.Nil(t, err)      //assert that there is no error
	assert.True(t, isValid) //assert that the token is valid

	// if it is not a valid token
	invalidToken := "invalid_token"
	isValid, err = dal.ValidateToken(invalidToken)
	assert.NotNil(t, err)    //assert that there is an error
	assert.False(t, isValid) // assert that the token is invalid
}

// Function does not pass the test because crypto/bcrypt: hashedSecret too short to be a bcrypted password
func TestAuthenticateUser(t *testing.T) {
	username := "johnpork"
	password := "he's'calling"

	// authenticate the user
	token, authErr := dal.AuthenticateUser(username, password)
	if authErr != nil {
		t.Errorf("Authentication failed: %v", authErr)

	}
	if token == "" {
		t.Errorf("Authentication succeeded, but the token is empty.")
	}
}

// this function doesn't pass the test because of sql: no rows in result set. Nothing in the table
func TestRefreshToken(t *testing.T) {
	// Replace with actual values
	oldRefreshToken := "valid_refresh_token"

	// Generate a new access token and refresh token using the old refresh token
	newAccessToken, newRefreshToken, err := dal.RefreshToken(oldRefreshToken)
	if err != nil {
		t.Errorf("Refresh token failed: %v", err)
	}

	// Check if the new access and refresh tokens are not empty
	if newAccessToken == "" {
		t.Errorf("New access token is empty.")
	}
	if newRefreshToken == "" {
		t.Errorf("New refresh token is empty.")
	}

	// Validate the new access token
	isValid, err := dal.ValidateToken(newAccessToken)
	if err != nil {
		t.Errorf("Token validation failed: %v", err)
	}
	if !isValid {
		t.Errorf("New access token is not valid.")
	}

	// Validate the new refresh token
	isValid, err = dal.ValidateToken(newRefreshToken)
	if err != nil {
		t.Errorf("Token validation failed: %v", err)
	}
	if !isValid {
		t.Errorf("New refresh token is not valid.")
	}
}

// RegisterUser function currently does not pass the test for foreign key violation. Trying to solve this problem
func TestRegisterUser(t *testing.T) {
	username := "JohnPork" //replace all of these with actual variables
	login := "hescalling"
	role := "NA"
	password := "jp123"
	userID, err := dal.RegisterUser(username, login, role, password, true)
	if err != nil {
		t.Errorf("User registration failed: %v", err)
	}
	if userID == "" { //check if userID is not empty
		t.Errorf("User registration succeeded, but userID is empty.")
	}
}

func TestLogoutUser(t *testing.T) {
	userID := "testUserID"
	err := dal.LogoutUser(userID)
	if err != nil {
		t.Errorf("User logout failed: %v", err)
	}
}

// function to change password
func TestChangePassword(t *testing.T) {
	userID := "testUserID"
	newPassword := "newPassword"
	err := dal.ChangePassword(userID, newPassword)
	if err != nil {
		t.Errorf("ChangePassword failed: %v", err)

	}
}