package DAL

import (
	"database/sql"
	"errors"
	//	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Userdata struct {
	user_id         string
	user_name       string
	user_role       string
	user_password   string
	active_or_not   bool
	user_date_added time.Time
}

// ReadSQLFile reads an SQL file located in the project directory and returns its content as a string.
// Most likely won't be utilizing this function as much (will be kept for test purposes).
func ReadSQLFile(filename string) (string, error) {
	// Get the absolute path to the project directory

	projectDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Construct the full path to the SQL file
	filePath := filepath.Join(projectDirectory, filename)
	// Read the content of the SQL file
	sqlContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(sqlContent), nil
}

// ValidUser checks if a user with the given login and password exists in the database.
// It returns true if a matching user is found, false otherwise, along with any errors encountered.
func ValidUser(db *sql.DB, userLogin, userPassword string) (bool, error) {
	// Read the SQL query from the "scripts.sql" file
	sqlQuery, err := ReadSQLFile("scripts.sql")
	if err != nil {
		return false, err
	}

	// Update the SQL query to reference the "users" table
	// For example, if your query is something like "SELECT * FROM users WHERE username = ? AND password = ?",
	// you don't need to change it since it already references the "users" table.

	// Prepare the SQL statement
	preparedStatement, err := db.Prepare(sqlQuery)
	if err != nil {
		return false, err
	}
	defer preparedStatement.Close()

	// Execute the query and check if a user with the given userLogin and userPassword exists
	var rowCount int
	err = preparedStatement.QueryRow(userLogin, userPassword).Scan(&rowCount)
	if err != nil {
		// Handle the error from the SQL query
		return false, err
	}

	if rowCount == 0 {
		// No matching user found, return a custom error
		return false, errors.New("username and password do not match")
	}

	// A matching user was found
	return true, nil
}

// updating the user calling upon the sproc
func updateUser(db *sql.DB, user_id string, user_name string, user_role string) {
	_, err := db.Exec("CALL update_user(?, ?, ?)", user_id, user_name, user_role)
	if err != nil {
		log.Fatal(err)
	}
}

// deleting the user
func deleteUser(db *sql.DB, user_id string) error {
	_, err := db.Exec("CALL delete_user(?)", user_id)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
