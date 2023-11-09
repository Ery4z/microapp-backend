package main

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "db.sqlite")
	if err != nil {
		log.Fatal("Failed to instantiate the db object: " + err.Error())
	}

	initGroupsTable()
}

func initGroupsTable() {
	// Create the groups table if it doesn't exist.
	createTableSQL := `CREATE TABLE IF NOT EXISTS groups (
        groupId TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT
    );`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating groups table: %s", err)
	}
}
