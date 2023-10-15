package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
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

// CreateSession creates a new session in the database for a given user with the provided token.
// It takes a userID and a token as parameters and returns an error if the operation fails.
func CreateSession(userID string, token string) error {
	// Execute the 'create_session' stored procedure with the provided userID and token.
	_, err := db.Exec("CALL create_session(?, ?)", userID, token)
	return err
}

// ValidateToken checks the validity of a token in the database and retrieves the associated userID and validation status.
// It takes a token as a parameter and returns the userID, a boolean indicating validity, and an error if any.
func ValidateToken(token string) (userID string, isValid bool, err error) {
	// Query the database using the 'validate_token' stored procedure with the provided token.
	err = db.QueryRow("CALL validate_token(?)", token).Scan(&userID, &isValid)
	return userID, isValid, err
}

// AddPermission adds a permission for a user role to perform a specific action on a given resource.
// It takes userRole, actionName, and resourceName as parameters and returns an error if the operation fails.
func AddPermission(userRole, actionName, resourceName string) error {
	// Execute the 'add_permission' stored procedure with the provided userRole, actionName, and resourceName.
	_, err := db.Exec("CALL add_permission(?, ?, ?)", userRole, actionName, resourceName)
	return err
}

// CheckPermission checks if a user role has permission to perform a specific action on a given resource.
// It takes userRole, actionName, and resourceName as parameters and returns a boolean indicating permission status and an error if any.
func CheckPermission(userRole, actionName, resourceName string) (bool, error) {
	// Query the database using the 'check_permission' stored procedure with the provided userRole, actionName, and resourceName.
	var hasPermission bool
	err := db.QueryRow("CALL check_permission(?, ?, ?)", userRole, actionName, resourceName).Scan(&hasPermission)
	return hasPermission, err
}

//

// CreateUser inserts a new user into the database.
func CreateUser(userName, userLogin, userRole string, userPassword string, activeOrNot bool) (string, error) {
	var userID string
	err := db.QueryRow("CALL create_user(?, ?, ?, AES_ENCRYPT(?, 'IST888IST888'), ?)", userName, userLogin, userRole, userPassword, activeOrNot).Scan(&userID)
	if err != nil {
		return "", err
	} else { // If no error, log the user ID
		log.Printf("User created with ID: %s", userID)
	}
	return userID, nil
}

// GetUserByID retrieves a specific user by their ID.
func GetUserByID(userID string) (*User, error) {
	var u User
	row := db.QueryRow("CALL get_user_by_ID(?)", userID)
	if err := row.Scan(&u.UserID, &u.UserName, &u.UserLogin, &u.UserRole, &u.UserPassword, &u.ActiveOrNot, &u.UserDateAdded); err != nil {
		return nil, err
	} else {
		log.Printf("Get User: %+v", u)
	}
	return &u, nil
}

// GetUsersByRole fetches all users with a specific role.
func GetUsersByRole(role string) ([]*User, error) {
	rows, err := db.Query("CALL get_users_by_role(?)", role)
	if err != nil {
		return nil, err
	} else {
		log.Printf("Open query for getting Users by Role: %+v", rows)
	}
	defer rows.Close()
	log.Printf("Closing Rows: %+v", rows)
	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.UserName, &u.UserLogin, &u.UserRole, &u.UserPassword, &u.ActiveOrNot, &u.UserDateAdded); err != nil {
			return nil, err
		}
		users = append(users, &u)
		log.Printf("Get Users by Role: %+v", u)
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
	log.Printf("Closing Rows: %+v", rows)
	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.UserName, &u.UserLogin, &u.UserRole, &u.UserPassword, &u.ActiveOrNot, &u.UserDateAdded); err != nil {
			return nil, err
		}
		users = append(users, &u)
		log.Printf("User: %+v", u)
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
	log.Printf("User ID: %s", userID, " was fetched successfully.")
	return userID, nil
}

// ValidateUserCredentials verifies user login details.
func ValidateUserCredentials(userLogin, userPassword string) (bool, error) {
	var isValid bool
	err := db.QueryRow("CALL validate_user(?, ?)", userLogin, userPassword).Scan(&isValid)
	log.Printf("User Login: %s", userLogin, " was validated successfully.")
	return isValid, err
}

// UpdateUser updates the details of a user.
func UpdateUser(userID, userName, userLogin, userRole, userPassword string) error {
	_, err := db.Exec("CALL update_user(?, ?, ?, ?, AES_ENCRYPT(?, 'IST888IST888'))", userID, userName, userLogin, userRole, userPassword)
	log.Printf("User: %s", userID, " was updated successfully.")
	return err
}

// DeleteUser removes a user from the database.
func DeleteUser(userID string) error {
	_, err := db.Exec("CALL delete_user(?)", userID)
	log.Printf("User: %s", userID, " was deleted successfully.")
	return err

}
