package dal

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"log"
)

// This code defines a struct called "User" with fields representing userID, name, login, role, password, active status, and date added.
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
//
// it creates a user in a database, logs the user ID if successful, and returns the user's ID or an error.
func CreateUser(userName, userLogin, userRole string, userPassword string, activeOrNot bool) (string, error) {
	var userID string
	err := DB.QueryRow("CALL create_user(?, ?, ?, ?, ?)", userName, userLogin, userRole, userPassword, activeOrNot).Scan(&userID)
	if err != nil {
		return "", err
	} else { // If no error, log the user ID
		log.Printf("User created with ID: %s", userID)
	}
	return userID, nil
}

// UpdateUser updates the details of a user.
//
// It defines a function "UpdateUser" that calls a stored procedure to update a user's information in a database, logs the user's ID, and returns any encountered error.
func UpdateUser(userID, userName, userLogin, userRole, userPassword string) error {
	_, err := DB.Exec("CALL update_user(?, ?, ?, ?, ?)", userID, userName, userLogin, userRole, userPassword)
	log.Printf("User: %s", userID)
	return err
}

// DeleteUser removes a user from the database.
//
// It defines a function that deletes a user with the given userID from a database using a stored procedure and logs the operation, returning any potential errors.
func DeleteUser(userID string) error {
	_, err := DB.Exec("CALL delete_user(?)", userID)
	log.Printf("User: %s", userID)
	return err

}

// This code defines a function that retrieves a user from a database using a stored procedure based on a given user login,
// and returns the user's information or an error.
func GetUserByLogin(userLogin string) (*User, error) {
	var u User
	row := DB.QueryRow("CALL get_user_by_login(?)", userLogin)
	if err := row.Scan(&u.UserID, &u.UserName, &u.UserLogin, &u.UserRole, &u.UserPassword, &u.ActiveOrNot, &u.UserDateAdded); err != nil {
		return nil, err
	} else {
		log.Printf("Get User: %+v", u)
	}
	return &u, nil
}

// GetUserByID retrieves a specific user by their ID.
//
// This code defines a function called GetUserByID that retrieves a user's information from a database by their ID
// and returns a pointer to a User struct along with an error.
func GetUserByID(userID string) (*User, error) {
	var u User
	row := DB.QueryRow("CALL get_user_by_ID(?)", userID)
	if err := row.Scan(&u.UserID, &u.UserName, &u.UserLogin, &u.UserRole, &u.UserPassword, &u.ActiveOrNot, &u.UserDateAdded); err != nil {
		return nil, err
	} else {
		log.Printf("Get User: %+v", u)
	}
	return &u, nil
}

// GetUsersByRole fetches all users with a specific role.
//
// This code defines a function that queries a database to retrieve a list of users by their role and logs various steps in the process,
// returning the list of users and any encountered errors.
func GetUsersByRole(role string) ([]*User, error) {
	rows, err := DB.Query("CALL get_users_by_role(?)", role)
	if err != nil {
		return nil, err
	} else {
		log.Printf("Open query for getting Users by Role: %+v", rows)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	log.Printf("Closing Rows: %+v", rows)
	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.UserName, &u.UserLogin, &u.UserRole, &u.UserPassword, &u.ActiveOrNot, &u.UserDateAdded); err != nil {
			return nil, err
		} else {
			log.Printf("Scan Rows: %+v", rows)
		}
		users = append(users, &u)
		log.Printf("Get Users by Role: %+v", u)
	}
	return users, rows.Err()
}

// GetAllUsers retrieves all registered users.
//
// This code defines a function, GetAllUsers, that retrieves user data from a database, processes it,
// and returns a  user objects while handling potential errors and resource cleanup.
func GetAllUsers() ([]*User, error) {
	rows, err := DB.Query("CALL get_users()")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
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
//
// This function retrieves a user's ID by calling a stored procedure in a database and logs the result, handling any errors that may occur.
func FetchUserIDByName(userName string) (string, error) {
	var userID string
	err := DB.QueryRow("CALL fetch_user_id(?)", userName).Scan(&userID)
	if err != nil {
		return "", err
	}
	log.Printf("User ID: %s", userID)
	return userID, nil
}
