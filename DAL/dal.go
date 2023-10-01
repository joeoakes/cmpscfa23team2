package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

type DBConfig struct {
    Username string
    Password string
    HostName string
    Database string
}

var db *sql.DB

func Connection(config DBConfig) (*sql.DB, error) {
    connDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Username, config.Password, config.HostName, config.Database))

    if err != nil {
       return nil, err
    }

    err = connDB.Ping()
    if err != nil {
       return nil, err
    }

    return connDB, nil
}

func CloseDb() {
    if db != nil {
       err := db.Close()
       if err != nil {
          log.Fatalf("Could not establish a connection with the database: %v", err)
       }
    }
}

func CreateUserDAL(user_id string, user_name string, user_login string, user_role string, user_password string, active_or_not bool) error {
    return CreateUser(db, user_name, user_login, user_role, user_password, active_or_not)
}

func UpdateUserDAL(user_id string, user_name string, user_login string, user_role string, user_password string) error {
    return UpdateUser(db, user_id, user_name, user_login, user_role, user_password)
}

func DeleteUsersDAL(user_id string) error {
    return DeleteUser(db, user_id)
}

// GetUser isnt working properly need to be updated and then can properly implement function here
//
//  func GetUsersDAL([]Userdata, error) {
//     return GetUsers(db)
//  }
func GetGitHubUserData(username string) (*GitHubUser, error) {

    return GetGitHubUser(username)

}
func main() {
    config := DBConfig{
       Username: "root",
       Password: "password",
       HostName: "localhost:3306",
       Database: "goengine",
    }
    var err error
    db, err = Connection(config) // Notice the removal of the := which creates a new local variable
    if err != nil {
       log.Fatalf("Failed to initialize database: %s", err)
    }
    defer CloseDb()
    //err = CreateUserDAL("user3", "matt4", "mfa54", "dev", "pass2", true)
    //if err != nil {
    // log.Printf("Failed to create a user: %s", err)
    //} else {
    // log.Println("Successfully created user.")
    //}

    err = UpdateUserDAL("2e52d4b4-6013-11ee-a9a7-00ffd9472baa", "Matt7", "mfa34", "adm", "newPass")
    if err != nil {
       log.Printf("Failed to update user: %s", err)
    } else {
       log.Println("Successfully updated user.")
    }
    err = DeleteUsersDAL("eed05ded-6012-11ee-a9a7-00ffd9472baa")
    if err != nil {
       log.Printf("Failed to delete user: @%s", err)
    } else {
       log.Println("Successfully deleted user.")
    }
    //users, err := GetUsersDAL()if err != nil {
    //
    // log.Printf("Failed to get users: %s", err)
    //
    //} else {
    //
    // log.Println("Successfully retrieved users:")
    //
    // for _, user := range users {
    //
    //    log.Printf("User ID: %s, User Name: %s, User Login: %s, User Role: %s, User Password: %s, Active: %v, User Date Added: %s",
    //
    //       user.user_id, user.user_name, user.user_login, user.user_role, user.user_password, user.active_or_not, user.user_date_added)
    //
    // }
    //
    //}

    githubUser, err := GetGitHubUser("Matt5")
    if err != nil {
       log.Printf("Failed to fetch GitHub user: %s", err)
    } else {
       log.Printf("Successfully fetched GitHub user: %+v", githubUser)
    }
}
 