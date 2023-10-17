package main

import (
	"database/sql"
)

func Initialize(database *sql.DB) {
	db = database
}

// ... [Additional functions for other sprocs as required]
