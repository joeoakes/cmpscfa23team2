package dal

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const SECRET_KEY = "SECRETKEY123!"

// Add the bcrypt hashing utility functions
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ComparePassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func AuthenticateUser(username string, password string) (string, error) {
	var userID string
	var hashedPassword []byte
	var token string

	err := DB.QueryRow("CALL authenticate_user(?, ?)", username, password).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return "", err
	}

	token, err = GenerateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	return token.SignedString([]byte(SECRET_KEY))
}

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

func LogoutUser(userID string) error {
	_, err := DB.Exec("CALL logout_user(?)", userID)
	return err
}

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

func ChangePassword(userID string, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DB.Exec("CALL change_user_password(?, ?)", userID, hashedPassword)
	return err
}
