package dal

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

const SECRET_KEY = "SECRETKEY123!"

// Add the bcrypt hashing utility functions
//
// It defines function that hashes a provided password using the bcrypt hashing algorithm with a default cost and returns the hashed password as a byte slice or an error if encountered.
func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		InsertLog("400", "Failed to hash password", "HashPassword()")
		return nil, err
	}
	InsertLog("200", "Password hashed successfully", "HashPassword()")
	return hashedPassword, nil
}

// defines a function that compares a hashed password stored as a byte slice with a provided password string using the bcrypt library for secure password authentication.
func ComparePassword(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		InsertLog("400", "Password comparison failed", "ComparePassword()")
		return err
	}
	InsertLog("200", "Password compared successfully", "ComparePassword()")
	return nil
}

// AuthenticateUser authenticates a user by verifying their credentials.
//
// It takes a username and password as input, retrieves the hashed password from the database,
// and compares it with the provided password. If the credentials are valid, it generates a JWT token
// for the user and returns it. If authentication fails, it returns an error.
func AuthenticateUser(username string, password string) (string, error) {
	var userID, hashedPasswordStr string

	err := DB.QueryRow("CALL authenticate_user(?)", username).Scan(&userID, &hashedPasswordStr)
	if err != nil {
		InsertLog("400", "Error in DB Query during authentication", "AuthenticateUser()")
		return "", err
	}

	if userID == "" {
		InsertLog("400", "User not found during authentication", "AuthenticateUser()")
		return "", fmt.Errorf("user not found")
	}

	hashedPassword := []byte(hashedPasswordStr)

	if !strings.HasPrefix(hashedPasswordStr, "$2a$") {
		InsertLog("400", "Invalid bcrypt hash format", "AuthenticateUser()")
		return "", fmt.Errorf("invalid bcrypt hash format")
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		InsertLog("400", "Password comparison failed during authentication", "AuthenticateUser()")
		return "", err
	}

	token, err := GenerateToken(userID)
	if err != nil {
		InsertLog("400", "Error generating token during authentication", "AuthenticateUser()")
		return "", err
	}

	if token == "" {
		InsertLog("400", "Generated token is empty during authentication", "AuthenticateUser()")
		return "", fmt.Errorf("generated token is empty")
	}

	InsertLog("200", fmt.Sprintf("Generated token for user %s", username), "AuthenticateUser()")
	return token, nil
}

// This code generates a JWT token with a user ID and expiration time, using HMAC-SHA256 for signing.
func GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		InsertLog("400", "Error signing token", "GenerateToken()")
		return "", err
	}
	InsertLog("200", "Token generated successfully", "GenerateToken()")
	return tokenString, nil
}

// This code defines a function that validates a JSON Web Token (JWT) by parsing it
// verifying its signature using a secret key,
// and checking its expiration time
func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		InsertLog("400", "Error parsing token", "ValidateToken()")
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		InsertLog("400", "Failed to validate token", "ValidateToken()")
		return false, fmt.Errorf("failed to validate token")
	}

	InsertLog("200", "Token validated successfully", "ValidateToken()")
	return true, nil
}

// It defines a function called RefreshToken that takes an old refresh token as input, validates it against a database, generates a new access token and refresh token,
// and updates the database with the new refresh token
//func RefreshToken(oldRefreshToken string) (string, string, error) {
//	var userID, token string
//	var expiryBytes []byte
//	var expiry time.Time
//
//	// Retrieve the token from the database
//	err := DB.QueryRow("CALL validate_refresh_token(?)", userID).Scan(&userID, &token, &expiryBytes)
//	if err != nil {
//		return "", "", fmt.Errorf("error retrieving token: %v", err)
//	}
//
//	// Convert expiryBytes to a string
//	expiryStr := string(expiryBytes)
//
//	// Parse the string into time.Time
//	expiry, err = time.Parse("2006-01-02 15:04:05", expiryStr) // Adjust the layout string as per your date-time format
//	if err != nil {
//		return "", "", fmt.Errorf("error parsing expiry time: %v", err)
//	}
//
//	// Compare using bcrypt
//	if err := bcrypt.CompareHashAndPassword([]byte(token), []byte(oldRefreshToken)); err != nil {
//		return "", "", fmt.Errorf("token mismatch")
//	}
//
//	// Check expiry
//	if expiry.Before(time.Now()) {
//		return "", "", fmt.Errorf("token expired")
//	}
//	// Generate a new access token
//	newAccessToken, err := GenerateToken(userID)
//	if err != nil {
//		return "", "", err
//	}
//
//	// Generate a new refresh token
//	newRefreshTokenBytes, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()), bcrypt.DefaultCost)
//	if err != nil {
//		return "", "", err
//	}
//	newRefreshToken := string(newRefreshTokenBytes)
//
//	// Update the refresh token in the database
//	_, err = DB.Exec("CALL issue_refresh_token(?, ?)", userID, newRefreshToken)
//	return newAccessToken, newRefreshToken, err
//}

func ExtractToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// This code defines a function called LogoutUser that takes a userID as a parameter and it uses the database connection.
// (DB) to execute a SQL stored procedure to log out a user with the specified userID,
// returning any potential errors encountered during the database operation.
func LogoutUser(userID string) error {
	_, err := DB.Exec("CALL logout_user(?)", userID)
	if err != nil {
		InsertLog("400", "Failed to logout user", "LogoutUser()")
		return err
	}
	InsertLog("200", "User logged out successfully", "LogoutUser()")
	return nil
}

// It defines a function "RegisterUser" that securely registers a user by hashing their password
// and storing their information in a database, returning a user ID or an error.
func RegisterUser(username string, login string, role string, password string, active bool) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		InsertLog("400", "Failed to hash password during registration", "RegisterUser()")
		return "", err
	}

	var userID string
	err = DB.QueryRow("CALL user_registration(?, ?, ?, ?, ?)", username, login, role, hashedPassword, active).Scan(&userID)
	if err != nil {
		InsertLog("400", "Failed to register user", "RegisterUser()")
		return "", err
	}
	InsertLog("200", "User registered successfully", "RegisterUser()")

	return userID, nil
}

// Takes a user ID and a new password as input and returns an error if there is any issue with the passowrd change process
func ChangePassword(userID string, newPassword string) error {
	// Generate a hashed password from the new password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		InsertLog("400", "Error generating hashed password during password change", "ChangePassword()")
		return err
	}

	// Update the user's password in the database.
	_, err = DB.Exec("CALL change_user_password(?, ?)", userID, hashedPassword)
	if err != nil {
		InsertLog("400", "Error updating password in the database during password change", "ChangePassword()")
		return err
	}

	InsertLog("200", fmt.Sprintf("Password changed for user with ID %s", userID), "ChangePassword()")
	return nil
}
