package carp

import (
	"cmpscfa23team2/dal"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// ValidateUserCredentials checks if the provided login credentials are valid.
// It returns the user's information if valid, otherwise returns an error.
func ValidateUserCredentials(userLogin, userPassword string) (*dal.User, error) {
	// Retrieve user information by login
	user, err := dal.GetUserByLogin(userLogin)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return nil, errors.New("authentication failed: invalid login credentials")
	}

	// Validate the password
	if !comparePasswords(userPassword, string(user.UserPassword)) {
		log.Printf("Authentication failed: invalid password for user %s", userLogin)
		return nil, errors.New("authentication failed: invalid login credentials")
	}

	return user, nil
}

// comparePasswords compares a hashed password with its plaintext version.
// It returns true if the passwords match, otherwise returns false.
func comparePasswords(plainPwd, hashedPwd string) bool {
	// Use bcrypt to compare hashed and plaintext passwords
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

// HashPassword hashes the provided plaintext password.
// It returns the hashed password or an error if hashing fails.
func HashPassword(plainPwd string) (string, error) {
	// Use bcrypt to hash the password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(plainPwd), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", errors.New("internal server error")
	}
	return string(hashedPwd), nil
}

// For Testing
func ExampleLogin(userLogin, userPassword string) {
	user, err := ValidateUserCredentials(userLogin, userPassword)
	if err != nil {
		log.Printf("Login failed: %v", err)
		return
	}

	log.Printf("Login successful. User: %+v", user)
	// You can perform additional actions after a successful login.
}
