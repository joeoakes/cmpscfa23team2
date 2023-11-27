package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

const SECRET_KEY = "SECRETKEY123!"

// Add the bcrypt hashing utility functions
//
// It defines function that hashes a provided password using the bcrypt hashing algorithm with a default cost and returns the hashed password as a byte slice or an error if encountered.
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// defines a function that compares a hashed password stored as a byte slice with a provided password string using the bcrypt library for secure password authentication.
func ComparePassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

// The code defines a function that authenticates a user by querying a database with a username and password,
// comparing the hashed password with the provided one, and generating a token for the user,
// returning the token or an error.
//
//	func AuthenticateUser(username string, password string) (string, error) {
//		var userID string
//		var hashedPassword []byte
//		var token string
//
//		err := DB.QueryRow("CALL authenticate_user(?, ?)", username, password).Scan(&userID, &hashedPassword)
//		if err != nil {
//			return "", err
//		}
//
//		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
//		if err != nil {
//			return "", err
//		}
//
//		token, err = GenerateToken(userID)
//		if err != nil {
//			return "", err
//		}
//
//		log.Printf("Generated token for user %s: %s", username, token)
//
//		return token, nil
//	}
func AuthenticateUser(username string, password string) (string, error) {
	var userID string
	var hashedPasswordStr string

	err := DB.QueryRow("CALL authenticate_user(?)", username).Scan(&userID, &hashedPasswordStr)
	if err != nil {
		log.Printf("Error in DB Query: %v", err)
		return "", err
	}

	hashedPassword := []byte(hashedPasswordStr) // Convert back to []byte

	if userID == "" {
		return "", fmt.Errorf("user not found")
	}

	// Verify that the hashed password has the correct bcrypt format
	if !strings.HasPrefix(hashedPasswordStr, "$2a$") {
		return "", fmt.Errorf("invalid bcrypt hash format")
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
		return "", err
	}

	token, err := GenerateToken(userID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	if token == "" {
		return "", fmt.Errorf("generated token is empty")
	}

	log.Printf("Generated token for user %s: %s", username, token)

	return token, nil
}

// This code generates a JWT token with a user ID and expiration time, using HMAC-SHA256 for signing.
func GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(), // 1 hour expiration
	})

	return token.SignedString([]byte(SECRET_KEY))
}

// This code generates a JWT token with a user ID and expiration time, using HMAC-SHA256 for signing.
//func GenerateToken(userID, userRole string) (string, error) {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"uid":  userID,
//		"role": userRole,
//		"exp":  time.Now().Add(time.Hour * 1).Unix(), // 1 hour expiration
//	})
//
//	return token.SignedString([]byte(SECRET_KEY))
//}

// This code defines a function that validates a JSON Web Token (JWT) by parsing it
// verifying its signature using a secret key,
// and checking its expiration time
func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("Failed to parse token claims")
	}

	return token.Valid && claims.VerifyExpiresAt(time.Now().Unix(), true), nil
}

// It defines a function called RefreshToken that takes an old refresh token as input, validates it against a database, generates a new access token and refresh token,
// and updates the database with the new refresh token
func RefreshToken(oldRefreshToken string) (string, string, error) {
	// This function should be used to generate a new access token given a valid refresh token.
	// For now, this just demonstrates the generation process.
	var userID string
	err := DB.QueryRow("CALL validate_refresh_token(?)", oldRefreshToken).Scan(&userID)
	if err != nil {
		return "", "", err
	}
	newAccessToken, err := GenerateToken(userID)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := GenerateToken(userID)
	if err != nil {
		return "", "", err
	}
	_, err = DB.Exec("CALL issue_refresh_token(?, ?)", userID, newRefreshToken)
	return newAccessToken, newRefreshToken, err
}

// This code defines a function called LogoutUser that takes a userID as a parameter and it uses the database connection.
// (DB) to execute a SQL stored procedure to log out a user with the specified userID,
// returning any potential errors encountered during the database operation.
func LogoutUser(userID string) error {
	_, err := DB.Exec("CALL logout_user(?)", userID)
	return err
}

// It defines a function "RegisterUser" that securely registers a user by hashing their password
// and storing their information in a database, returning a user ID or an error.
func RegisterUser(username string, login string, role string, password string, active bool) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	var userID string
	err = DB.QueryRow("CALL user_registration(?, ?, ?, ?, ?)", username, login, role, hashedPassword, active).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

// Takes a user ID and a new password as input and returns an error if there is any issue with the passowrd change process
func ChangePassword(userID string, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DB.Exec("CALL change_user_password(?, ?)", userID, hashedPassword)
	return err
}
