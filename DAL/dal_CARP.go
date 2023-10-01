package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"time"
)

type Userdata struct {
	user_id         string
	user_name       string
	user_login      string
	user_role       string
	user_password   string
	active_or_not   bool
	user_date_added time.Time
}

// creating the user calling upon the sproc
func CreateUser(db *sql.DB, user_name string, user_login string, user_role string, user_password string, active_or_not bool) error {
	_, err := db.Exec("CALL create_user(?, ?, ?, ?, ?)", user_name, user_login, user_role, user_password, active_or_not)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// updating the user calling upon the sproc
func UpdateUser(db *sql.DB, user_id string, user_name string, user_login string, user_role string, user_password string) error {
	_, err := db.Exec("CALL update_user(?, ?, ?, ?, ?)", user_id, user_name, user_login, user_role, user_password)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// deleting the user
func DeleteUser(db *sql.DB, user_id string) error {
	_, err := db.Exec("CALL delete_user(?)", user_id)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// get users. scans rows
func GetUsers(db *sql.DB) ([]Userdata, error) {
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
		var dateTimeString string
		err := rows.Scan(
			&user.user_id,
			&user.user_name,
			&user.user_login,
			&user.user_role,
			&user.user_password,
			&user.active_or_not,
			&dateTimeString,
		)
		if err != nil {
			log.Println(err) // Handle the error appropriately
			continue
		}
		user.user_date_added, err = time.Parse("YYYY-MM-DD 15:01:05", dateTimeString)
		if err != nil {
			log.Println(err)
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

type JSON_Data_Connect struct {
	Username string
	Password string
	Hostname string
	Database string
}

// This function creates a connection to the database using given configuration.
func ConnectToDB(config JSON_Data_Connect) (*sql.DB, error) {
	connDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.Hostname, config.Database))
	if err != nil {
		return nil, err
	}

	err = connDB.Ping()
	if err != nil {
		return nil, err
	}

	return connDB, nil
}

// Main function for testing
//func main() {
//  // Mockup for database connection details. Replace with your real credentials.
//  config := JSON_Data_Connect{
//     Username: "root",
//     Password: "Matthew@1499.",
//     Hostname: "localhost:3306",
//     Database: "goengine",
//  }
//
//  db, err := ConnectToDB(config)
//  if err != nil {
//     log.Fatalf("Could not connect to database: %s", err)
//  }
//  defer db.Close()
//
//  // Test createUser function
//  err = CreateUser(db, "id1", "user1", "dev", "pas1")
//  if err != nil {
//     log.Printf("Failed to create user: %s", err)
//  } else {
//     log.Println("Successfully created user.")
//  }
//
//  // Test other functions similarly...
//
//  // For demonstration purposes, fetching a GitHub user
//  githubUser, err := GetGitHubUser("octocat")
//  if err != nil {
//     log.Printf("Failed to fetch GitHub user: %s", err)
//  } else {
//     log.Printf("Fetched GitHub user: %+v", githubUser)
//  }
//
//  //err = updateUser(db, "id1", "user2", "adm")
//  //if err != nil {
//  // log.Printf("Failed to update user: %s", err)
//  //
//  //} else {
//  // log.Println(("Successfully updated user"))
//  //}
//  err = deleteUser(db, "id1")
//  if err != nil {
//     log.Println("Failed to delete user: %s", err)
//
//  } else {
//
//     log.Println("Successfully deleted user")
//  }
//  users, err := getUsers(db)
//  if err != nil {
//     log.Println("Failed to retrieve users: %s", err)
//  } else {
//     log.Println("Successfully retrieved users: ")
//     for _, user := range users {
//        log.Println("User ID: %s, User Name: %s, User Login: %s,"+
//           " User Role: %s, User Password: %s, active or not: %s, user date added: %s",
//           user.user_id, user.user_name, user.user_login, user.user_role, user.user_password, user.active_or_not, user.user_date_added)
//     }
//  }
//}
