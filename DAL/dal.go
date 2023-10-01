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

func main() {
	config := DBConfig{
		Username: "root",
		Password: "Matthew@1499.",
		HostName: "localhost:3306",
		Database: "goengine",
	}
	var err error
	db, err = Connection(config) // Notice the removal of the := which creates a new local variable
	if err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
	}
	defer CloseDb()
	err = CreateUserDAL("user3", "matt4", "mfa54", "dev", "pass2", true)
	if err != nil {
		log.Printf("Failed to create a user: %s", err)
	} else {
		log.Println("Successfully created user.")
	}

	err = UpdateUserDAL("user3", "Matt5", "mfa34", "adm", "newPass")
	if err != nil {
		log.Printf("Failed to update user: %s", err)
	} else {
		log.Println("Successfully updated user.")
	}
}
