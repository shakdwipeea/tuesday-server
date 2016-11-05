package main

import "database/sql"

func createSchema(db *sql.DB) error {
	userTable := `
	CREATE TABLE IF NOT EXISTS user (
		uid VARCHAR(256) NOT NULL UNIQUE,
		name VARCHAR(256),
		pic VARCHAR(256)
	);
	`

	_, err := db.Exec(userTable)
	return err
}
