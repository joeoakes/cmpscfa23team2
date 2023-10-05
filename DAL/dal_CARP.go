package DAL

import (
	"database/sql"
	"errors"
	//	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
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


// creating the user calling upon the sproc
func createUser(db *sql.DB, user_id string, user_name string, user_role string, user_password string) error {
	_, err := db.Exec("CALL create_user(?, ?, ?, ?, ?)", user_name, user_id, user_role, user_password, true)

	if err != nil {
		return "", err
	}


// get users. scans rows
func getUsers(db *sql.DB) ([]Userdata, error) {
	rows, err := db.Query("CALL get_users()")
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



func createUsertest(t *testing.T) {
	db, err := sql.Open("MySQL", "goengine")
	err = createUser(db, "user_id", "user_name", "user_role", "user_password")

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

