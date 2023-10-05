package main

import (
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserID        string
	UserName      string
	UserLogin     string
	UserRole      string
	UserPassword  []byte
	ActiveOrNot   bool
	UserDateAdded string
}

// CreateUser inserts a new user into the database.
func CreateUser(userName, userLogin, userRole string, userPassword string, activeOrNot bool) (string, error) {
	var userID string
	err := db.QueryRow("CALL create_user(?, ?, ?, AES_ENCRYPT(?, 'IST888IST888'), ?)", userName, userLogin, userRole, userPassword, activeOrNot).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// GetUserByID retrieves a specific user by their ID.
func GetUserByID(userID string) (*User, error) {
	var u User
	row := db.QueryRow("CALL get_user_by_ID(?)", userID)
	if err := row.Scan(&u.UserID, &u.UserName, &u.UserLogin, &u.UserRole, &u.UserPassword, &u.ActiveOrNot, &u.UserDateAdded); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUsersByRole fetches all users with a specific role.
func GetUsersByRole(role string) ([]*User, error) {
	rows, err := db.Query("CALL get_users_by_role(?)", role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.UserName, &u.UserLogin, &u.UserRole, &u.UserPassword, &u.ActiveOrNot, &u.UserDateAdded); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}

// GetAllUsers retrieves all registered users.
func GetAllUsers() ([]*User, error) {
	rows, err := db.Query("CALL get_users()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.UserName, &u.UserLogin, &u.UserRole, &u.UserPassword, &u.ActiveOrNot, &u.UserDateAdded); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}

// FetchUserIDByName retrieves a user's ID using their username.
func FetchUserIDByName(userName string) (string, error) {
	var userID string
	err := db.QueryRow("CALL fetch_user_id(?)", userName).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// ValidateUserCredentials verifies user login details.
func ValidateUserCredentials(userLogin, userPassword string) (bool, error) {
	var isValid bool
	err := db.QueryRow("CALL validate_user(?, ?)", userLogin, userPassword).Scan(&isValid)
	return isValid, err
}

// UpdateUser updates the details of a user.
func UpdateUser(userID, userName, userLogin, userRole, userPassword string) error {
	_, err := db.Exec("CALL update_user(?, ?, ?, ?, AES_ENCRYPT(?, 'IST888IST888'))", userID, userName, userLogin, userRole, userPassword)
	return err
}

// DeleteUser removes a user from the database.
func DeleteUser(userID string) error {
	_, err := db.Exec("CALL delete_user(?)", userID)
	return err

}
