package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"log"
	"os"
	"testing"
)

// init initializes the program, reading the database configuration and establishing a connection
func init() {
	config, err := readJSONConfig("config.json")
	if err != nil {
		log.Fatal("Error reading JSON config:", err)
	}

	var connErr error
	db, connErr = Connection(config)
	if connErr != nil {
		log.Fatal("Error establishing database connection:", connErr)
	}
}
func TestMain(m *testing.M) {
	config, err := readJSONConfig("config.json")
	if err != nil {
		log.Fatal("Error reading JSON config:", err)
	}

	db, err := Connection(config)
	if err != nil {
		log.Fatal("Error establishing database connection:", err)
	}

	if db == nil {
		log.Fatal("Database connection is not initialized.")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin a transaction:", err)
	}

	_, err = tx.Exec("CALL goengine.user_registration(?, ?, ?, ?, ?)", "test_user", "test_login", "ADM", "test_password", true)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return
		}
		log.Fatalf("Failed to populate sample user: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	code := m.Run()
	os.Exit(code)
}

func TestAuthenticatingUser(t *testing.T) {
	var userId sql.NullString
	var userName sql.NullString
	var userRole sql.NullString

	err := db.QueryRow("CALL goengine.user_login(?, ?)", "test_login", "test_password").Scan(&userId, &userName, &userRole)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found, possibly add debug code here
		}
		t.Fatalf("Failed to authenticate user: %v", err)
	}
	if !userId.Valid || !userName.Valid || !userRole.Valid {
		// One of the returned values is NULL, possibly add debug code here
		t.Fatalf("Failed to get valid user information")
	}
}

func TestGettingUserRole(t *testing.T) {
	var userRole sql.NullString
	err := db.QueryRow("CALL goengine.get_user_role(?)", "ebcbee99-67a7-11ee-8fa6-4c796ed97681").Scan(&userRole) // Changed "test_user_id" to an actual UUID
	if err != nil {
		t.Fatalf("Failed to get user role: %v", err)
	}
	if !userRole.Valid {
		// Role is NULL, possibly add debug code here
		t.Fatalf("Failed to get a valid user role")
	}
}

func TestIsUserActive(t *testing.T) {
	var isActive sql.NullBool
	err := db.QueryRow("CALL goengine.is_user_active(?)", "b081bf7b-67aa-11ee-8fa6-4c796ed97681").Scan(&isActive) // Changed "test_user_id" to an actual UUID
	if err != nil {
		t.Fatalf("Failed to check if user is active: %v", err)
	}
	if !isActive.Valid || !isActive.Bool {
		// User is inactive, possibly add debug code here
		t.Fatalf("User is not active")
	}
}

func TestAuthorizeUser(t *testing.T) {
	var isAuthorized sql.NullBool
	err := db.QueryRow("CALL goengine.authorize_user(?, ?)", "b081bf7b-67aa-11ee-8fa6-4c796ed97681", "STD").Scan(&isAuthorized) // Changed "test_user_id" to an actual UUID
	if err != nil {
		t.Fatalf("Failed to authorize user: %v", err)
	}
	if !isAuthorized.Valid || !isAuthorized.Bool {
		// User is unauthorized, possibly add debug code here
		t.Fatalf("User is not authorized")
	}
}

func Connection(config JSON_Data_Connect) (*sql.DB, error) {
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
