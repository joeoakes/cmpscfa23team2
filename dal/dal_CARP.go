package dal

import (
	"database/sql"

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

// CreateUser inserts a new user into the database.
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
func UpdateUser(userID, userName, userLogin, userRole, userPassword string) error {
	_, err := DB.Exec("CALL update_user(?, ?, ?, ?, ?)", userID, userName, userLogin, userRole, userPassword)
	log.Printf("User: %s", userID)
	return err
}

// DeleteUser removes a user from the database.
func DeleteUser(userID string) error {
	_, err := DB.Exec("CALL delete_user(?)", userID)
	log.Printf("User: %s", userID)
	return err

}

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
func FetchUserIDByName(userName string) (string, error) {
	var userID string
	err := DB.QueryRow("CALL fetch_user_id(?)", userName).Scan(&userID)
	if err != nil {
		return "", err
	}
	log.Printf("User ID: %s", userID)
	return userID, nil
}
