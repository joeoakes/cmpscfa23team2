package DAL

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"testing"
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

// updating the user calling upon the sproc
func updateUser(db *sql.DB, user_id string, user_name string, user_role string) {
	_, err := db.Exec("CALL update_user(?, ?, ?)", user_id, user_name, user_role)
	if err != nil {
		log.Fatal(err)
	}
}

// creating the user calling upon the sproc
func createUser(db *sql.DB, user_id string, user_name string, user_role string, user_password string) error {
	_, err := db.Exec("CALL create_user(?, ?, ?, ?, ?)", user_name, user_id, user_role, user_password, true)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// deleting the user
func deleteUser(db *sql.DB, user_id string) error {
	_, err := db.Exec("CALL delete_user(?)", user_id)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// get users. scans rows
func getUsers(db *sql.DB) ([]Userdata, error) {
	rows, err := db.Query("CALL get_users()")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var users []Userdata

	for rows.Next() {
		var user Userdata
		err := rows.Scan(
			&user.user_id,
			&user.user_name,
			&user.user_role,
			&user.user_password,
			&user.active_or_not,
			&user.user_date_added,
		)
		if err != nil {
			log.Println(err) // Handle the error appropriately
			continue
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

type GitHubUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	Bio         string `json:"bio"`
	PublicRepos int    `json:"public_repos"`
}

// This function fetches information about a github user by their username.
func GetGitHubUser(username string) (*GitHubUser, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	//get request to api
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API request failed with status code %d", resp.StatusCode)
	}

	// decode json into github struct
	var user GitHubUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func createUsertest(t *testing.T) {
	db, err := sql.Open("MySQL", "goengine")
	err = createUser(db, "user_id", "user_name", "user_role", "user_password")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	user_id := "testUserID"
	user_name := "testUserName"
	user_role := "testUserRole"
	user_password := "testPassword"
	err = createUser(db, user_id, user_name, user_role, user_password)
	if err != nil {
		log.Fatal(err)

	}
	// does this work?
}
