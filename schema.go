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

	tuesTable := `
	CREATE TABLE IF NOT EXISTS tuesid (
		id VARCHAR(256) NOT NULL UNIQUE,
		used BOOL NOT NULL
	);
	`

	_, err := db.Exec(userTable)
	_, err = db.Exec(tuesTable)
	return err
}
